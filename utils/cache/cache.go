/*
 * TTLCache - an in-memory cache with expiration
 * TTLCache is a simple key/value cache in golang with the following functions:
 *
 * 1. Thread-safe
 * 2. Individual expiring time or global expiring time, you can choose
 * 3. Auto-Extending expiration on Get -or- DNS style TTL, see SkipTtlExtensionOnHit(bool)
 * 4. Fast and memory efficient
 * 5. Can trigger callback on key expiration
 * 6. Cleanup resources by calling Close() at end of lifecycle.
 *
 * Based on https://github.com/ReneKroon/ttlcache
 */

package cache

import (
	"sync"
	"time"
)

// checkExpireCallback is a function type used as a callback to externally verify if a cached item should expire.
// This allows for custom expiration logic beyond the standard TTL.
type checkExpireCallback[K comparable, T any] func(key K, value T) bool

// expireCallback is a function type used as a callback when a cached item expires or when a new item is added.
// It can be used for logging, cleanup, or other notification purposes.
type expireCallback[K comparable, T any] func(key K, value T)

// Cache is a thread-safe, in-memory key/value store with support for time-to-live (TTL) expiration.
// It allows for both global and per-item TTL settings and provides callbacks for various events.
//
// Type Parameters:
//
//	K: The type of the keys, which must be comparable.
//	T: The type of the values stored in the cache.
type Cache[K comparable, T any] struct {
	mutex                  sync.Mutex
	ttl                    time.Duration
	items                  map[K]*cachedItem[K, T]
	expireCallback         expireCallback[K, T]
	checkExpireCallback    checkExpireCallback[K, T]
	newItemCallback        expireCallback[K, T]
	priorityQueue          *priorityQueue[K, T]
	expirationNotification chan bool
	expirationTime         time.Time
	skipTTLExtension       bool
	shutdownSignal         chan (chan struct{})
	isShutDown             bool
}

// getItem retrieves an item from the cache, handling expiration and TTL extension.
// It is not thread-safe and must be called within a locked mutex.
func (cache *Cache[K, T]) getItem(key K) (*cachedItem[K, T], bool, bool) {
	item, exists := cache.items[key]
	if !exists || item.expired() {
		return nil, false, false
	}

	// If the item has a TTL, update its expiration time.
	if item.ttl >= 0 && (item.ttl > 0 || cache.ttl > 0) {
		if cache.ttl > 0 && item.ttl == 0 {
			item.ttl = cache.ttl
		}

		if !cache.skipTTLExtension {
			item.touch()
		}
		cache.priorityQueue.update(item)
	}

	// Check if an expiration notification is needed.
	expirationNotification := cache.expirationTime.After(time.Now().Add(item.ttl))
	return item, exists, expirationNotification
}

// startExpirationProcessing runs in a background goroutine to handle the expiration of items.
// It uses a priority queue to efficiently find the next item to expire.
func (cache *Cache[K, T]) startExpirationProcessing() {
	timer := time.NewTimer(time.Hour)
	for {
		sleepTime := cache.getSleepTime()
		timer.Reset(sleepTime)

		select {
		case shutdownFeedback := <-cache.shutdownSignal:
			timer.Stop()
			shutdownFeedback <- struct{}{}
			return
		case <-timer.C:
			timer.Stop()
			cache.processExpiredItems()
		case <-cache.expirationNotification:
			timer.Stop()
		}
	}
}

// getSleepTime calculates the duration until the next item expires.
func (cache *Cache[K, T]) getSleepTime() time.Duration {
	cache.mutex.Lock()
	defer cache.mutex.Unlock()

	if cache.priorityQueue.Len() == 0 {
		if cache.ttl > 0 {
			return cache.ttl
		} else {
			return time.Hour
		}
	}

	item := cache.priorityQueue.items[0]
	sleepTime := time.Until(item.expireAt)

	if sleepTime < 0 {
		sleepTime = time.Microsecond // Expired, process immediately.
	}
	if cache.ttl > 0 && sleepTime > cache.ttl {
		sleepTime = cache.ttl
	}

	cache.expirationTime = time.Now().Add(sleepTime)
	return sleepTime
}

// processExpiredItems removes expired items from the cache and triggers callbacks.
func (cache *Cache[K, T]) processExpiredItems() {
	cache.mutex.Lock()
	defer cache.mutex.Unlock()

	for i := 0; i < cache.priorityQueue.Len(); {
		item := cache.priorityQueue.items[i]
		if !item.expired() {
			break
		}

		if cache.checkExpireCallback != nil && !cache.checkExpireCallback(item.key, item.data) {
			item.touch()
			cache.priorityQueue.update(item)
			i++
			continue
		}

		cache.priorityQueue.remove(item)
		delete(cache.items, item.key)

		if cache.expireCallback != nil {
			go cache.expireCallback(item.key, item.data)
		}
	}
}

// Close gracefully shuts down the cache, stopping the expiration processing goroutine and clearing all items.
// It is safe to call Close multiple times.
func (cache *Cache[K, T]) Close() {
	cache.mutex.Lock()
	if cache.isShutDown {
		cache.mutex.Unlock()
		return
	}
	cache.isShutDown = true
	cache.mutex.Unlock()

	feedback := make(chan struct{})
	cache.shutdownSignal <- feedback
	<-feedback
	close(cache.shutdownSignal)

	cache.Purge()
}

// Set adds or updates an item in the cache with the global TTL.
func (cache *Cache[K, T]) Set(key K, data T) {
	cache.SetWithTTL(key, data, ItemExpireWithGlobalTTL)
}

// SetWithTTL adds or updates an item in the cache with a specific TTL.
func (cache *Cache[K, T]) SetWithTTL(key K, data T, ttl time.Duration) {
	cache.mutex.Lock()
	defer cache.mutex.Unlock()

	item, exists, _ := cache.getItem(key)

	if exists {
		item.data = data
		item.ttl = ttl
	} else {
		item = newItem(key, data, ttl)
		cache.items[key] = item
	}

	if item.ttl >= 0 && (item.ttl > 0 || cache.ttl > 0) {
		if cache.ttl > 0 && item.ttl == 0 {
			item.ttl = cache.ttl
		}
		item.touch()
	}

	if exists {
		cache.priorityQueue.update(item)
	} else {
		cache.priorityQueue.push(item)
	}

	if !exists && cache.newItemCallback != nil {
		go cache.newItemCallback(key, data)
	}
	cache.expirationNotification <- true
}

// Get retrieves an item from the cache. It also extends the item's TTL unless disabled.
func (cache *Cache[K, T]) Get(key K) (T, bool) {
	cache.mutex.Lock()
	item, exists, triggerExpirationNotification := cache.getItem(key)
	var dataToReturn T
	if exists {
		dataToReturn = item.data
	}
	cache.mutex.Unlock()

	if triggerExpirationNotification {
		cache.expirationNotification <- true
	}
	return dataToReturn, exists
}

// Remove deletes an item from the cache.
func (cache *Cache[K, T]) Remove(key K) bool {
	cache.mutex.Lock()
	defer cache.mutex.Unlock()

	object, exists := cache.items[key]
	if !exists {
		return false
	}

	delete(cache.items, object.key)
	cache.priorityQueue.remove(object)
	return true
}

// Count returns the number of items currently in the cache.
func (cache *Cache[K, T]) Count() int {
	cache.mutex.Lock()
	defer cache.mutex.Unlock()
	return len(cache.items)
}

// SetTTL sets the global TTL for the cache.
func (cache *Cache[K, T]) SetTTL(ttl time.Duration) {
	cache.mutex.Lock()
	defer cache.mutex.Unlock()
	cache.ttl = ttl
	cache.expirationNotification <- true
}

// SetExpirationCallback sets the callback function to be executed when an item expires.
func (cache *Cache[K, T]) SetExpirationCallback(callback expireCallback[K, T]) {
	cache.expireCallback = callback
}

// SetCheckExpirationCallback sets a callback that allows for external validation before an item expires.
func (cache *Cache[K, T]) SetCheckExpirationCallback(callback checkExpireCallback[K, T]) {
	cache.checkExpireCallback = callback
}

// SetNewItemCallback sets the callback function to be executed when a new item is added to the cache.
func (cache *Cache[K, T]) SetNewItemCallback(callback expireCallback[K, T]) {
	cache.newItemCallback = callback
}

// SkipTtlExtensionOnHit configures whether the TTL of an item is extended on a cache hit.
func (cache *Cache[K, T]) SkipTtlExtensionOnHit(value bool) {
	cache.skipTTLExtension = value
}

// Purge removes all items from the cache.
func (cache *Cache[K, T]) Purge() {
	cache.mutex.Lock()
	defer cache.mutex.Unlock()
	cache.items = make(map[K]*cachedItem[K, T])
	cache.priorityQueue = newPriorityQueue[K, T]()
}

// NewTtlCache creates and returns a new Cache instance.
func NewTtlCache[K comparable, T any]() *Cache[K, T] {
	cache := &Cache[K, T]{
		items:                  make(map[K]*cachedItem[K, T]),
		priorityQueue:          newPriorityQueue[K, T](),
		expirationNotification: make(chan bool, 1),
		shutdownSignal:         make(chan chan struct{}),
	}
	go cache.startExpirationProcessing()
	return cache
}

// Load is an alias for Get.
func (cache *Cache[K, T]) Load(key K) (T, bool) {
	return cache.Get(key)
}

// Store is an alias for Set.
func (cache *Cache[K, T]) Store(key K, value T) {
	cache.Set(key, value)
}

// StoreWithTTL is an alias for SetWithTTL.
func (cache *Cache[K, T]) StoreWithTTL(key K, value T, ttl time.Duration) {
	cache.SetWithTTL(key, value, ttl)
}

// Delete is an alias for Remove.
func (cache *Cache[K, T]) Delete(key K) {
	cache.Remove(key)
}

// Range iterates over all items in the cache and applies a callback.
func (cache *Cache[K, T]) Range(cb func(k K, v T) bool) {
	cache.mutex.Lock()
	defer cache.mutex.Unlock()
	for k, v := range cache.items {
		if !cb(k, v.data) {
			break
		}
	}
}

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
	"fmt"
	"sync"
	"time"
)

// CheckExpireCallback is used as a callback for an external check on cachedItem expiration
type checkExpireCallback func(key string, value any) bool

// ExpireCallback is used as a callback on cachedItem expiration or when notifying of an cachedItem new to the cache
type expireCallback func(key string, value any)

// Cache is a synchronized map of items that can auto-expire once stale
type Cache struct {
	mutex                  sync.Mutex
	ttl                    time.Duration
	items                  map[string]*cachedItem
	expireCallback         expireCallback
	checkExpireCallback    checkExpireCallback
	newItemCallback        expireCallback
	priorityQueue          *priorityQueue
	expirationNotification chan bool
	expirationTime         time.Time
	skipTTLExtension       bool
	shutdownSignal         chan (chan struct{})
	isShutDown             bool
}

func (cache *Cache) getItem(key string) (*cachedItem, bool, bool) {
	item, exists := cache.items[key]
	if !exists || item.expired() {
		return nil, false, false
	}

	if item.ttl >= 0 && (item.ttl > 0 || cache.ttl > 0) {
		if cache.ttl > 0 && item.ttl == 0 {
			item.ttl = cache.ttl
		}

		if !cache.skipTTLExtension {
			item.touch()
		}
		cache.priorityQueue.update(item)
	}

	expirationNotification := false
	if cache.expirationTime.After(time.Now().Add(item.ttl)) {
		expirationNotification = true
	}
	return item, exists, expirationNotification
}

func (cache *Cache) startExpirationProcessing() {
	timer := time.NewTimer(time.Hour)
	for {
		var sleepTime time.Duration
		cache.mutex.Lock()
		if cache.priorityQueue.Len() > 0 {
			sleepTime = time.Until(cache.priorityQueue.items[0].expireAt)
			if sleepTime < 0 && cache.priorityQueue.items[0].expireAt.IsZero() {
				sleepTime = time.Hour
			} else if sleepTime < 0 {
				sleepTime = time.Microsecond
			}
			if cache.ttl > 0 {
				sleepTime = min(sleepTime, cache.ttl)
			}

		} else if cache.ttl > 0 {
			sleepTime = cache.ttl
		} else {
			sleepTime = time.Hour
		}

		cache.expirationTime = time.Now().Add(sleepTime)
		cache.mutex.Unlock()

		timer.Reset(sleepTime)
		select {
		case shutdownFeedback := <-cache.shutdownSignal:
			timer.Stop()
			shutdownFeedback <- struct{}{}
			return
		case <-timer.C:
			timer.Stop()
			cache.mutex.Lock()
			if cache.priorityQueue.Len() == 0 {
				cache.mutex.Unlock()
				continue
			}

			// index will only be advanced if the current entry will not be evicted
			i := 0
			for item := cache.priorityQueue.items[i]; item.expired(); item = cache.priorityQueue.items[i] {

				if cache.checkExpireCallback != nil {
					if !cache.checkExpireCallback(item.key, item.data) {
						item.touch()
						cache.priorityQueue.update(item)
						i++
						if i == cache.priorityQueue.Len() {
							break
						}
						continue
					}
				}

				cache.priorityQueue.remove(item)
				delete(cache.items, item.key)
				if cache.expireCallback != nil {
					go cache.expireCallback(item.key, item.data)
				}
				if cache.priorityQueue.Len() == 0 {
					goto done
				}
			}
		done:
			cache.mutex.Unlock()

		case <-cache.expirationNotification:
			timer.Stop()
			continue
		}
	}
}

// Close calls Purge, and then stops the goroutine that does ttl checking, for a clean shutdown.
// The cache is no longer cleaning up after the first call to Close, repeated calls are safe though.
func (cache *Cache) Close() {

	cache.mutex.Lock()
	if !cache.isShutDown {
		cache.isShutDown = true
		cache.mutex.Unlock()
		feedback := make(chan struct{})
		cache.shutdownSignal <- feedback
		<-feedback
		close(cache.shutdownSignal)
	} else {
		cache.mutex.Unlock()
	}
	cache.Purge()
}

// Set is a thread-safe way to add new items to the map
func (cache *Cache) Set(key string, data any) {
	cache.SetWithTTL(key, data, ItemExpireWithGlobalTTL)
}

// SetWithTTL is a thread-safe way to add new items to the map with individual ttl
func (cache *Cache) SetWithTTL(key string, data any, ttl time.Duration) {
	cache.mutex.Lock()
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

	cache.mutex.Unlock()
	if !exists && cache.newItemCallback != nil {
		cache.newItemCallback(key, data)
	}
	cache.expirationNotification <- true
}

// Get is a thread-safe way to lookup items
// Every lookup, also touches the cachedItem, hence extending it's life
func (cache *Cache) Get(key string) (any, bool) {
	cache.mutex.Lock()
	item, exists, triggerExpirationNotification := cache.getItem(key)

	var dataToReturn any
	if exists {
		dataToReturn = item.data
	}
	cache.mutex.Unlock()
	if triggerExpirationNotification {
		cache.expirationNotification <- true
	}
	return dataToReturn, exists
}

func (cache *Cache) Remove(key string) bool {
	cache.mutex.Lock()
	object, exists := cache.items[key]
	if !exists {
		cache.mutex.Unlock()
		return false
	}
	delete(cache.items, object.key)
	cache.priorityQueue.remove(object)
	cache.mutex.Unlock()

	return true
}

// Count returns the number of items in the cache
func (cache *Cache) Count() int {
	cache.mutex.Lock()
	length := len(cache.items)
	cache.mutex.Unlock()
	return length
}

func (cache *Cache) SetTTL(ttl time.Duration) {
	cache.mutex.Lock()
	cache.ttl = ttl
	cache.mutex.Unlock()
	cache.expirationNotification <- true
}

// SetExpirationCallback sets a callback that will be called when an cachedItem expires
func (cache *Cache) SetExpirationCallback(callback expireCallback) {
	cache.expireCallback = callback
}

// SetCheckExpirationCallback sets a callback that will be called when an cachedItem is about to expire
// in order to allow external code to decide whether the cachedItem expires or remains for another TTL cycle
func (cache *Cache) SetCheckExpirationCallback(callback checkExpireCallback) {
	cache.checkExpireCallback = callback
}

// SetNewItemCallback sets a callback that will be called when a new cachedItem is added to the cache
func (cache *Cache) SetNewItemCallback(callback expireCallback) {
	cache.newItemCallback = callback
}

// SkipTtlExtensionOnHit allows the user to change the cache behaviour. When this flag is set to true it will
// no longer extend TTL of items when they are retrieved using Get, or when their expiration condition is evaluated
// using SetCheckExpirationCallback.
func (cache *Cache) SkipTtlExtensionOnHit(value bool) {
	cache.skipTTLExtension = value
}

// Purge will remove all entries
func (cache *Cache) Purge() {
	cache.mutex.Lock()
	cache.items = make(map[string]*cachedItem)
	cache.priorityQueue = newPriorityQueue()
	cache.mutex.Unlock()
}

// NewCache is a helper to create instance of the Cache struct
func NewTtlCache() *Cache {

	shutdownChan := make(chan chan struct{})

	cache := &Cache{
		items:                  make(map[string]*cachedItem),
		priorityQueue:          newPriorityQueue(),
		expirationNotification: make(chan bool),
		expirationTime:         time.Now(),
		shutdownSignal:         shutdownChan,
		isShutDown:             false,
	}
	go cache.startExpirationProcessing()
	return cache
}

func min(duration time.Duration, second time.Duration) time.Duration {
	if duration < second {
		return duration
	}
	return second
}

// Load returns key value.
func (cache *Cache) Load(key any) (any, bool) {
	strKey := fmt.Sprintf("%v", key)
	return cache.Get(strKey)
}

// Store sets the key value.
func (cache *Cache) Store(key any, value any) {
	strKey := fmt.Sprintf("%v", key)
	cache.Set(strKey, value)
}

// StoreWithTTL sets the key value with TTL overrides the default.
func (cache *Cache) StoreWithTTL(key any, value any, ttl time.Duration) {
	strKey := fmt.Sprintf("%v", key)
	cache.SetWithTTL(strKey, value, ttl)
}

// Delete deletes the key value.
func (cache *Cache) Delete(key any) {
	strKey := fmt.Sprintf("%v", key)
	cache.Remove(strKey)
}

// Iterate over all items in the cache
func (cache *Cache) Range(cb func(k, v any) bool) {

	for _, val := range cache.items {
		if cb(val.key, val.data) == false {
			return
		}
	}
}

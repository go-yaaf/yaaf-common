// Package cache Cache item implementation
// Based on https://github.com/ReneKroon/ttlcache
package cache

import (
	"time"
)

const (
	// ItemNotExpire is a special TTL value indicating that the item should not expire based on time.
	// However, it can still be removed from the cache through other means, such as manual deletion or callbacks.
	ItemNotExpire time.Duration = -1

	// ItemExpireWithGlobalTTL is a special TTL value indicating that the item should use the cache's global TTL setting.
	ItemExpireWithGlobalTTL time.Duration = 0
)

// newItem creates and initializes a new cachedItem.
// It sets the key, data, and TTL, and also sets the initial expiration time by "touching" the item.
func newItem[K any, T any](key K, data T, ttl time.Duration) *cachedItem[K, T] {
	item := &cachedItem[K, T]{
		key:  key,
		data: data,
		ttl:  ttl,
	}
	// The item is new and not yet in the cache, so it's safe to call touch without a lock.
	item.touch()
	return item
}

// cachedItem represents an item stored in the cache.
// It contains the key, data, TTL, expiration time, and its index in the priority queue.
type cachedItem[K any, T any] struct {
	key        K             // The key of the item.
	data       T             // The data (value) of the item.
	ttl        time.Duration // The time-to-live for this specific item.
	expireAt   time.Time     // The time at which the item will expire.
	queueIndex int           // The index of the item in the priority queue, used for efficient updates.
}

// touch updates the expiration time of the item to the current time plus its TTL.
// This is typically called when an item is accessed or updated.
func (item *cachedItem[K, T]) touch() {
	if item.ttl > 0 {
		item.expireAt = time.Now().Add(item.ttl)
	} else {
		// If TTL is not positive, set a zero expiration time to indicate it doesn't expire.
		item.expireAt = time.Time{}
	}
}

// expired checks if the item has passed its expiration time.
// Items with a non-positive TTL are not considered expired by this check.
func (item *cachedItem[K, T]) expired() bool {
	if item.ttl <= 0 {
		return false
	}
	return item.expireAt.Before(time.Now())
}

/*
 * Cache item implementation
 *
 * Based on https://github.com/ReneKroon/ttlcache
 */
package cache

import (
	"time"
)

const (
	// ItemNotExpire Will avoid the cachedItem being expired by TTL, but can still be exired by callback etc.
	ItemNotExpire time.Duration = -1

	// ItemExpireWithGlobalTTL will use the global TTL when set.
	ItemExpireWithGlobalTTL time.Duration = 0
)

func newItem(key string, data any, ttl time.Duration) *cachedItem {
	item := &cachedItem{
		data: data,
		ttl:  ttl,
		key:  key,
	}
	// since nobody is aware yet of this cachedItem, it's safe to touch without lock here
	item.touch()
	return item
}

type cachedItem struct {
	key        string
	data       any
	ttl        time.Duration
	expireAt   time.Time
	queueIndex int
}

// Reset the cachedItem expiration time
func (item *cachedItem) touch() {
	if item.ttl > 0 {
		item.expireAt = time.Now().Add(item.ttl)
	}
}

// Verify if the cachedItem is expired
func (item *cachedItem) expired() bool {
	if item.ttl <= 0 {
		return false
	}
	return item.expireAt.Before(time.Now())
}

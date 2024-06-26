// Package cache Cache item implementation
// Based on https://github.com/ReneKroon/ttlcache
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

func newItem[K any, T any](key K, data T, ttl time.Duration) *cachedItem[K, T] {
	item := &cachedItem[K, T]{
		data: data,
		ttl:  ttl,
		key:  key,
	}
	// since nobody is aware yet of this cachedItem, it's safe to touch without lock here
	item.touch()
	return item
}

type cachedItem[K any, T any] struct {
	key        K
	data       T
	ttl        time.Duration
	expireAt   time.Time
	queueIndex int
}

// Reset the cachedItem expiration time
func (item *cachedItem[K, T]) touch() {
	if item.ttl > 0 {
		item.expireAt = time.Now().Add(item.ttl)
	}
}

// Verify if the cachedItem is expired
func (item *cachedItem[K, T]) expired() bool {
	if item.ttl <= 0 {
		return false
	}
	return item.expireAt.Before(time.Now())
}

package collections

import (
	"sync"
)

type ConcurrentStringMap[T any] struct {
	sync.RWMutex
	m map[string]T
}

func (c *ConcurrentStringMap[T]) Get(key string) (T, bool) {
	c.RLock()
	val, found := c.m[key]
	c.RUnlock()
	return val, found
}

func (c *ConcurrentStringMap[T]) Put(key string, val T) {
	c.Lock()
	c.m[key] = val
	c.Unlock()
}

func NewConcurrentStringMap[T any]() ConcurrentStringMap[T] {
	return ConcurrentStringMap[T]{
		m: make(map[string]T),
	}
}

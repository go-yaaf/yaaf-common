package collections

import (
	"sync"
)

// ConcurrentStringMap enable safe shared map with Read/Write locks
type ConcurrentStringMap[T any] struct {
	sync.RWMutex
	m map[string]T
}

// NewConcurrentStringMap factory method
func NewConcurrentStringMap[T any]() ConcurrentStringMap[T] {
	return ConcurrentStringMap[T]{
		m: make(map[string]T),
	}
}

// Get retrieve item from map
func (c *ConcurrentStringMap[T]) Get(key string) (T, bool) {
	c.RLock()
	val, found := c.m[key]
	c.RUnlock()
	return val, found
}

// Put set item in the map
func (c *ConcurrentStringMap[T]) Put(key string, val T) {
	c.Lock()
	c.m[key] = val
	c.Unlock()
}

// Keys returns all the keys in the map
func (c *ConcurrentStringMap[T]) Keys() (result []string) {
	c.Lock()
	for k, _ := range c.m {
		result = append(result, k)
	}
	c.Unlock()
	return result
}

// Values returns all the values in the map
func (c *ConcurrentStringMap[T]) Values() (result []T) {
	c.Lock()
	for _, v := range c.m {
		result = append(result, v)
	}
	c.Unlock()
	return result
}

// KeysAndValues returns all the keys and the values in the map
func (c *ConcurrentStringMap[T]) KeysAndValues() (keys []string, values []T) {
	c.Lock()
	for k, v := range c.m {
		keys = append(keys, k)
		values = append(values, v)
	}
	c.Unlock()
	return keys, values
}

// Delete remove item from map
func (c *ConcurrentStringMap[T]) Delete(key string) {
	c.RLock()
	delete(c.m, key)
	c.RUnlock()
}

// DeleteAll remove all items from the map
func (c *ConcurrentStringMap[T]) DeleteAll() {
	c.RLock()
	c.m = make(map[string]T)
	c.RUnlock()
}

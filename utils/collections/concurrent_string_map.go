package collections

import (
	"sync"
)

// ConcurrentStringMap is a thread-safe map with string keys.
// It uses a RWMutex to allow for concurrent read access and exclusive write access.
//
// Type Parameters:
//
//	T: The type of the values stored in the map.
type ConcurrentStringMap[T any] struct {
	sync.RWMutex
	m map[string]T
}

// NewConcurrentStringMap creates and returns a new ConcurrentStringMap.
//
// Type Parameters:
//
//	T: The type of the values to be stored in the map.
//
// Returns:
//
//	An initialized ConcurrentStringMap.
func NewConcurrentStringMap[T any]() *ConcurrentStringMap[T] {
	return &ConcurrentStringMap[T]{
		m: make(map[string]T),
	}
}

// Get retrieves a value from the map by its key.
// It is safe for concurrent use.
//
// Parameters:
//
//	key: The key of the value to retrieve.
//
// Returns:
//
//	The value associated with the key, and a boolean indicating if the key was found.
func (c *ConcurrentStringMap[T]) Get(key string) (T, bool) {
	c.RLock()
	defer c.RUnlock()
	val, found := c.m[key]
	return val, found
}

// Put adds or updates a value in the map.
// It is safe for concurrent use.
//
// Parameters:
//
//	key: The key of the value to set.
//	val: The value to set.
func (c *ConcurrentStringMap[T]) Put(key string, val T) {
	c.Lock()
	defer c.Unlock()
	c.m[key] = val
}

// Keys returns a slice of all keys in the map.
// The order of the keys is not guaranteed.
// It is safe for concurrent use.
//
// Returns:
//
//	A slice of strings containing all the keys in the map.
func (c *ConcurrentStringMap[T]) Keys() []string {
	c.RLock()
	defer c.RUnlock()
	keys := make([]string, 0, len(c.m))
	for k := range c.m {
		keys = append(keys, k)
	}
	return keys
}

// Values returns a slice of all values in the map.
// The order of the values is not guaranteed.
// It is safe for concurrent use.
//
// Returns:
//
//	A slice containing all the values in the map.
func (c *ConcurrentStringMap[T]) Values() []T {
	c.RLock()
	defer c.RUnlock()
	values := make([]T, 0, len(c.m))
	for _, v := range c.m {
		values = append(values, v)
	}
	return values
}

// KeysAndValues returns two slices: one with all the keys and one with all the values.
// The order of keys and values is not guaranteed, but the correspondence between a key and a value at the same index is maintained.
// It is safe for concurrent use.
//
// Returns:
//
//	A slice of keys and a slice of values.
func (c *ConcurrentStringMap[T]) KeysAndValues() ([]string, []T) {
	c.RLock()
	defer c.RUnlock()
	keys := make([]string, 0, len(c.m))
	values := make([]T, 0, len(c.m))
	for k, v := range c.m {
		keys = append(keys, k)
		values = append(values, v)
	}
	return keys, values
}

// Delete removes a key-value pair from the map.
// It is safe for concurrent use.
//
// Parameters:
//
//	key: The key of the item to delete.
func (c *ConcurrentStringMap[T]) Delete(key string) {
	c.Lock()
	defer c.Unlock()
	delete(c.m, key)
}

// DeleteAll removes all key-value pairs from the map.
// It is safe for concurrent use.
func (c *ConcurrentStringMap[T]) DeleteAll() {
	c.Lock()
	defer c.Unlock()
	c.m = make(map[string]T)
}

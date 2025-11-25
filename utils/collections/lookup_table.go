package collections

// LookupTable is a generic lookup table, essentially a map with string keys and values of a specified type.
// It provides a set of convenient methods for common map operations.
//
// Type Parameters:
//
//	T: The type of the values stored in the lookup table.
type LookupTable[T any] map[string]T

// Get retrieves the value associated with the given key.
//
// Parameters:
//
//	key: The key to look up.
//
// Returns:
//
//	The value associated with the key, and a boolean indicating if the key was found.
func (lt LookupTable[T]) Get(key string) (T, bool) {
	val, ok := lt[key]
	return val, ok
}

// Keys returns a slice containing all the keys in the lookup table.
// The order of the keys is not guaranteed.
//
// Returns:
//
//	A slice of strings representing the keys.
func (lt LookupTable[T]) Keys() []string {
	keys := make([]string, 0, len(lt))
	for k := range lt {
		keys = append(keys, k)
	}
	return keys
}

// Values returns a slice containing all the values in the lookup table.
// The order of the values is not guaranteed.
//
// Returns:
//
//	A slice containing the values.
func (lt LookupTable[T]) Values() []T {
	vals := make([]T, 0, len(lt))
	for _, v := range lt {
		vals = append(vals, v)
	}
	return vals
}

// Delete removes the entry for the given key from the lookup table.
// If the key does not exist, this is a no-op.
//
// Parameters:
//
//	key: The key of the entry to delete.
func (lt LookupTable[T]) Delete(key string) {
	delete(lt, key)
}

// Len returns the number of entries in the lookup table.
//
// Returns:
//
//	The number of entries.
func (lt LookupTable[T]) Len() int {
	return len(lt)
}

// Clear removes all entries from the lookup table, making it empty.
func (lt LookupTable[T]) Clear() {
	for k := range lt {
		delete(lt, k)
	}
}

// Contains checks if the given key exists in the lookup table.
//
// Parameters:
//
//	key: The key to check for.
//
// Returns:
//
//	`true` if the key exists, `false` otherwise.
func (lt LookupTable[T]) Contains(key string) bool {
	_, ok := lt[key]
	return ok
}

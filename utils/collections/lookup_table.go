package collections

// LookupTable is a generic lookup table with string keys and values of type T
type LookupTable[T any] map[string]T

// Get returns the value and a bool indicating if it was found
func (lt LookupTable[T]) Get(key string) (T, bool) {
	val, ok := lt[key]
	return val, ok
}

// Keys returns all keys in the table
func (lt LookupTable[T]) Keys() []string {
	keys := make([]string, 0, len(lt))
	for k := range lt {
		keys = append(keys, k)
	}
	return keys
}

// Values returns all values in the table
func (lt LookupTable[T]) Values() []T {
	vals := make([]T, 0, len(lt))
	for _, v := range lt {
		vals = append(vals, v)
	}
	return vals
}

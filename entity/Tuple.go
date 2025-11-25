package entity

// Tuple represents a generic key-value pair.
// It is useful for storing associated data where a full map is not needed or order matters.
//
// Type Parameters:
//   - K: The type of the key.
//   - V: The type of the value.
type Tuple[K, V any] struct {
	Key   K `json:"key"`   // Key is the first element of the tuple
	Value V `json:"value"` // Value is the second element of the tuple
}

package entity

// Tuple model represents a generic key-value pair
type Tuple[K, V any] struct {
	Key   K `json:"key"`   // Tuple key
	Value V `json:"value"` // Tuple value
}

package collections

// MakeSet creates a set from a slice of items. A set is represented as a map
// with the items as keys and empty structs as values, which is a common and
// memory-efficient way to implement sets in Go.
//
// Type Parameters:
//
//	T: The type of the items, which must be comparable.
//
// Parameters:
//
//	items: A slice of items to be converted into a set.
//
// Returns:
//
//	A map representing the set of items.
func MakeSet[T comparable](items []T) map[T]struct{} {
	set := make(map[T]struct{}, len(items))
	for _, item := range items {
		set[item] = struct{}{}
	}
	return set
}

// InSet checks if an item is present in a set.
//
// Type Parameters:
//
//	T: The type of the item, which must be comparable.
//
// Parameters:
//
//	set: The set to check for membership in.
//	item: The item to check for.
//
// Returns:
//
//	`true` if the item is in the set, `false` otherwise.
func InSet[T comparable](set map[T]struct{}, item T) bool {
	_, exists := set[item]
	return exists
}

// InSetAny checks if any of the items in a slice are present in a set.
// It returns true as soon as the first match is found.
//
// Type Parameters:
//
//	T: The type of the items, which must be comparable.
//
// Parameters:
//
//	haystack: The set to search within.
//	needles: A slice of items to search for.
//
// Returns:
//
//	`true` if any of the `needles` are found in the `haystack`, `false` otherwise.
func InSetAny[T comparable](haystack map[T]struct{}, needles []T) bool {
	for _, needle := range needles {
		if _, found := haystack[needle]; found {
			return true
		}
	}
	return false
}

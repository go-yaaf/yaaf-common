package collections

// Generic MakeSet function that builds a set from a slice
func MakeSet[T comparable](items []T) map[T]struct{} {
	set := make(map[T]struct{}, len(items))
	for _, item := range items {
		set[item] = struct{}{}
	}
	return set
}

// Generic InSet function that checks membership in a set
func InSet[T comparable](set map[T]struct{}, item T) bool {
	_, exists := set[item]
	return exists
}

// InSetAny returns true if any item in 'needles' is found in 'haystack'
func InSetAny[T comparable](haystack map[T]struct{}, needles []T) bool {
	for _, needle := range needles {
		if _, found := haystack[needle]; found {
			return true
		}
	}
	return false
}

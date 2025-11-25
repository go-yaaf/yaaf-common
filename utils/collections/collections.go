// Package collections provides a set of utility functions for working with slices and collections.
// These functions are designed to be generic and reusable, simplifying common operations
// such as filtering, mapping, and searching.

package collections

import (
	"fmt"
	"strings"
)

// Index returns the first index of the target value `t` in the slice `vs`, or -1 if no match is found.
//
// Type Parameters:
//
//	T: The type of the elements in the slice, which must be comparable.
//
// Parameters:
//
//	vs: The slice to search in.
//	t: The value to search for.
//
// Returns:
//
//	The index of the first occurrence of `t`, or -1 if not found.
func Index[T comparable](vs []T, t T) int {
	for i, v := range vs {
		if v == t {
			return i
		}
	}
	return -1
}

// Include returns `true` if the target value `t` is in the slice `vs`.
//
// Type Parameters:
//
//	T: The type of the elements in the slice, which must be comparable.
//
// Parameters:
//
//	vs: The slice to search in.
//	t: The value to search for.
//
// Returns:
//
//	`true` if `t` is found in `vs`, `false` otherwise.
func Include[T comparable](vs []T, t T) bool {
	return Index(vs, t) >= 0
}

// IncludeMask checks if any integer in the slice `vs` has all the bits of the mask `t` set.
//
// Parameters:
//
//	vs: The slice of integers to check.
//	t: The bitmask to check against.
//
// Returns:
//
//	`true` if any integer in `vs` includes the mask `t`, `false` otherwise.
func IncludeMask(vs []int, t int) bool {
	for _, r := range vs {
		if (r & t) == r {
			return true
		}
	}
	return false
}

// Any returns `true` if at least one element in the slice `vs` satisfies the predicate `f`.
//
// Type Parameters:
//
//	T: The type of the elements in the slice.
//
// Parameters:
//
//	vs: The slice to check.
//	f: The predicate function.
//
// Returns:
//
//	`true` if any element satisfies `f`, `false` otherwise.
func Any[T any](vs []T, f func(T) bool) bool {
	for _, v := range vs {
		if f(v) {
			return true
		}
	}
	return false
}

// All returns `true` if all elements in the slice `vs` satisfy the predicate `f`.
//
// Type Parameters:
//
//	T: The type of the elements in the slice.
//
// Parameters:
//
//	vs: The slice to check.
//	f: The predicate function.
//
// Returns:
//
//	`true` if all elements satisfy `f`, `false` otherwise.
func All[T any](vs []T, f func(T) bool) bool {
	for _, v := range vs {
		if !f(v) {
			return false
		}
	}
	return true
}

// Filter returns a new slice containing all elements from `vs` that satisfy the predicate `f`.
//
// Type Parameters:
//
//	T: The type of the elements in the slice.
//
// Parameters:
//
//	vs: The slice to filter.
//	f: The predicate function.
//
// Returns:
//
//	A new slice with the filtered elements.
func Filter[T any](vs []T, f func(T) bool) []T {
	vsf := make([]T, 0)
	for _, v := range vs {
		if f(v) {
			vsf = append(vsf, v)
		}
	}
	return vsf
}

// Map returns a new slice containing the results of applying the function `f` to each element in the original slice `vs`.
//
// Type Parameters:
//
//	T: The type of the elements in the input slice.
//	U: The type of the elements in the output slice.
//
// Parameters:
//
//	vs: The slice to map.
//	f: The mapping function.
//
// Returns:
//
//	A new slice with the mapped elements.
func Map[T any, U any](vs []T, f func(T) U) []U {
	vsm := make([]U, len(vs))
	for i, v := range vs {
		vsm[i] = f(v)
	}
	return vsm
}

// Distinct returns a new slice with all duplicate values removed.
//
// Type Parameters:
//
//	T: The type of the elements in the slice, which must be comparable.
//
// Parameters:
//
//	vs: The slice to remove duplicates from.
//
// Returns:
//
//	A new slice with unique elements.
func Distinct[T comparable](vs []T) []T {
	set := make(map[T]struct{})
	result := make([]T, 0)
	for _, s := range vs {
		if _, ok := set[s]; !ok {
			set[s] = struct{}{}
			result = append(result, s)
		}
	}
	return result
}

// Concat concatenates multiple slices into a single slice.
//
// Type Parameters:
//
//	T: The type of the elements in the slices.
//
// Parameters:
//
//	slices: The slices to concatenate.
//
// Returns:
//
//	A new slice containing all elements from the input slices.
func Concat[T any](slices ...[]T) []T {
	total := 0
	for _, slc := range slices {
		total += len(slc)
	}

	result := make([]T, total)
	var i int
	for _, s := range slices {
		i += copy(result[i:], s)
	}
	return result
}

// Intersect returns a new slice containing only the elements that exist in all given slices.
//
// Type Parameters:
//
//	T: The type of the elements in the slices, which must be comparable.
//
// Parameters:
//
//	slices: The slices to find the intersection of.
//
// Returns:
//
//	A new slice with the common elements.
func Intersect[T comparable](slices ...[]T) []T {
	if len(slices) == 0 {
		return []T{}
	}

	set := make(map[T]struct{})
	for _, item := range slices[0] {
		set[item] = struct{}{}
	}

	for _, slc := range slices[1:] {
		nextSet := make(map[T]struct{})
		for _, item := range slc {
			if _, ok := set[item]; ok {
				nextSet[item] = struct{}{}
			}
		}
		set = nextSet
	}

	result := make([]T, 0, len(set))
	for item := range set {
		result = append(result, item)
	}
	return result
}

// AddIfNotExists adds an element to the slice if it is not already present.
//
// Type Parameters:
//
//	T: The type of the elements in the slice, which must be comparable.
//
// Parameters:
//
//	vs: The slice to add to.
//	t: The element to add.
//
// Returns:
//
//	A new slice with the element added if it was not present, or the original slice.
func AddIfNotExists[T comparable](vs []T, t T) []T {
	if Index(vs, t) >= 0 {
		return vs
	}
	return append(vs, t)
}

// Remove removes the first occurrence of an element from the slice.
//
// Type Parameters:
//
//	T: The type of the elements in the slice, which must be comparable.
//
// Parameters:
//
//	vs: The slice to remove from.
//	t: The element to remove.
//
// Returns:
//
//	A new slice with the element removed.
func Remove[T comparable](vs []T, t T) []T {
	for i, v := range vs {
		if v == t {
			return append(vs[:i], vs[i+1:]...)
		}
	}
	return vs
}

// BitMaskInclude checks if a source bitmask includes a given flag.
//
// Parameters:
//
//	src: The source bitmask.
//	flag: The flag to check for.
//
// Returns:
//
//	`true` if the source bitmask includes the flag, `false` otherwise.
func BitMaskInclude(src, flag int) bool {
	return src&flag == flag
}

// JoinN converts a slice of integers to a string, with elements separated by a given separator.
//
// Parameters:
//
//	slice: The slice of integers to join.
//	sep: The separator string.
//
// Returns:
//
//	A string representation of the integer slice.
func JoinN(slice []int, sep string) string {
	list := make([]string, len(slice))
	for i, v := range slice {
		list[i] = fmt.Sprintf("%d", v)
	}
	return strings.Join(list, sep)
}

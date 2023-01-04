// We often need our programs to perform operations on
// collections of data, like selecting all items that
// satisfy a given predicate or mapping all items to a new
// collection with a custom function.
//

package collections

import (
	"fmt"
	"strings"
)

// Index returns the first index of the target string `t`, or -1 if no match is found.
func Index(vs []string, t string) int {
	for i, v := range vs {
		if v == t {
			return i
		}
	}
	return -1
}

// IndexN returns the first index of the target int `t`, or -1 if no match is found.
func IndexN(vs []int, t int) int {
	for i, v := range vs {
		if v == t {
			return i
		}
	}
	return -1
}

// Include returns `true` if the target string t is in the slice.
func Include(vs []string, t string) bool {
	return Index(vs, t) >= 0
}

// IncludeN returns `true` if the target int t is in the slice.
func IncludeN(vs []int, t int) bool {
	return IndexN(vs, t) >= 0
}

func IncludeMask(vs []int, t int) bool {
	for _, r := range vs {
		if (r & t) == r {
			return true
		}
	}
	return false
}

// Any returns `true` if one of the strings in the slice satisfies the predicate `f`.
func Any(vs []string, f func(string) bool) bool {
	for _, v := range vs {
		if f(v) {
			return true
		}
	}
	return false
}

// AnyN returns `true` if one of the integers in the slice satisfies the predicate `f`.
func AnyN(vs []int, f func(int) bool) bool {
	for _, v := range vs {
		if f(v) {
			return true
		}
	}
	return false
}

// All returns `true` if all the strings in the slice satisfy the predicate `f`.
func All(vs []string, f func(string) bool) bool {
	for _, v := range vs {
		if !f(v) {
			return false
		}
	}
	return true
}

// AllN returns `true` if all the integers in the slice satisfy the predicate `f`.
func AllN(vs []int, f func(int) bool) bool {
	for _, v := range vs {
		if !f(v) {
			return false
		}
	}
	return true
}

// Filter returns a new slice containing all strings in the slice that satisfy the predicate `f`.
func Filter(vs []string, f func(string) bool) []string {
	vsf := make([]string, 0)
	for _, v := range vs {
		if f(v) {
			vsf = append(vsf, v)
		}
	}
	return vsf
}

// FilterN returns a new slice containing all integers in the slice that satisfy the predicate `f`.
func FilterN(vs []int, f func(int) bool) []int {
	vsf := make([]int, 0)
	for _, v := range vs {
		if f(v) {
			vsf = append(vsf, v)
		}
	}
	return vsf
}

// Map returns a new slice containing the results of applying the function `f` to each string in the original slice.
func Map(vs []string, f func(string) string) []string {
	vsm := make([]string, len(vs))
	for i, v := range vs {
		vsm[i] = f(v)
	}
	return vsm
}

// MapN returns a new slice containing the results of applying the function `f` to each integers in the original slice.
func MapN(vs []int, f func(int) int) []int {
	vsm := make([]int, len(vs))
	for i, v := range vs {
		vsm[i] = f(v)
	}
	return vsm
}

// Distinct returns a new slice without duplications
func Distinct(vs []string) []string {

	strMap := make(map[string]string)
	for _, s := range vs {
		strMap[s] = s
	}

	vsm := make([]string, 0)
	for k, _ := range strMap {
		vsm = append(vsm, k)
	}
	return vsm
}

// DistinctN returns a new slice without duplications
func DistinctN(vs []int) []int {

	strMap := make(map[int]int)
	for _, s := range vs {
		strMap[s] = s
	}

	vsm := make([]int, 0)
	for k, _ := range strMap {
		vsm = append(vsm, k)
	}
	return vsm
}

// Concatenate multiple string slices efficiently
func Concat(slices ...[]string) []string {

	total := 0
	for _, slc := range slices {
		total += len(slc)
	}

	result := make([]string, total)

	var i int
	for _, s := range slices {
		i += copy(result[i:], s)
	}
	return result
}

// Concatenate multiple int slices efficiently
func ConcatN(slices ...[]int) []int {

	total := 0
	for _, slc := range slices {
		total += len(slc)
	}

	result := make([]int, total)

	var i int
	for _, s := range slices {
		i += copy(result[i:], s)
	}
	return result
}

// Intersect returns only the string values in all slices
func Intersect(slices ...[]string) []string {

	result := make([]string, 0)

	for _, v := range slices[0] {
		existsInAll := true
		for _, slc := range slices {
			existsInAll = existsInAll && Include(slc, v)
		}
		if existsInAll {
			result = append(result, v)
		}
	}
	return result
}

// Intersect returns only the values in all slices
func IntersectN(slices ...[]int) []int {

	result := make([]int, 0)

	for _, v := range slices[0] {
		existsInAll := true
		for _, slc := range slices {
			existsInAll = existsInAll && IncludeN(slc, v)
		}
		if existsInAll {
			result = append(result, v)
		}
	}
	return result
}

// Add item to list only if it does not exist
func AddIfNotExists(vs []string, t string) []string {
	if Index(vs, t) >= 0 {
		return vs
	} else {
		return append(vs, t)
	}
}

// Remove an item from the array
func Remove(vs []string, t string) []string {

	for i, v := range vs {
		if v == t {
			vs = append(vs[:i], vs[i+1:]...)
			break
		}
	}
	return vs
}

// BitMaskInclude checks if the src bitmask including the flag
func BitMaskInclude(src, flag int) bool {
	return src&flag == flag
}

// JoinN convert all integers in the slice to strings and joins them together as a single string with separator
func JoinN(slice []int, sep string) string {

	list := make([]string, 0)
	for _, v := range slice {
		list = append(list, fmt.Sprintf("%d", v))
	}
	return strings.Join(list, sep)
}

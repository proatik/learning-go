package freq

import (
	"cmp"
	"slices"
)

// Frequencies counts how many times each value appears in xs.
func Frequencies[T comparable](xs []T) map[T]int {
	counts := make(map[T]int)
	for _, x := range xs {
		counts[x]++
	}
	return counts
}

// SortedKeys returns the keys of m in ascending order.
func SortedKeys[K cmp.Ordered, V any](m map[K]V) []K {
	keys := make([]K, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	slices.Sort(keys)
	return keys
}

// MostCommon returns the most frequent value in xs and its count.
// Ties go to the smallest value. ok is false when xs is empty.
func MostCommon[T cmp.Ordered](xs []T) (value T, count int, ok bool) {
	counts := Frequencies(xs)
	for _, k := range SortedKeys(counts) {
		if counts[k] > count {
			value, count, ok = k, counts[k], true
		}
	}
	return value, count, ok
}

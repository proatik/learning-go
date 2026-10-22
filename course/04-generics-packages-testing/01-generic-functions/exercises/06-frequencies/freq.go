package freq

import "cmp"

// Frequencies counts how many times each value appears in xs.
func Frequencies[T comparable](xs []T) map[T]int {
	// TODO
	return nil
}

// SortedKeys returns the keys of m in ascending order.
func SortedKeys[K cmp.Ordered, V any](m map[K]V) []K {
	// TODO
	return nil
}

// MostCommon returns the most frequent value in xs and its count.
// Ties go to the smallest value. ok is false when xs is empty.
func MostCommon[T cmp.Ordered](xs []T) (value T, count int, ok bool) {
	// TODO: combine Frequencies and SortedKeys
	return value, count, ok
}

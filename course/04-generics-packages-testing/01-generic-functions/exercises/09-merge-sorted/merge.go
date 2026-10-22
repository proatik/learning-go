package merge

import "cmp"

// Merge merges a and b, both sorted by compare, into a new sorted slice.
// On ties the element from a comes first.
func Merge[T any](a, b []T, compare func(x, y T) int) []T {
	// TODO: two indexes i (into a) and j (into b)
	return nil
}

// MergeOrdered merges two ascending slices of an ordered type.
func MergeOrdered[T cmp.Ordered](a, b []T) []T {
	// TODO: reuse Merge
	return nil
}

// MergeAll merges any number of lists sorted by compare.
// On ties, elements of earlier lists come first.
func MergeAll[T any](compare func(x, y T) int, lists ...[]T) []T {
	// TODO
	return nil
}

// Desc returns a comparator that orders the opposite way to compare.
func Desc[T any](compare func(x, y T) int) func(x, y T) int {
	// TODO
	return func(x, y T) int { return 0 }
}

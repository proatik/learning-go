package merge

import "cmp"

// Merge merges a and b, both sorted by compare, into a new sorted slice.
// On ties the element from a comes first.
func Merge[T any](a, b []T, compare func(x, y T) int) []T {
	out := make([]T, 0, len(a)+len(b))
	i, j := 0, 0
	for i < len(a) && j < len(b) {
		// <= 0 takes a on ties, which keeps the merge stable.
		if compare(a[i], b[j]) <= 0 {
			out = append(out, a[i])
			i++
		} else {
			out = append(out, b[j])
			j++
		}
	}
	out = append(out, a[i:]...)
	return append(out, b[j:]...)
}

// MergeOrdered merges two ascending slices of an ordered type.
func MergeOrdered[T cmp.Ordered](a, b []T) []T {
	return Merge(a, b, cmp.Compare[T])
}

// MergeAll merges any number of lists sorted by compare.
// On ties, elements of earlier lists come first.
func MergeAll[T any](compare func(x, y T) int, lists ...[]T) []T {
	var out []T
	for _, l := range lists {
		out = Merge(out, l, compare)
	}
	return out
}

// Desc returns a comparator that orders the opposite way to compare.
func Desc[T any](compare func(x, y T) int) func(x, y T) int {
	return func(x, y T) int { return compare(y, x) }
}

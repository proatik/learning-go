package psort

import "cmp"

// SortFunc sorts s in place by cmp, sorting the two halves concurrently down
// to depth levels and sequentially below that.
func SortFunc[T any](s []T, cmp func(a, b T) int, depth int) {
	// TODO
}

// Sort sorts s in place in ascending order, like SortFunc with cmp.Compare.
func Sort[T cmp.Ordered](s []T, depth int) {
	// TODO
}

package psort

import (
	"cmp"
	"slices"
	"sync"
)

// SortFunc sorts s in place by cmp, sorting the two halves concurrently down
// to depth levels and sequentially below that.
func SortFunc[T any](s []T, cmp func(a, b T) int, depth int) {
	if depth <= 0 || len(s) < 2 {
		slices.SortFunc(s, cmp)
		return
	}
	mid := len(s) / 2
	var wg sync.WaitGroup
	wg.Go(func() { SortFunc(s[:mid], cmp, depth-1) })
	SortFunc(s[mid:], cmp, depth-1) // the current goroutine takes the other half
	wg.Wait()
	merge(s, mid, cmp)
}

// merge merges the sorted runs s[:mid] and s[mid:] into s.
func merge[T any](s []T, mid int, cmp func(a, b T) int) {
	out := make([]T, 0, len(s))
	i, j := 0, mid
	for i < mid && j < len(s) {
		if cmp(s[j], s[i]) < 0 {
			out = append(out, s[j])
			j++
		} else {
			out = append(out, s[i])
			i++
		}
	}
	out = append(out, s[i:mid]...)
	out = append(out, s[j:]...)
	copy(s, out)
}

// Sort sorts s in place in ascending order, like SortFunc with cmp.Compare.
func Sort[T cmp.Ordered](s []T, depth int) {
	SortFunc(s, cmp.Compare[T], depth)
}

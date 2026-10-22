package search

import (
	"slices"
	"sync"
)

// FirstMatch calls pred on every item concurrently and returns the smallest
// index whose item matches, or -1 if none do.
func FirstMatch[T any](items []T, pred func(T) bool) int {
	matched := make([]bool, len(items))
	var wg sync.WaitGroup
	for i, item := range items {
		wg.Go(func() { matched[i] = pred(item) })
	}
	wg.Wait()
	return slices.Index(matched, true)
}

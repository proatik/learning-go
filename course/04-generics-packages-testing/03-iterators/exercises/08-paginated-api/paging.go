package paging

import "iter"

// Fetch loads one page of results; pages are numbered from 1.
// more reports whether a later page exists.
type Fetch[T any] func(page int) (items []T, more bool, err error)

// Items returns an iterator over every item of every page, fetching pages
// lazily. A fetch error is yielded once as (zero value, err) and ends the
// sequence.
func Items[T any](fetch Fetch[T]) iter.Seq2[T, error] {
	return func(yield func(T, error) bool) {
		// TODO: loop over page numbers from 1; fetch, handle the error,
		// yield each item, stop when there are no more pages
	}
}

// CollectAll returns every item of seq, or the items before the first
// error together with that error.
func CollectAll[T any](seq iter.Seq2[T, error]) ([]T, error) {
	// TODO
	return nil, nil
}

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
		for page := 1; ; page++ {
			items, more, err := fetch(page)
			if err != nil {
				var zero T
				yield(zero, err)
				return
			}
			for _, item := range items {
				if !yield(item, nil) {
					return
				}
			}
			if !more {
				return
			}
		}
	}
}

// CollectAll returns every item of seq, or the items before the first
// error together with that error.
func CollectAll[T any](seq iter.Seq2[T, error]) ([]T, error) {
	var all []T
	for item, err := range seq {
		if err != nil {
			return all, err
		}
		all = append(all, item)
	}
	return all, nil
}

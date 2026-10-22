// Package pipe provides lazy adapters over iter.Seq.
package pipe

import "iter"

// Map yields f(v) for every v in seq.
func Map[T, U any](seq iter.Seq[T], f func(T) U) iter.Seq[U] {
	return func(yield func(U) bool) {
		for v := range seq {
			if !yield(f(v)) {
				return
			}
		}
	}
}

// Filter yields the values of seq for which keep returns true.
func Filter[T any](seq iter.Seq[T], keep func(T) bool) iter.Seq[T] {
	return func(yield func(T) bool) {
		for v := range seq {
			if keep(v) && !yield(v) {
				return
			}
		}
	}
}

// Take yields at most the first n values of seq and pulls no more than that.
func Take[T any](seq iter.Seq[T], n int) iter.Seq[T] {
	return func(yield func(T) bool) {
		if n <= 0 {
			return
		}
		count := 0
		for v := range seq {
			if !yield(v) {
				return
			}
			count++
			if count == n {
				return
			}
		}
	}
}

// Batch yields the values of seq in fresh slices of size values; the last
// batch may be shorter. It yields nothing if size < 1.
func Batch[T any](seq iter.Seq[T], size int) iter.Seq[[]T] {
	return func(yield func([]T) bool) {
		if size < 1 {
			return
		}
		var batch []T
		for v := range seq {
			batch = append(batch, v)
			if len(batch) == size {
				if !yield(batch) {
					return
				}
				batch = nil
			}
		}
		if len(batch) > 0 {
			yield(batch)
		}
	}
}

// Fold combines the values of seq into one result, starting from init.
func Fold[T, A any](seq iter.Seq[T], init A, f func(A, T) A) A {
	acc := init
	for v := range seq {
		acc = f(acc, v)
	}
	return acc
}

package window

import "iter"

// Windows returns an iterator over every sub-slice of s of length k,
// from left to right. Each window shares memory with s and has its
// capacity clipped to k. It yields nothing if k < 1 or k > len(s).
func Windows[T any](s []T, k int) iter.Seq[[]T] {
	return func(yield func([]T) bool) {
		if k < 1 {
			return
		}
		for i := 0; i+k <= len(s); i++ {
			if !yield(s[i : i+k : i+k]) {
				return
			}
		}
	}
}

// Dedupe returns an iterator over seq with consecutive duplicate values
// removed.
func Dedupe[V comparable](seq iter.Seq[V]) iter.Seq[V] {
	return func(yield func(V) bool) {
		var prev V
		first := true
		for v := range seq {
			if !first && v == prev {
				continue
			}
			first = false
			prev = v
			if !yield(v) {
				return
			}
		}
	}
}

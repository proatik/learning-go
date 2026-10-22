package combine

import "iter"

// Enumerate returns an iterator over (index, value) pairs of seq,
// counting from 0.
func Enumerate[V any](seq iter.Seq[V]) iter.Seq2[int, V] {
	return func(yield func(int, V) bool) {
		i := 0
		for v := range seq {
			if !yield(i, v) {
				return
			}
			i++
		}
	}
}

// Chain returns an iterator over every value of each sequence in seqs,
// in order. Nil sequences are skipped.
func Chain[V any](seqs ...iter.Seq[V]) iter.Seq[V] {
	return func(yield func(V) bool) {
		for _, seq := range seqs {
			if seq == nil {
				continue
			}
			for v := range seq {
				if !yield(v) {
					return
				}
			}
		}
	}
}

// Keys returns an iterator over the first value of each pair in seq.
func Keys[K, V any](seq iter.Seq2[K, V]) iter.Seq[K] {
	return func(yield func(K) bool) {
		for k := range seq {
			if !yield(k) {
				return
			}
		}
	}
}

// Values returns an iterator over the second value of each pair in seq.
func Values[K, V any](seq iter.Seq2[K, V]) iter.Seq[V] {
	return func(yield func(V) bool) {
		for _, v := range seq {
			if !yield(v) {
				return
			}
		}
	}
}

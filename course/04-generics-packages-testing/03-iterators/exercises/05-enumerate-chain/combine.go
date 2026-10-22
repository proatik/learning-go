package combine

import "iter"

// Enumerate returns an iterator over (index, value) pairs of seq,
// counting from 0.
func Enumerate[V any](seq iter.Seq[V]) iter.Seq2[int, V] {
	return func(yield func(int, V) bool) {
		// TODO: keep a counter next to a range loop over seq
	}
}

// Chain returns an iterator over every value of each sequence in seqs,
// in order. Nil sequences are skipped.
func Chain[V any](seqs ...iter.Seq[V]) iter.Seq[V] {
	return func(yield func(V) bool) {
		// TODO: one loop over seqs, one range loop over each seq
	}
}

// Keys returns an iterator over the first value of each pair in seq.
func Keys[K, V any](seq iter.Seq2[K, V]) iter.Seq[K] {
	return func(yield func(K) bool) {
		// TODO: for k := range seq
	}
}

// Values returns an iterator over the second value of each pair in seq.
func Values[K, V any](seq iter.Seq2[K, V]) iter.Seq[V] {
	return func(yield func(V) bool) {
		// TODO: for _, v := range seq
	}
}

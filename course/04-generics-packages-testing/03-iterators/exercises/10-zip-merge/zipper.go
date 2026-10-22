package zipper

import (
	"cmp"
	"iter"
)

// Zip returns an iterator over pairs (a, b) taken from as and bs in step.
// It stops as soon as either sequence ends.
func Zip[A, B any](as iter.Seq[A], bs iter.Seq[B]) iter.Seq2[A, B] {
	return func(yield func(A, B) bool) {
		// TODO: range over as, pull from bs with iter.Pull (defer stop!)
	}
}

// Merge returns an iterator over the values of the sorted sequences a and b,
// in ascending order. On a tie, the value from a comes first.
func Merge[V cmp.Ordered](a, b iter.Seq[V]) iter.Seq[V] {
	return func(yield func(V) bool) {
		// TODO: pull from both; yield the smaller head, then advance that side
	}
}

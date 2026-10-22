package zipper

import (
	"cmp"
	"iter"
)

// Zip returns an iterator over pairs (a, b) taken from as and bs in step.
// It stops as soon as either sequence ends.
func Zip[A, B any](as iter.Seq[A], bs iter.Seq[B]) iter.Seq2[A, B] {
	return func(yield func(A, B) bool) {
		nextB, stop := iter.Pull(bs)
		defer stop()
		for a := range as {
			b, ok := nextB()
			if !ok || !yield(a, b) {
				return
			}
		}
	}
}

// Merge returns an iterator over the values of the sorted sequences a and b,
// in ascending order. On a tie, the value from a comes first.
func Merge[V cmp.Ordered](a, b iter.Seq[V]) iter.Seq[V] {
	return func(yield func(V) bool) {
		nextA, stopA := iter.Pull(a)
		defer stopA()
		nextB, stopB := iter.Pull(b)
		defer stopB()

		va, okA := nextA()
		vb, okB := nextB()
		for okA && okB {
			if va <= vb {
				if !yield(va) {
					return
				}
				va, okA = nextA()
			} else {
				if !yield(vb) {
					return
				}
				vb, okB = nextB()
			}
		}
		for ; okA; va, okA = nextA() {
			if !yield(va) {
				return
			}
		}
		for ; okB; vb, okB = nextB() {
			if !yield(vb) {
				return
			}
		}
	}
}

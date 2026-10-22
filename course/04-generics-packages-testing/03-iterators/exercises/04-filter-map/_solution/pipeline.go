package pipeline

import "iter"

// Filter returns an iterator over the values of seq for which keep returns true.
func Filter[V any](seq iter.Seq[V], keep func(V) bool) iter.Seq[V] {
	return func(yield func(V) bool) {
		for v := range seq {
			if keep(v) && !yield(v) {
				return
			}
		}
	}
}

// Map returns an iterator over f(v) for every value v of seq.
func Map[In, Out any](seq iter.Seq[In], f func(In) Out) iter.Seq[Out] {
	return func(yield func(Out) bool) {
		for v := range seq {
			if !yield(f(v)) {
				return
			}
		}
	}
}

// Reduce folds seq into a single value: acc starts at init and becomes
// f(acc, v) for every value v. It returns init for an empty sequence.
func Reduce[V, A any](seq iter.Seq[V], init A, f func(A, V) A) A {
	acc := init
	for v := range seq {
		acc = f(acc, v)
	}
	return acc
}

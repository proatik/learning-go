package countdown

import "iter"

// Countdown returns an iterator over n, n-1, …, 1.
// For n <= 0 the sequence is empty.
func Countdown(n int) iter.Seq[int] {
	return func(yield func(int) bool) {
		// TODO: yield n, n-1, …, 1; return as soon as yield returns false
	}
}

// Repeat returns an iterator that yields v exactly n times.
// For n <= 0 the sequence is empty.
func Repeat[V any](v V, n int) iter.Seq[V] {
	return func(yield func(V) bool) {
		// TODO: yield v n times; return as soon as yield returns false
	}
}

package fibo

import "iter"

// Fibonacci returns an infinite iterator over the Fibonacci numbers
// 0, 1, 1, 2, 3, 5, …
func Fibonacci() iter.Seq[int] {
	return func(yield func(int) bool) {
		// TODO: keep two numbers a, b; yield a, then move to b, a+b
	}
}

// Take returns an iterator over the first n values of seq (fewer if seq
// is shorter). It never pulls more than n values from seq.
func Take[V any](seq iter.Seq[V], n int) iter.Seq[V] {
	return func(yield func(V) bool) {
		// TODO: range over seq, count the values, stop after the n-th
	}
}

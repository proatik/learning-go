// Package pq provides a generic priority queue backed by a binary heap.
package pq

import (
	"cmp"
	"iter"
)

// Queue is a priority queue: Pop returns the item that sorts first by less.
type Queue[T any] struct {
	items []T               // binary heap: items[i] is never less than items[(i-1)/2]
	less  func(a, b T) bool // a comes out before b
}

// New returns an empty queue ordered by less.
func New[T any](less func(a, b T) bool) *Queue[T] {
	// TODO
	return nil
}

// NewOrdered returns an empty min-queue: the smallest value comes out first.
func NewOrdered[T cmp.Ordered]() *Queue[T] {
	// TODO
	return nil
}

// Push adds values to the queue.
func (q *Queue[T]) Push(values ...T) {
	// TODO
}

// Pop removes and returns the first item, or the zero value and false if
// the queue is empty.
func (q *Queue[T]) Pop() (T, bool) {
	// TODO
	var zero T
	return zero, false
}

// Peek returns the first item without removing it.
func (q *Queue[T]) Peek() (T, bool) {
	// TODO
	var zero T
	return zero, false
}

// Len returns the number of queued items.
func (q *Queue[T]) Len() int {
	// TODO
	return 0
}

// Drain pops and yields items in priority order. Breaking out of the loop
// leaves the remaining items queued.
func (q *Queue[T]) Drain() iter.Seq[T] {
	return func(yield func(T) bool) {
		// TODO
	}
}

// TopK returns the k largest values of seq, largest first. It keeps at most
// k values in memory.
func TopK[T cmp.Ordered](seq iter.Seq[T], k int) []T {
	// TODO
	return nil
}

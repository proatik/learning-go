// Package pq provides a generic priority queue backed by a binary heap.
package pq

import (
	"cmp"
	"iter"
	"slices"
)

// Queue is a priority queue: Pop returns the item that sorts first by less.
type Queue[T any] struct {
	items []T               // binary heap: items[i] is never less than items[(i-1)/2]
	less  func(a, b T) bool // a comes out before b
}

// New returns an empty queue ordered by less.
func New[T any](less func(a, b T) bool) *Queue[T] {
	return &Queue[T]{less: less}
}

// NewOrdered returns an empty min-queue: the smallest value comes out first.
func NewOrdered[T cmp.Ordered]() *Queue[T] {
	return New(cmp.Less[T])
}

// Push adds values to the queue.
func (q *Queue[T]) Push(values ...T) {
	for _, v := range values {
		q.items = append(q.items, v)
		q.up(len(q.items) - 1)
	}
}

// Pop removes and returns the first item, or the zero value and false if
// the queue is empty.
func (q *Queue[T]) Pop() (T, bool) {
	var zero T
	if len(q.items) == 0 {
		return zero, false
	}
	top := q.items[0]
	last := len(q.items) - 1
	q.items[0] = q.items[last]
	q.items[last] = zero // let the garbage collector reclaim it
	q.items = q.items[:last]
	q.down(0)
	return top, true
}

// Peek returns the first item without removing it.
func (q *Queue[T]) Peek() (T, bool) {
	if len(q.items) == 0 {
		var zero T
		return zero, false
	}
	return q.items[0], true
}

// Len returns the number of queued items.
func (q *Queue[T]) Len() int {
	return len(q.items)
}

// Drain pops and yields items in priority order. Breaking out of the loop
// leaves the remaining items queued.
func (q *Queue[T]) Drain() iter.Seq[T] {
	return func(yield func(T) bool) {
		for {
			v, ok := q.Pop()
			if !ok || !yield(v) {
				return
			}
		}
	}
}

// up moves items[i] towards the root until its parent is not after it.
func (q *Queue[T]) up(i int) {
	for i > 0 {
		parent := (i - 1) / 2
		if !q.less(q.items[i], q.items[parent]) {
			return
		}
		q.items[i], q.items[parent] = q.items[parent], q.items[i]
		i = parent
	}
}

// down moves items[i] towards the leaves until no child comes before it.
func (q *Queue[T]) down(i int) {
	n := len(q.items)
	for {
		first := i
		for _, c := range []int{2*i + 1, 2*i + 2} {
			if c < n && q.less(q.items[c], q.items[first]) {
				first = c
			}
		}
		if first == i {
			return
		}
		q.items[i], q.items[first] = q.items[first], q.items[i]
		i = first
	}
}

// TopK returns the k largest values of seq, largest first. It keeps at most
// k values in memory.
func TopK[T cmp.Ordered](seq iter.Seq[T], k int) []T {
	if k <= 0 {
		return []T{}
	}
	smallest := NewOrdered[T]()
	for v := range seq {
		if smallest.Len() < k {
			smallest.Push(v)
		} else if top, _ := smallest.Peek(); v > top {
			smallest.Pop()
			smallest.Push(v)
		}
	}
	out := slices.Collect(smallest.Drain())
	slices.Reverse(out)
	if out == nil {
		out = []T{}
	}
	return out
}

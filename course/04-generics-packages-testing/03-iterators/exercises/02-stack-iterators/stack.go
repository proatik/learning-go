package stack

import "iter"

// Stack is a last-in, first-out collection. The zero value is an empty stack.
type Stack[T any] struct {
	items []T // items[len(items)-1] is the top
}

// Push puts v on top of the stack.
func (s *Stack[T]) Push(v T) { s.items = append(s.items, v) }

// Pop removes and returns the top item. ok is false if the stack is empty.
func (s *Stack[T]) Pop() (v T, ok bool) {
	if len(s.items) == 0 {
		return v, false
	}
	v = s.items[len(s.items)-1]
	s.items = s.items[:len(s.items)-1]
	return v, true
}

// Len returns the number of items on the stack.
func (s *Stack[T]) Len() int { return len(s.items) }

// All returns an iterator over the items from top to bottom.
// It does not modify the stack.
func (s *Stack[T]) All() iter.Seq[T] {
	return func(yield func(T) bool) {
		// TODO: walk s.items from the last index down to 0
	}
}

// Backward returns an iterator over the items from bottom to top
// (the order they were pushed). It does not modify the stack.
func (s *Stack[T]) Backward() iter.Seq[T] {
	return func(yield func(T) bool) {
		// TODO: walk s.items from index 0 up
	}
}

// Drain returns an iterator that pops each item and yields it, until the
// stack is empty or the loop stops. Items not yet yielded stay on the stack.
func (s *Stack[T]) Drain() iter.Seq[T] {
	return func(yield func(T) bool) {
		// TODO: pop, yield, repeat
	}
}

package stack

// Stack is a last-in, first-out stack of T values. The zero value is an
// empty stack.
type Stack[T any] struct {
	items []T
}

// Push puts v on top of the stack.
func (s *Stack[T]) Push(v T) {
	s.items = append(s.items, v)
}

// Pop removes and returns the top value; ok is false if the stack is empty.
func (s *Stack[T]) Pop() (T, bool) {
	v, ok := s.Peek()
	if ok {
		var zero T
		s.items[len(s.items)-1] = zero // let the garbage collector reclaim it
		s.items = s.items[:len(s.items)-1]
	}
	return v, ok
}

// Peek returns the top value without removing it; ok is false if the stack
// is empty.
func (s *Stack[T]) Peek() (T, bool) {
	if len(s.items) == 0 {
		var zero T
		return zero, false
	}
	return s.items[len(s.items)-1], true
}

// Len returns the number of values on the stack.
func (s *Stack[T]) Len() int {
	return len(s.items)
}

// Container is anything you can push T values onto and pop them from.
type Container[T any] interface {
	Push(T)
	Pop() (T, bool)
	Len() int
}

// Drain pops every value from c and returns them in pop order.
func Drain[T any](c Container[T]) []T {
	out := make([]T, 0, c.Len())
	for {
		v, ok := c.Pop()
		if !ok {
			return out
		}
		out = append(out, v)
	}
}

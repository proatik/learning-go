package stack

// Stack is a last-in, first-out stack of T values. The zero value is an
// empty stack.
type Stack[T any] struct {
	items []T
}

// Push puts v on top of the stack.
func (s *Stack[T]) Push(v T) {
	// TODO
}

// Pop removes and returns the top value; ok is false if the stack is empty.
func (s *Stack[T]) Pop() (T, bool) {
	// TODO
	var zero T
	return zero, false
}

// Peek returns the top value without removing it; ok is false if the stack
// is empty.
func (s *Stack[T]) Peek() (T, bool) {
	// TODO
	var zero T
	return zero, false
}

// Len returns the number of values on the stack.
func (s *Stack[T]) Len() int {
	// TODO
	return 0
}

// Container is anything you can push T values onto and pop them from.
type Container[T any] interface {
	Push(T)
	Pop() (T, bool)
	Len() int
}

// Drain pops every value from c and returns them in pop order.
func Drain[T any](c Container[T]) []T {
	// TODO
	return nil
}

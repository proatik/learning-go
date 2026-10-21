package stack

// Stack is a last-in, first-out stack of ints. The zero value is an empty
// stack.
type Stack []int

// Push puts v on top of the stack.
func (s *Stack) Push(v int) {
	// TODO
}

// Pop removes and returns the top value; ok is false if the stack is empty.
func (s *Stack) Pop() (int, bool) {
	// TODO
	return 0, false
}

// Peek returns the top value without removing it; ok is false if the stack
// is empty.
func (s *Stack) Peek() (int, bool) {
	// TODO
	return 0, false
}

// Len returns the number of values on the stack.
func (s *Stack) Len() int {
	// TODO
	return 0
}

// Values returns a copy of the stack contents from bottom to top.
func (s *Stack) Values() []int {
	// TODO
	return nil
}

// EvalRPN evaluates integer Reverse Polish Notation and reports whether the
// expression was valid.
func EvalRPN(tokens []string) (int, bool) {
	// TODO
	return 0, false
}

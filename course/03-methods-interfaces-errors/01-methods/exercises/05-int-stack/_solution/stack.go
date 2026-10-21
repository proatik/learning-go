package stack

import (
	"slices"
	"strconv"
)

// Stack is a last-in, first-out stack of ints. The zero value is an empty
// stack.
type Stack []int

// Push puts v on top of the stack.
func (s *Stack) Push(v int) {
	*s = append(*s, v)
}

// Pop removes and returns the top value; ok is false if the stack is empty.
func (s *Stack) Pop() (int, bool) {
	v, ok := s.Peek()
	if ok {
		*s = (*s)[:len(*s)-1]
	}
	return v, ok
}

// Peek returns the top value without removing it; ok is false if the stack
// is empty.
func (s *Stack) Peek() (int, bool) {
	if len(*s) == 0 {
		return 0, false
	}
	return (*s)[len(*s)-1], true
}

// Len returns the number of values on the stack.
func (s *Stack) Len() int {
	return len(*s)
}

// Values returns a copy of the stack contents from bottom to top.
func (s *Stack) Values() []int {
	return slices.Clone(*s)
}

// EvalRPN evaluates integer Reverse Polish Notation and reports whether the
// expression was valid.
func EvalRPN(tokens []string) (int, bool) {
	var s Stack
	for _, tok := range tokens {
		switch tok {
		case "+", "-", "*", "/":
			b, ok1 := s.Pop()
			a, ok2 := s.Pop()
			if !ok1 || !ok2 {
				return 0, false
			}
			switch tok {
			case "+":
				s.Push(a + b)
			case "-":
				s.Push(a - b)
			case "*":
				s.Push(a * b)
			case "/":
				if b == 0 {
					return 0, false
				}
				s.Push(a / b)
			}
		default:
			n, err := strconv.Atoi(tok)
			if err != nil {
				return 0, false
			}
			s.Push(n)
		}
	}
	if s.Len() != 1 {
		return 0, false
	}
	return s.Pop()
}

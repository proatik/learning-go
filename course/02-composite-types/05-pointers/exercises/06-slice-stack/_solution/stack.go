package stack

// Push puts v on top of the stack *s.
func Push(s *[]int, v int) {
	// append may allocate a new backing array; storing the result through
	// the pointer updates the caller's slice header (length and maybe array).
	*s = append(*s, v)
}

// PushAll pushes each of vs in order, so the last one ends up on top.
func PushAll(s *[]int, vs ...int) {
	for _, v := range vs {
		Push(s, v)
	}
}

// Pop removes and returns the top of the stack *s.
// It returns 0, false if the stack is empty.
func Pop(s *[]int) (int, bool) {
	n := len(*s)
	if n == 0 {
		return 0, false
	}
	top := (*s)[n-1]
	*s = (*s)[:n-1]
	return top, true
}

// Peek returns the top of the stack s without removing it.
// It returns 0, false if the stack is empty.
func Peek(s []int) (int, bool) {
	// Reading needs no pointer: the slice header is enough.
	if len(s) == 0 {
		return 0, false
	}
	return s[len(s)-1], true
}

package counter

// Counter returns a function that returns 1, 2, 3, ... on successive calls.
func Counter() func() int {
	n := 0 // captured by the closure; each Counter call gets a fresh n
	return func() int {
		n++
		return n
	}
}

// CountBy returns a function that returns start, start+step, start+2*step, ...
// on successive calls.
func CountBy(start, step int) func() int {
	next := start
	return func() int {
		cur := next
		next += step
		return cur
	}
}

// Accumulator returns a function that adds its argument to a running total
// (beginning at initial) and returns the new total.
func Accumulator(initial int) func(int) int {
	// Parameters are local variables too, so the closure can capture initial.
	return func(x int) int {
		initial += x
		return initial
	}
}

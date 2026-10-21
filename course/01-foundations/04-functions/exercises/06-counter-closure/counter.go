package counter

// Counter returns a function that returns 1, 2, 3, ... on successive calls.
func Counter() func() int {
	// TODO
	return nil
}

// CountBy returns a function that returns start, start+step, start+2*step, ...
// on successive calls.
func CountBy(start, step int) func() int {
	// TODO
	return nil
}

// Accumulator returns a function that adds its argument to a running total
// (beginning at initial) and returns the new total.
func Accumulator(initial int) func(int) int {
	// TODO
	return nil
}

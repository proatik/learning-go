package fib

import "fmt"

// Fibonacci returns a generator that yields 0, 1, 1, 2, 3, 5, ...
// one number per call.
func Fibonacci() func() int {
	a, b := 0, 1
	return func() int {
		cur := a
		a, b = b, a+b // both right-hand sides are evaluated before assigning
		return cur
	}
}

// Nth returns the n-th Fibonacci number, with Nth(0) == 0 and Nth(1) == 1.
// It returns an error if n is negative.
func Nth(n int) (int, error) {
	if n < 0 {
		return 0, fmt.Errorf("negative index %d", n)
	}
	next := Fibonacci()
	for range n { // skip F(0) .. F(n-1)
		next()
	}
	return next(), nil
}

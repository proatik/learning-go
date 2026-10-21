package hof

// IntFunc is any function from int to int.
type IntFunc func(int) int

// ApplyN applies f to x n times. If n <= 0 it returns x unchanged.
func ApplyN(f IntFunc, n, x int) int {
	// TODO
	return 0
}

// Compose returns a function that computes f(g(x)).
func Compose(f, g IntFunc) IntFunc {
	// TODO
	return nil
}

// Twice returns a function that computes f(f(x)).
func Twice(f IntFunc) IntFunc {
	// TODO
	return nil
}

package hof

// IntFunc is any function from int to int.
type IntFunc func(int) int

// ApplyN applies f to x n times. If n <= 0 it returns x unchanged.
func ApplyN(f IntFunc, n, x int) int {
	for range n { // ranging over a negative or zero n runs no iterations
		x = f(x)
	}
	return x
}

// Compose returns a function that computes f(g(x)).
func Compose(f, g IntFunc) IntFunc {
	return func(x int) int {
		return f(g(x))
	}
}

// Twice returns a function that computes f(f(x)).
func Twice(f IntFunc) IntFunc {
	return Compose(f, f)
}

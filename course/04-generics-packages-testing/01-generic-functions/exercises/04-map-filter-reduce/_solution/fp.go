package fp

// Map returns a new slice holding f(x) for every x in xs.
func Map[T, U any](xs []T, f func(T) U) []U {
	out := make([]U, 0, len(xs))
	for _, x := range xs {
		out = append(out, f(x))
	}
	return out
}

// Filter returns a new slice with the elements of xs for which keep is true.
func Filter[T any](xs []T, keep func(T) bool) []T {
	var out []T
	for _, x := range xs {
		if keep(x) {
			out = append(out, x)
		}
	}
	return out
}

// Reduce folds xs into one value, starting from init: acc = f(acc, x).
func Reduce[T, A any](xs []T, init A, f func(A, T) A) A {
	acc := init
	for _, x := range xs {
		acc = f(acc, x)
	}
	return acc
}

// Find returns the first element for which pred is true.
// If there is none it returns the zero value and false.
func Find[T any](xs []T, pred func(T) bool) (T, bool) {
	for _, x := range xs {
		if pred(x) {
			return x, true
		}
	}
	var zero T
	return zero, false
}

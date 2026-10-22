package fp

// Map returns a new slice holding f(x) for every x in xs.
func Map[T, U any](xs []T, f func(T) U) []U {
	// TODO: make([]U, 0, len(xs)) then append
	return nil
}

// Filter returns a new slice with the elements of xs for which keep is true.
func Filter[T any](xs []T, keep func(T) bool) []T {
	// TODO
	return nil
}

// Reduce folds xs into one value, starting from init: acc = f(acc, x).
func Reduce[T, A any](xs []T, init A, f func(A, T) A) A {
	// TODO
	var zero A
	return zero
}

// Find returns the first element for which pred is true.
// If there is none it returns the zero value and false.
func Find[T any](xs []T, pred func(T) bool) (T, bool) {
	// TODO
	var zero T
	return zero, false
}

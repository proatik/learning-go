package first

import "time"

// FirstOf returns the first value received from a or b, with from set to 1
// or 2 to say which channel it came from. If neither delivers within
// timeout, it returns the zero T and 0.
func FirstOf[T any](a, b <-chan T, timeout time.Duration) (v T, from int) {
	// TODO
	var zero T
	return zero, 0
}

package timeout

import (
	"errors"
	"time"
)

// ErrTimeout is returned by Do when f doesn't finish within the deadline.
var ErrTimeout = errors.New("timeout")

// Do runs f in a new goroutine and returns its result, or the zero T and
// ErrTimeout if f hasn't returned after d.
func Do[T any](d time.Duration, f func() T) (T, error) {
	// TODO
	var zero T
	return zero, nil
}

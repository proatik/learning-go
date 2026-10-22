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
	var zero T
	if d <= 0 {
		return zero, ErrTimeout
	}
	result := make(chan T, 1) // buffered: the goroutine can finish even after a timeout
	go func() { result <- f() }()
	select {
	case v := <-result:
		return v, nil
	case <-time.After(d):
		return zero, ErrTimeout
	}
}

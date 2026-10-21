package retry

import (
	"errors"
	"fmt"
)

// ErrNoAttempts is returned when Retry is asked for fewer than one attempt.
var ErrNoAttempts = errors.New("attempts must be at least 1")

// TemporaryError marks any error as temporary.
type TemporaryError struct {
	Err error
}

// Error returns the message of the wrapped error.
func (e *TemporaryError) Error() string { return e.Err.Error() }

// Temporary reports that the error is temporary.
func (e *TemporaryError) Temporary() bool { return true }

// Unwrap returns the wrapped error.
func (e *TemporaryError) Unwrap() error { return e.Err }

// temporary is the behaviour Retry looks for.
type temporary interface {
	error
	Temporary() bool
}

// IsTemporary reports whether the first error in err's tree that has a
// Temporary() bool method says it is temporary.
func IsTemporary(err error) bool {
	t, ok := errors.AsType[temporary](err)
	return ok && t.Temporary()
}

// Retry calls f up to attempts times, retrying only temporary errors.
func Retry(attempts int, f func() error) error {
	if attempts < 1 {
		return ErrNoAttempts
	}
	var last error
	for i := 1; i <= attempts; i++ {
		err := f()
		if err == nil {
			return nil
		}
		if !IsTemporary(err) {
			return fmt.Errorf("attempt %d: %w", i, err)
		}
		last = err
	}
	return fmt.Errorf("gave up after %d attempts: %w", attempts, last)
}

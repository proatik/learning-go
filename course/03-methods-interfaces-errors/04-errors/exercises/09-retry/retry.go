package retry

import "errors"

// ErrNoAttempts is returned when Retry is asked for fewer than one attempt.
var ErrNoAttempts = errors.New("attempts must be at least 1")

// TemporaryError marks any error as temporary.
type TemporaryError struct {
	Err error
}

// Error returns the message of the wrapped error.
func (e *TemporaryError) Error() string {
	// TODO
	return ""
}

// Temporary reports that the error is temporary.
func (e *TemporaryError) Temporary() bool {
	// TODO
	return false
}

// Unwrap returns the wrapped error.
func (e *TemporaryError) Unwrap() error {
	// TODO
	return nil
}

// IsTemporary reports whether the first error in err's tree that has a
// Temporary() bool method says it is temporary.
func IsTemporary(err error) bool {
	// TODO
	return false
}

// Retry calls f up to attempts times, retrying only temporary errors.
func Retry(attempts int, f func() error) error {
	// TODO
	return nil
}

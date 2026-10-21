package middleware

import "errors"

// ErrAbort is a deliberate panic value that Protect must not swallow.
var ErrAbort = errors.New("abort handler")

// PanicError reports a panic recovered from a handler.
type PanicError struct {
	Input string // the input the handler was called with
	Value any    // the recovered panic value
}

// Error returns `panic handling "<Input>": <Value>`.
func (e *PanicError) Error() string {
	// TODO
	return ""
}

// Unwrap returns Value if it is an error, otherwise nil.
func (e *PanicError) Unwrap() error {
	// TODO
	return nil
}

// Protect wraps h so that a panic becomes a *PanicError, except panics with
// an error matching ErrAbort, which are re-panicked unchanged.
func Protect(h func(string) string) func(string) (string, error) {
	// TODO
	return func(string) (string, error) { return "", nil }
}

// Pipeline runs input through every step (each protected by Protect) and
// returns the final output, or "" and "step <i>: ..." on the first error.
func Pipeline(input string, steps []func(string) string) (string, error) {
	// TODO
	return "", nil
}

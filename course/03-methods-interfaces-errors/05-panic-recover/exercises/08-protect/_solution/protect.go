package middleware

import (
	"errors"
	"fmt"
)

// ErrAbort is a deliberate panic value that Protect must not swallow.
var ErrAbort = errors.New("abort handler")

// PanicError reports a panic recovered from a handler.
type PanicError struct {
	Input string // the input the handler was called with
	Value any    // the recovered panic value
}

// Error returns `panic handling "<Input>": <Value>`.
func (e *PanicError) Error() string {
	return fmt.Sprintf("panic handling %q: %v", e.Input, e.Value)
}

// Unwrap returns Value if it is an error, otherwise nil.
func (e *PanicError) Unwrap() error {
	err, _ := e.Value.(error) // nil when Value is not an error
	return err
}

// Protect wraps h so that a panic becomes a *PanicError, except panics with
// an error matching ErrAbort, which are re-panicked unchanged.
func Protect(h func(string) string) func(string) (string, error) {
	return func(input string) (out string, err error) {
		defer func() {
			r := recover()
			if r == nil {
				return
			}
			if e, ok := r.(error); ok && errors.Is(e, ErrAbort) {
				panic(r) // deliberate abort: not ours to handle
			}
			out, err = "", &PanicError{Input: input, Value: r}
		}()
		return h(input), nil
	}
}

// Pipeline runs input through every step (each protected by Protect) and
// returns the final output, or "" and "step <i>: ..." on the first error.
func Pipeline(input string, steps []func(string) string) (string, error) {
	out := input
	for i, step := range steps {
		next, err := Protect(step)(out)
		if err != nil {
			return "", fmt.Errorf("step %d: %w", i, err)
		}
		out = next
	}
	return out, nil
}

// Package result provides a generic value-or-error type with chaining.
package result

import (
	"errors"
	"iter"
)

// ErrPanic is wrapped by results that Try built from a panic.
var ErrPanic = errors.New("result: panic")

// Result holds either a value or an error.
type Result[T any] struct {
	val T
	err error
}

// Ok returns a successful result holding v.
func Ok[T any](v T) Result[T] {
	// TODO
	return Result[T]{}
}

// Fail returns a failed result holding err.
func Fail[T any](err error) Result[T] {
	// TODO
	return Result[T]{}
}

// Of turns a (value, error) pair into a Result; the value is dropped if err
// is non-nil.
func Of[T any](v T, err error) Result[T] {
	// TODO
	return Result[T]{}
}

// Try calls f and returns its value, or a failed result wrapping ErrPanic if
// f panics.
func Try[T any](f func() T) (r Result[T]) {
	// TODO
	return Result[T]{}
}

// Get returns the value and error. The value is the zero T on failure.
func (r Result[T]) Get() (T, error) {
	// TODO
	var zero T
	return zero, nil
}

// IsOk reports whether r holds a value.
func (r Result[T]) IsOk() bool {
	// TODO
	return false
}

// OrElse returns the value, or fallback if r failed.
func (r Result[T]) OrElse(fallback T) T {
	// TODO
	var zero T
	return zero
}

// Context wraps the error of a failed result as "msg: err". A successful
// result is returned unchanged.
func (r Result[T]) Context(msg string) Result[T] {
	// TODO
	return r
}

// Map applies f to the value. A failed result passes its error through and f
// is not called.
func (r Result[T]) Map[U any](f func(T) U) Result[U] {
	// TODO
	return Result[U]{}
}

// Then applies a step that can fail. A failed result passes its error
// through and f is not called.
func (r Result[T]) Then[U any](f func(T) (U, error)) Result[U] {
	// TODO
	return Result[U]{}
}

// Collect gathers the values of seq. It stops at the first failure and
// returns its error as "item i: err", with i counted from 0.
func Collect[T any](seq iter.Seq[Result[T]]) ([]T, error) {
	// TODO
	return nil, nil
}

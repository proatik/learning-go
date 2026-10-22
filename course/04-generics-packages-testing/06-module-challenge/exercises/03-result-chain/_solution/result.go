// Package result provides a generic value-or-error type with chaining.
package result

import (
	"errors"
	"fmt"
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
	return Result[T]{val: v}
}

// Fail returns a failed result holding err.
func Fail[T any](err error) Result[T] {
	return Result[T]{err: err}
}

// Of turns a (value, error) pair into a Result; the value is dropped if err
// is non-nil.
func Of[T any](v T, err error) Result[T] {
	if err != nil {
		return Fail[T](err)
	}
	return Ok(v)
}

// Try calls f and returns its value, or a failed result wrapping ErrPanic if
// f panics.
func Try[T any](f func() T) (r Result[T]) {
	defer func() {
		if p := recover(); p != nil {
			r = Fail[T](fmt.Errorf("%w: %v", ErrPanic, p))
		}
	}()
	return Ok(f())
}

// Get returns the value and error. The value is the zero T on failure.
func (r Result[T]) Get() (T, error) {
	return r.val, r.err
}

// IsOk reports whether r holds a value.
func (r Result[T]) IsOk() bool {
	return r.err == nil
}

// OrElse returns the value, or fallback if r failed.
func (r Result[T]) OrElse(fallback T) T {
	if r.err != nil {
		return fallback
	}
	return r.val
}

// Context wraps the error of a failed result as "msg: err". A successful
// result is returned unchanged.
func (r Result[T]) Context(msg string) Result[T] {
	if r.err != nil {
		return Fail[T](fmt.Errorf("%s: %w", msg, r.err))
	}
	return r
}

// Map applies f to the value. A failed result passes its error through and f
// is not called.
func (r Result[T]) Map[U any](f func(T) U) Result[U] {
	if r.err != nil {
		return Fail[U](r.err)
	}
	return Ok(f(r.val))
}

// Then applies a step that can fail. A failed result passes its error
// through and f is not called.
func (r Result[T]) Then[U any](f func(T) (U, error)) Result[U] {
	if r.err != nil {
		return Fail[U](r.err)
	}
	return Of(f(r.val))
}

// Collect gathers the values of seq. It stops at the first failure and
// returns its error as "item i: err", with i counted from 0.
func Collect[T any](seq iter.Seq[Result[T]]) ([]T, error) {
	out := []T{}
	i := 0
	for r := range seq {
		v, err := r.Get()
		if err != nil {
			return nil, fmt.Errorf("item %d: %w", i, err)
		}
		out = append(out, v)
		i++
	}
	return out, nil
}

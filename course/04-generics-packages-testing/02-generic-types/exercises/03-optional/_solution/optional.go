package optional

import "fmt"

// Optional holds either a value of type T or nothing. The zero value is
// None.
type Optional[T any] struct {
	value T
	ok    bool
}

// Some returns an Optional holding v.
func Some[T any](v T) Optional[T] {
	return Optional[T]{value: v, ok: true}
}

// None returns an empty Optional. T can't be inferred: call None[int]().
func None[T any]() Optional[T] {
	return Optional[T]{}
}

// Get returns the value and true, or the zero T and false for None.
func (o Optional[T]) Get() (T, bool) {
	return o.value, o.ok
}

// IsSome reports whether o holds a value.
func (o Optional[T]) IsSome() bool {
	return o.ok
}

// OrElse returns the value, or fallback for None.
func (o Optional[T]) OrElse(fallback T) T {
	if o.ok {
		return o.value
	}
	return fallback
}

// String formats o as "Some(<v>)" or "None".
func (o Optional[T]) String() string {
	if !o.ok {
		return "None"
	}
	return fmt.Sprintf("Some(%v)", o.value)
}

// First returns the first element of s, or None if s is empty.
func First[T any](s []T) Optional[T] {
	if len(s) == 0 {
		return None[T]()
	}
	return Some(s[0])
}

// Lookup returns m[k] if the key is present, or None.
func Lookup[K comparable, V any](m map[K]V, k K) Optional[V] {
	v, ok := m[k]
	if !ok {
		return None[V]()
	}
	return Some(v)
}

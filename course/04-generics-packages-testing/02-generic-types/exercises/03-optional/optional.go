package optional

// Optional holds either a value of type T or nothing. The zero value is
// None.
type Optional[T any] struct {
	value T
	ok    bool
}

// Some returns an Optional holding v.
func Some[T any](v T) Optional[T] {
	// TODO
	return Optional[T]{}
}

// None returns an empty Optional. T can't be inferred: call None[int]().
func None[T any]() Optional[T] {
	// TODO
	return Optional[T]{}
}

// Get returns the value and true, or the zero T and false for None.
func (o Optional[T]) Get() (T, bool) {
	// TODO
	var zero T
	return zero, false
}

// IsSome reports whether o holds a value.
func (o Optional[T]) IsSome() bool {
	// TODO
	return false
}

// OrElse returns the value, or fallback for None.
func (o Optional[T]) OrElse(fallback T) T {
	// TODO
	var zero T
	return zero
}

// String formats o as "Some(<v>)" or "None".
func (o Optional[T]) String() string {
	// TODO
	return ""
}

// First returns the first element of s, or None if s is empty.
func First[T any](s []T) Optional[T] {
	// TODO
	return Optional[T]{}
}

// Lookup returns m[k] if the key is present, or None.
func Lookup[K comparable, V any](m map[K]V, k K) Optional[V] {
	// TODO
	return Optional[V]{}
}

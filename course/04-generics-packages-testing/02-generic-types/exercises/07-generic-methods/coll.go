package coll

// List is an immutable-style list: methods never modify the receiver, they
// return new lists. The zero value is an empty list.
type List[T any] struct {
	items []T
}

// Of returns a list holding a copy of vals.
func Of[T any](vals ...T) List[T] {
	// TODO
	return List[T]{}
}

// Items returns a copy of the elements.
func (l List[T]) Items() []T {
	// TODO
	return nil
}

// Len returns the number of elements.
func (l List[T]) Len() int {
	// TODO
	return 0
}

// Filter returns the elements for which keep returns true.
func (l List[T]) Filter(keep func(T) bool) List[T] {
	// TODO
	return List[T]{}
}

// Map returns a list holding f(v) for every element v. U is the method's own
// type parameter (Go 1.27), so the result can have a different element type.
func (l List[T]) Map[U any](f func(T) U) List[U] {
	// TODO
	return List[U]{}
}

// FlatMap returns the concatenation of f(v) for every element v.
func (l List[T]) FlatMap[U any](f func(T) []U) List[U] {
	// TODO
	return List[U]{}
}

// Reduce folds the elements from left to right into one value, starting
// from init.
func (l List[T]) Reduce[A any](init A, f func(A, T) A) A {
	// TODO
	var zero A
	return zero
}

// GroupBy splits the list by key(v), keeping the original order inside each
// group.
func (l List[T]) GroupBy[K comparable](key func(T) K) map[K]List[T] {
	// TODO
	return nil
}

// Number is the set of numeric types Sum accepts.
type Number interface {
	~int | ~int64 | ~float64
}

// Sum adds up the elements. It is a function, not a method: a method of
// List[T] only knows T is "any", and a method's own type parameters can't
// tighten T.
func Sum[N Number](l List[N]) N {
	// TODO
	return 0
}

package set

import "cmp"

// Set is an unordered collection of distinct comparable values. A nil Set is
// a valid empty set for reading; use New or make before adding.
type Set[T comparable] map[T]struct{}

// New returns a non-nil set holding vals.
func New[T comparable](vals ...T) Set[T] {
	// TODO
	return nil
}

// Add inserts every value in vals.
func (s Set[T]) Add(vals ...T) {
	// TODO
}

// Remove deletes v and reports whether it was present.
func (s Set[T]) Remove(v T) bool {
	// TODO
	return false
}

// Has reports whether v is in the set.
func (s Set[T]) Has(v T) bool {
	// TODO
	return false
}

// Len returns the number of elements.
func (s Set[T]) Len() int {
	// TODO
	return 0
}

// Union returns a new set with the elements of s or o.
func (s Set[T]) Union(o Set[T]) Set[T] {
	// TODO
	return nil
}

// Intersect returns a new set with the elements of both s and o.
func (s Set[T]) Intersect(o Set[T]) Set[T] {
	// TODO
	return nil
}

// Difference returns a new set with the elements of s that are not in o.
func (s Set[T]) Difference(o Set[T]) Set[T] {
	// TODO
	return nil
}

// SubsetOf reports whether every element of s is in o.
func (s Set[T]) SubsetOf(o Set[T]) bool {
	// TODO
	return false
}

// Equal reports whether s and o hold the same elements.
func (s Set[T]) Equal(o Set[T]) bool {
	// TODO
	return false
}

// Sorted returns the elements of s in ascending order. It is a function, not
// a method, because it needs cmp.Ordered, which is stricter than the
// comparable constraint every Set method is limited to.
func Sorted[T cmp.Ordered](s Set[T]) []T {
	// TODO
	return nil
}

package set

import (
	"cmp"
	"slices"
)

// Set is an unordered collection of distinct comparable values. A nil Set is
// a valid empty set for reading; use New or make before adding.
type Set[T comparable] map[T]struct{}

// New returns a non-nil set holding vals.
func New[T comparable](vals ...T) Set[T] {
	s := make(Set[T], len(vals))
	s.Add(vals...)
	return s
}

// Add inserts every value in vals.
func (s Set[T]) Add(vals ...T) {
	for _, v := range vals {
		s[v] = struct{}{}
	}
}

// Remove deletes v and reports whether it was present.
func (s Set[T]) Remove(v T) bool {
	if !s.Has(v) {
		return false
	}
	delete(s, v)
	return true
}

// Has reports whether v is in the set.
func (s Set[T]) Has(v T) bool {
	_, ok := s[v]
	return ok
}

// Len returns the number of elements.
func (s Set[T]) Len() int {
	return len(s)
}

// Union returns a new set with the elements of s or o.
func (s Set[T]) Union(o Set[T]) Set[T] {
	out := make(Set[T], len(s)+len(o))
	for v := range s {
		out.Add(v)
	}
	for v := range o {
		out.Add(v)
	}
	return out
}

// Intersect returns a new set with the elements of both s and o.
func (s Set[T]) Intersect(o Set[T]) Set[T] {
	out := make(Set[T])
	for v := range s {
		if o.Has(v) {
			out.Add(v)
		}
	}
	return out
}

// Difference returns a new set with the elements of s that are not in o.
func (s Set[T]) Difference(o Set[T]) Set[T] {
	out := make(Set[T])
	for v := range s {
		if !o.Has(v) {
			out.Add(v)
		}
	}
	return out
}

// SubsetOf reports whether every element of s is in o.
func (s Set[T]) SubsetOf(o Set[T]) bool {
	for v := range s {
		if !o.Has(v) {
			return false
		}
	}
	return true
}

// Equal reports whether s and o hold the same elements.
func (s Set[T]) Equal(o Set[T]) bool {
	return len(s) == len(o) && s.SubsetOf(o)
}

// Sorted returns the elements of s in ascending order. It is a function, not
// a method, because it needs cmp.Ordered, which is stricter than the
// comparable constraint every Set method is limited to.
func Sorted[T cmp.Ordered](s Set[T]) []T {
	out := make([]T, 0, len(s))
	for v := range s {
		out = append(out, v)
	}
	slices.Sort(out)
	return out
}

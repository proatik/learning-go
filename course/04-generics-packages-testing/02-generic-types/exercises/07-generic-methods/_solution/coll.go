package coll

import "slices"

// List is an immutable-style list: methods never modify the receiver, they
// return new lists. The zero value is an empty list.
type List[T any] struct {
	items []T
}

// Of returns a list holding a copy of vals.
func Of[T any](vals ...T) List[T] {
	return List[T]{items: slices.Clone(vals)}
}

// Items returns a copy of the elements.
func (l List[T]) Items() []T {
	return slices.Clone(l.items)
}

// Len returns the number of elements.
func (l List[T]) Len() int {
	return len(l.items)
}

// Filter returns the elements for which keep returns true.
func (l List[T]) Filter(keep func(T) bool) List[T] {
	var out []T
	for _, v := range l.items {
		if keep(v) {
			out = append(out, v)
		}
	}
	return List[T]{items: out}
}

// Map returns a list holding f(v) for every element v. U is the method's own
// type parameter (Go 1.27), so the result can have a different element type.
func (l List[T]) Map[U any](f func(T) U) List[U] {
	out := make([]U, len(l.items))
	for i, v := range l.items {
		out[i] = f(v)
	}
	return List[U]{items: out}
}

// FlatMap returns the concatenation of f(v) for every element v.
func (l List[T]) FlatMap[U any](f func(T) []U) List[U] {
	var out []U
	for _, v := range l.items {
		out = append(out, f(v)...)
	}
	return List[U]{items: out}
}

// Reduce folds the elements from left to right into one value, starting
// from init.
func (l List[T]) Reduce[A any](init A, f func(A, T) A) A {
	acc := init
	for _, v := range l.items {
		acc = f(acc, v)
	}
	return acc
}

// GroupBy splits the list by key(v), keeping the original order inside each
// group.
func (l List[T]) GroupBy[K comparable](key func(T) K) map[K]List[T] {
	groups := make(map[K]List[T])
	for _, v := range l.items {
		k := key(v)
		g := groups[k]
		g.items = append(g.items, v)
		groups[k] = g
	}
	return groups
}

// Number is the set of numeric types Sum accepts.
type Number interface {
	~int | ~int64 | ~float64
}

// Sum adds up the elements. It is a function, not a method: a method of
// List[T] only knows T is "any", and a method's own type parameters can't
// tighten T.
func Sum[N Number](l List[N]) N {
	return l.Reduce(0, func(acc, v N) N { return acc + v })
}

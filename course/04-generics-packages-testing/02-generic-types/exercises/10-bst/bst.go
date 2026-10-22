package bst

import "cmp"

// node is one tree node; it refers to its own instantiation for children.
type node[T cmp.Ordered] struct {
	val         T
	left, right *node[T]
}

// Tree is an unbalanced binary search tree of distinct values. The zero
// value is an empty tree.
type Tree[T cmp.Ordered] struct {
	root *node[T]
	size int
}

// From returns a tree holding vals, inserted in order.
func From[T cmp.Ordered](vals ...T) *Tree[T] {
	// TODO
	return &Tree[T]{}
}

// Insert adds v and reports whether it was new.
func (t *Tree[T]) Insert(v T) bool {
	// TODO
	return false
}

// Contains reports whether v is in the tree.
func (t *Tree[T]) Contains(v T) bool {
	// TODO
	return false
}

// Delete removes v and reports whether it was present.
func (t *Tree[T]) Delete(v T) bool {
	// TODO
	return false
}

// Len returns the number of values.
func (t *Tree[T]) Len() int {
	// TODO
	return 0
}

// Height returns the number of nodes on the longest root-to-leaf path.
func (t *Tree[T]) Height() int {
	// TODO
	return 0
}

// Min returns the smallest value; ok is false if the tree is empty.
func (t *Tree[T]) Min() (T, bool) {
	// TODO
	var zero T
	return zero, false
}

// Max returns the largest value; ok is false if the tree is empty.
func (t *Tree[T]) Max() (T, bool) {
	// TODO
	var zero T
	return zero, false
}

// Floor returns the largest value <= v; ok is false if there is none.
func (t *Tree[T]) Floor(v T) (T, bool) {
	// TODO
	var zero T
	return zero, false
}

// Ceil returns the smallest value >= v; ok is false if there is none.
func (t *Tree[T]) Ceil(v T) (T, bool) {
	// TODO
	var zero T
	return zero, false
}

// InOrder returns the values in ascending order.
func (t *Tree[T]) InOrder() []T {
	// TODO
	return nil
}

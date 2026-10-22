package bst

import (
	"cmp"
	"iter"
)

// Tree is a binary search tree of unique keys. The zero value is an empty tree.
type Tree[K cmp.Ordered] struct {
	root *node[K]
	size int
}

type node[K cmp.Ordered] struct {
	key         K
	left, right *node[K] // left: smaller keys, right: larger keys
}

// Insert adds k to the tree. It reports false if k was already present.
func (t *Tree[K]) Insert(k K) bool {
	p := &t.root
	for *p != nil {
		switch c := cmp.Compare(k, (*p).key); {
		case c < 0:
			p = &(*p).left
		case c > 0:
			p = &(*p).right
		default:
			return false
		}
	}
	*p = &node[K]{key: k}
	t.size++
	return true
}

// Len returns the number of keys in the tree.
func (t *Tree[K]) Len() int { return t.size }

// All returns an iterator over the keys in ascending order.
func (t *Tree[K]) All() iter.Seq[K] {
	return func(yield func(K) bool) {
		// TODO: call a recursive helper on t.root
	}
}

// Backward returns an iterator over the keys in descending order.
func (t *Tree[K]) Backward() iter.Seq[K] {
	return func(yield func(K) bool) {
		// TODO: the mirror image of All
	}
}

// Range returns an iterator over the keys k with lo <= k <= hi,
// in ascending order.
func (t *Tree[K]) Range(lo, hi K) iter.Seq[K] {
	return func(yield func(K) bool) {
		// TODO: like All, but skip subtrees that can't hold keys in [lo, hi]
	}
}

// TODO: write the recursive helpers here. Each one should report whether
// the walk may continue, e.g.
//
//	func walk[K cmp.Ordered](n *node[K], yield func(K) bool) bool

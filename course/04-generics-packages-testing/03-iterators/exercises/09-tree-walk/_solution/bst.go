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
		walk(t.root, yield)
	}
}

// Backward returns an iterator over the keys in descending order.
func (t *Tree[K]) Backward() iter.Seq[K] {
	return func(yield func(K) bool) {
		walkBackward(t.root, yield)
	}
}

// Range returns an iterator over the keys k with lo <= k <= hi,
// in ascending order.
func (t *Tree[K]) Range(lo, hi K) iter.Seq[K] {
	return func(yield func(K) bool) {
		walkRange(t.root, lo, hi, yield)
	}
}

// walk yields the keys under n in ascending order and reports whether
// the caller should keep going.
func walk[K cmp.Ordered](n *node[K], yield func(K) bool) bool {
	if n == nil {
		return true
	}
	return walk(n.left, yield) && yield(n.key) && walk(n.right, yield)
}

// walkBackward yields the keys under n in descending order.
func walkBackward[K cmp.Ordered](n *node[K], yield func(K) bool) bool {
	if n == nil {
		return true
	}
	return walkBackward(n.right, yield) && yield(n.key) && walkBackward(n.left, yield)
}

// walkRange yields the keys under n within [lo, hi], skipping subtrees
// that lie entirely outside the range.
func walkRange[K cmp.Ordered](n *node[K], lo, hi K, yield func(K) bool) bool {
	if n == nil {
		return true
	}
	if n.key > lo && !walkRange(n.left, lo, hi, yield) {
		return false
	}
	if lo <= n.key && n.key <= hi && !yield(n.key) {
		return false
	}
	if n.key < hi {
		return walkRange(n.right, lo, hi, yield)
	}
	return true
}

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
	t := &Tree[T]{}
	for _, v := range vals {
		t.Insert(v)
	}
	return t
}

// Insert adds v and reports whether it was new.
func (t *Tree[T]) Insert(v T) bool {
	link := &t.root
	for *link != nil {
		switch c := cmp.Compare(v, (*link).val); {
		case c < 0:
			link = &(*link).left
		case c > 0:
			link = &(*link).right
		default:
			return false
		}
	}
	*link = &node[T]{val: v}
	t.size++
	return true
}

// Contains reports whether v is in the tree.
func (t *Tree[T]) Contains(v T) bool {
	n := t.root
	for n != nil {
		switch c := cmp.Compare(v, n.val); {
		case c < 0:
			n = n.left
		case c > 0:
			n = n.right
		default:
			return true
		}
	}
	return false
}

// Delete removes v and reports whether it was present.
func (t *Tree[T]) Delete(v T) bool {
	var found bool
	t.root, found = deleteNode(t.root, v)
	if found {
		t.size--
	}
	return found
}

// deleteNode removes v from the subtree n and returns the new subtree root.
func deleteNode[T cmp.Ordered](n *node[T], v T) (*node[T], bool) {
	if n == nil {
		return nil, false
	}
	var found bool
	switch c := cmp.Compare(v, n.val); {
	case c < 0:
		n.left, found = deleteNode(n.left, v)
		return n, found
	case c > 0:
		n.right, found = deleteNode(n.right, v)
		return n, found
	}
	switch {
	case n.left == nil:
		return n.right, true
	case n.right == nil:
		return n.left, true
	}
	succ := n.right
	for succ.left != nil {
		succ = succ.left
	}
	n.val = succ.val
	n.right, _ = deleteNode(n.right, succ.val)
	return n, true
}

// Len returns the number of values.
func (t *Tree[T]) Len() int {
	return t.size
}

// Height returns the number of nodes on the longest root-to-leaf path.
func (t *Tree[T]) Height() int {
	return height(t.root)
}

// height returns the height of the subtree n.
func height[T cmp.Ordered](n *node[T]) int {
	if n == nil {
		return 0
	}
	return 1 + max(height(n.left), height(n.right))
}

// Min returns the smallest value; ok is false if the tree is empty.
func (t *Tree[T]) Min() (T, bool) {
	if t.root == nil {
		var zero T
		return zero, false
	}
	n := t.root
	for n.left != nil {
		n = n.left
	}
	return n.val, true
}

// Max returns the largest value; ok is false if the tree is empty.
func (t *Tree[T]) Max() (T, bool) {
	if t.root == nil {
		var zero T
		return zero, false
	}
	n := t.root
	for n.right != nil {
		n = n.right
	}
	return n.val, true
}

// Floor returns the largest value <= v; ok is false if there is none.
func (t *Tree[T]) Floor(v T) (T, bool) {
	var best T
	ok := false
	for n := t.root; n != nil; {
		switch c := cmp.Compare(v, n.val); {
		case c == 0:
			return n.val, true
		case c < 0:
			n = n.left
		default:
			best, ok = n.val, true
			n = n.right
		}
	}
	return best, ok
}

// Ceil returns the smallest value >= v; ok is false if there is none.
func (t *Tree[T]) Ceil(v T) (T, bool) {
	var best T
	ok := false
	for n := t.root; n != nil; {
		switch c := cmp.Compare(v, n.val); {
		case c == 0:
			return n.val, true
		case c > 0:
			n = n.right
		default:
			best, ok = n.val, true
			n = n.left
		}
	}
	return best, ok
}

// InOrder returns the values in ascending order.
func (t *Tree[T]) InOrder() []T {
	out := make([]T, 0, t.size)
	return appendInOrder(out, t.root)
}

// appendInOrder appends the values of the subtree n to out, in order.
func appendInOrder[T cmp.Ordered](out []T, n *node[T]) []T {
	if n == nil {
		return out
	}
	out = appendInOrder(out, n.left)
	out = append(out, n.val)
	return appendInOrder(out, n.right)
}

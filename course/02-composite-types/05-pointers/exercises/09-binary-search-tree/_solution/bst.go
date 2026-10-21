package bst

// Tree is a node of a binary search tree. Keys in Left are smaller than Key,
// keys in Right are larger. A nil *Tree is the empty tree.
type Tree struct {
	Key         int
	Left, Right *Tree
}

// Insert adds key to the tree rooted at t and returns the root.
// Duplicate keys are ignored.
func Insert(t *Tree, key int) *Tree {
	if t == nil {
		return &Tree{Key: key} // found the empty spot: a new leaf
	}
	switch {
	case key < t.Key:
		t.Left = Insert(t.Left, key) // re-link: Left may have been nil
	case key > t.Key:
		t.Right = Insert(t.Right, key)
	}
	return t
}

// Contains reports whether key is in the tree rooted at t.
func Contains(t *Tree, key int) bool {
	// Iterative: follow one path down, no recursion needed.
	for t != nil {
		switch {
		case key < t.Key:
			t = t.Left
		case key > t.Key:
			t = t.Right
		default:
			return true
		}
	}
	return false
}

// InOrder returns the keys of the tree in ascending (in-order) order,
// or nil for the empty tree.
func InOrder(t *Tree) []int {
	return appendInOrder(nil, t)
}

// appendInOrder appends the keys of t, in order, to out and returns the
// extended slice (the same convention as the built-in append).
func appendInOrder(out []int, t *Tree) []int {
	if t == nil {
		return out
	}
	out = appendInOrder(out, t.Left)
	out = append(out, t.Key)
	return appendInOrder(out, t.Right)
}

// Height returns the number of nodes on the longest path from the root
// to a leaf: 0 for the empty tree.
func Height(t *Tree) int {
	if t == nil {
		return 0
	}
	return 1 + max(Height(t.Left), Height(t.Right))
}

// Min returns the smallest key in the tree and true,
// or 0, false for the empty tree.
func Min(t *Tree) (int, bool) {
	if t == nil {
		return 0, false
	}
	for t.Left != nil { // the smallest key is the leftmost node
		t = t.Left
	}
	return t.Key, true
}

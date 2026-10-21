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
	// TODO
	return nil
}

// Contains reports whether key is in the tree rooted at t.
func Contains(t *Tree, key int) bool {
	// TODO
	return false
}

// InOrder returns the keys of the tree in ascending (in-order) order,
// or nil for the empty tree.
func InOrder(t *Tree) []int {
	// TODO
	return nil
}

// Height returns the number of nodes on the longest path from the root
// to a leaf: 0 for the empty tree.
func Height(t *Tree) int {
	// TODO
	return 0
}

// Min returns the smallest key in the tree and true,
// or 0, false for the empty tree.
func Min(t *Tree) (int, bool) {
	// TODO
	return 0, false
}

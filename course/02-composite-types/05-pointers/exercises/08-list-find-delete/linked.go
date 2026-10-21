package linked

// Node is one element of a singly linked list. A nil *Node is the empty list.
type Node struct {
	Val  int
	Next *Node
}

// Find returns the first node holding v, or nil if there is none.
func Find(head *Node, v int) *Node {
	// TODO
	return nil
}

// InsertAfter links a new node holding v directly after n.
// It does nothing if n is nil.
func InsertAfter(n *Node, v int) {
	// TODO
}

// Delete removes the first node holding v and returns the head of the
// resulting list.
func Delete(head *Node, v int) *Node {
	// TODO
	return nil
}

// DeleteAll removes every node holding v and returns the head of the
// resulting list.
func DeleteAll(head *Node, v int) *Node {
	// TODO
	return nil
}

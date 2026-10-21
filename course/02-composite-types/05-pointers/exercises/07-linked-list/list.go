package list

// Node is one element of a singly linked list. A nil *Node is the empty list.
type Node struct {
	Val  int
	Next *Node
}

// Push returns a new list with v in front of head.
func Push(head *Node, v int) *Node {
	// TODO
	return nil
}

// FromSlice returns a list holding vals in the same order,
// or nil if vals is empty.
func FromSlice(vals []int) *Node {
	// TODO
	return nil
}

// Len returns the number of nodes in the list.
func Len(head *Node) int {
	// TODO
	return 0
}

// ToSlice returns the values of the list in order, or nil for an empty list.
func ToSlice(head *Node) []int {
	// TODO
	return nil
}

// Reverse reverses the list in place and returns the new head.
func Reverse(head *Node) *Node {
	// TODO
	return nil
}

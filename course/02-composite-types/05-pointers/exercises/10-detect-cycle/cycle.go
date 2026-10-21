package cycle

// Node is one element of a singly linked list.
type Node struct {
	Val  int
	Next *Node
}

// HasCycle reports whether following Next from head ever revisits a node.
func HasCycle(head *Node) bool {
	// TODO
	return false
}

// CycleStart returns the first node of the cycle (the first node reached
// twice when walking from head), or nil if the list has no cycle.
func CycleStart(head *Node) *Node {
	// TODO
	return nil
}

// CycleLen returns the number of nodes in the cycle, or 0 if there is none.
func CycleLen(head *Node) int {
	// TODO
	return 0
}

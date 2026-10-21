package cycle

// Node is one element of a singly linked list.
type Node struct {
	Val  int
	Next *Node
}

// meet runs Floyd's tortoise and hare. It returns the node where the two
// pointers meet inside the cycle, or nil if the list has no cycle.
func meet(head *Node) *Node {
	slow, fast := head, head
	for fast != nil && fast.Next != nil {
		slow = slow.Next
		fast = fast.Next.Next
		if slow == fast { // same node (identity), not just same Val
			return slow
		}
	}
	return nil // fast fell off the end: the list terminates
}

// HasCycle reports whether following Next from head ever revisits a node.
func HasCycle(head *Node) bool {
	return meet(head) != nil
}

// CycleStart returns the first node of the cycle (the first node reached
// twice when walking from head), or nil if the list has no cycle.
func CycleStart(head *Node) *Node {
	m := meet(head)
	if m == nil {
		return nil
	}
	// The distance from head to the cycle start equals the distance from
	// the meeting point to the cycle start (going forward around the cycle).
	a, b := head, m
	for a != b {
		a, b = a.Next, b.Next
	}
	return a
}

// CycleLen returns the number of nodes in the cycle, or 0 if there is none.
func CycleLen(head *Node) int {
	m := meet(head)
	if m == nil {
		return 0
	}
	n := 1
	for p := m.Next; p != m; p = p.Next {
		n++
	}
	return n
}

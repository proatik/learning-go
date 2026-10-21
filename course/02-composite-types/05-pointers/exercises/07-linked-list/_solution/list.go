package list

// Node is one element of a singly linked list. A nil *Node is the empty list.
type Node struct {
	Val  int
	Next *Node
}

// Push returns a new list with v in front of head.
func Push(head *Node, v int) *Node {
	return &Node{Val: v, Next: head}
}

// FromSlice returns a list holding vals in the same order,
// or nil if vals is empty.
func FromSlice(vals []int) *Node {
	// Pushing puts each value in front, so walk the slice backwards.
	var head *Node
	for i := len(vals) - 1; i >= 0; i-- {
		head = Push(head, vals[i])
	}
	return head
}

// Len returns the number of nodes in the list.
func Len(head *Node) int {
	n := 0
	for node := head; node != nil; node = node.Next {
		n++
	}
	return n
}

// ToSlice returns the values of the list in order, or nil for an empty list.
func ToSlice(head *Node) []int {
	var vals []int // stays nil for an empty list
	for n := head; n != nil; n = n.Next {
		vals = append(vals, n.Val)
	}
	return vals
}

// Reverse reverses the list in place and returns the new head.
func Reverse(head *Node) *Node {
	var prev *Node
	cur := head
	for cur != nil {
		next := cur.Next // remember the rest before overwriting the link
		cur.Next = prev  // point this node backwards
		prev, cur = cur, next
	}
	return prev // the old tail
}

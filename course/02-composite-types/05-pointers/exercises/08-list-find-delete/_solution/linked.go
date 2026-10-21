package linked

// Node is one element of a singly linked list. A nil *Node is the empty list.
type Node struct {
	Val  int
	Next *Node
}

// Find returns the first node holding v, or nil if there is none.
func Find(head *Node, v int) *Node {
	for n := head; n != nil; n = n.Next {
		if n.Val == v {
			return n
		}
	}
	return nil
}

// InsertAfter links a new node holding v directly after n.
// It does nothing if n is nil.
func InsertAfter(n *Node, v int) {
	if n == nil {
		return
	}
	n.Next = &Node{Val: v, Next: n.Next}
}

// Delete removes the first node holding v and returns the head of the
// resulting list.
func Delete(head *Node, v int) *Node {
	if head == nil {
		return nil
	}
	if head.Val == v {
		return head.Next // the second node (or nil) becomes the head
	}
	// Look one node ahead so we can re-point prev.Next past the match.
	for prev := head; prev.Next != nil; prev = prev.Next {
		if prev.Next.Val == v {
			prev.Next = prev.Next.Next
			break
		}
	}
	return head
}

// DeleteAll removes every node holding v and returns the head of the
// resulting list.
func DeleteAll(head *Node, v int) *Node {
	// Drop matching nodes at the front first: each one changes the head.
	for head != nil && head.Val == v {
		head = head.Next
	}
	if head == nil {
		return nil
	}
	// head now holds a value we keep. Unlink matches after it; advance prev
	// only when its next node survives, so consecutive matches are all removed.
	prev := head
	for prev.Next != nil {
		if prev.Next.Val == v {
			prev.Next = prev.Next.Next
		} else {
			prev = prev.Next
		}
	}
	return head
}

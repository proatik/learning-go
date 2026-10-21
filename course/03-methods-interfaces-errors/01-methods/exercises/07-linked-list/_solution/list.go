package list

import (
	"strconv"
	"strings"
)

// List is a node of a singly linked list of ints. A nil *List is the empty
// list, and every method accepts a nil receiver.
type List struct {
	Val  int
	Next *List
}

// FromSlice builds a list holding vals in order (nil for no values).
func FromSlice(vals []int) *List {
	var head *List
	for i := len(vals) - 1; i >= 0; i-- {
		head = head.Push(vals[i])
	}
	return head
}

// Len returns the number of nodes.
func (l *List) Len() int {
	n := 0
	for ; l != nil; l = l.Next {
		n++
	}
	return n
}

// Sum returns the sum of all values.
func (l *List) Sum() int {
	total := 0
	for ; l != nil; l = l.Next {
		total += l.Val
	}
	return total
}

// Contains reports whether v is in the list.
func (l *List) Contains(v int) bool {
	for ; l != nil; l = l.Next {
		if l.Val == v {
			return true
		}
	}
	return false
}

// Values returns the values in order, or nil for the empty list.
func (l *List) Values() []int {
	var vals []int
	for ; l != nil; l = l.Next {
		vals = append(vals, l.Val)
	}
	return vals
}

// String formats the list like "[1 -> 2 -> 3]" ("[]" when empty).
func (l *List) String() string {
	parts := make([]string, 0, l.Len())
	for ; l != nil; l = l.Next {
		parts = append(parts, strconv.Itoa(l.Val))
	}
	return "[" + strings.Join(parts, " -> ") + "]"
}

// Push adds v at the front and returns the new head.
func (l *List) Push(v int) *List {
	return &List{Val: v, Next: l}
}

// Append adds v at the end and returns the head of the list.
func (l *List) Append(v int) *List {
	node := &List{Val: v}
	if l == nil {
		return node
	}
	last := l
	for last.Next != nil {
		last = last.Next
	}
	last.Next = node
	return l
}

// Reverse reverses the list in place and returns the new head.
func (l *List) Reverse() *List {
	var prev *List
	for l != nil {
		next := l.Next
		l.Next = prev
		prev, l = l, next
	}
	return prev
}

// RemoveAll unlinks every node holding v and returns the new head.
func (l *List) RemoveAll(v int) *List {
	for l != nil && l.Val == v {
		l = l.Next
	}
	for cur := l; cur != nil; {
		if cur.Next != nil && cur.Next.Val == v {
			cur.Next = cur.Next.Next
		} else {
			cur = cur.Next
		}
	}
	return l
}

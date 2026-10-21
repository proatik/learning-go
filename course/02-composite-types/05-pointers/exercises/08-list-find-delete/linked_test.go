package linked

import (
	"slices"
	"testing"
)

// build creates a list without using any of the functions under test.
func build(vals ...int) *Node {
	var head *Node
	for i := len(vals) - 1; i >= 0; i-- {
		head = &Node{Val: vals[i], Next: head}
	}
	return head
}

// values reads a list, stopping after 100 nodes so a broken list can't hang.
func values(head *Node) []int {
	var out []int
	for n := head; n != nil && len(out) < 100; n = n.Next {
		out = append(out, n.Val)
	}
	return out
}

// nodes returns the list's nodes in order.
func nodes(head *Node) []*Node {
	var out []*Node
	for n := head; n != nil && len(out) < 100; n = n.Next {
		out = append(out, n)
	}
	return out
}

func TestFind(t *testing.T) {
	head := build(1, 2, 3, 2)
	ns := nodes(head)
	tests := []struct {
		v    int
		want *Node
	}{
		{1, ns[0]},
		{2, ns[1]}, // first match only
		{3, ns[2]},
		{9, nil},
	}
	for _, tt := range tests {
		if got := Find(head, tt.v); got != tt.want {
			t.Errorf("Find([1 2 3 2], %d) returned the wrong node (got %p, want %p)", tt.v, got, tt.want)
		}
	}
	if got := Find(nil, 1); got != nil {
		t.Errorf("Find(nil, 1) = %p, want nil", got)
	}
}

func TestInsertAfter(t *testing.T) {
	head := build(1, 3)
	InsertAfter(head, 2)
	if got := values(head); !slices.Equal(got, []int{1, 2, 3}) {
		t.Errorf("InsertAfter(first of [1 3], 2): list = %v, want [1 2 3]", got)
	}
	tail := head
	for tail.Next != nil {
		tail = tail.Next
	}
	InsertAfter(tail, 4)
	if got := values(head); !slices.Equal(got, []int{1, 2, 3, 4}) {
		t.Errorf("InsertAfter(tail, 4): list = %v, want [1 2 3 4]", got)
	}
	single := build(5)
	InsertAfter(single, 6)
	if got := values(single); !slices.Equal(got, []int{5, 6}) {
		t.Errorf("InsertAfter on single node: list = %v, want [5 6]", got)
	}
}

func TestInsertAfterNil(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("InsertAfter(nil, 1) panicked: %v", r)
		}
	}()
	InsertAfter(nil, 1)
}

func TestDelete(t *testing.T) {
	tests := []struct {
		in   []int
		v    int
		want []int
	}{
		{[]int{1, 2, 3}, 1, []int{2, 3}},
		{[]int{1, 2, 3}, 2, []int{1, 3}},
		{[]int{1, 2, 3}, 3, []int{1, 2}},
		{[]int{1, 2, 2}, 2, []int{1, 2}},
		{[]int{2, 2}, 2, []int{2}},
		{[]int{1, 2, 3}, 9, []int{1, 2, 3}},
		{[]int{7}, 7, nil},
		{[]int{7}, 8, []int{7}},
		{nil, 1, nil},
	}
	for _, tt := range tests {
		got := values(Delete(build(tt.in...), tt.v))
		if !slices.Equal(got, tt.want) {
			t.Errorf("Delete(%v, %d) = %v, want %v", tt.in, tt.v, got, tt.want)
		}
	}
}

func TestDeleteKeepsNodes(t *testing.T) {
	head := build(1, 2, 3)
	ns := nodes(head)
	got := nodes(Delete(head, 2))
	if len(got) != 2 || got[0] != ns[0] || got[1] != ns[2] {
		t.Errorf("Delete([1 2 3], 2) must keep the original first and third nodes")
	}
	head = build(1, 2, 3)
	second := head.Next
	if Delete(head, 1) != second {
		t.Errorf("Delete([1 2 3], 1) must return the original second node as the new head")
	}
}

func TestDeleteAll(t *testing.T) {
	tests := []struct {
		in   []int
		v    int
		want []int
	}{
		{[]int{2, 1, 2, 2}, 2, []int{1}},
		{[]int{2, 2}, 2, nil},
		{[]int{1, 2, 3}, 2, []int{1, 3}},
		{[]int{1, 2, 2, 2, 3, 2}, 2, []int{1, 3}},
		{[]int{1, 2, 3}, 9, []int{1, 2, 3}},
		{[]int{5, 5, 6, 5}, 5, []int{6}},
		{nil, 1, nil},
	}
	for _, tt := range tests {
		got := values(DeleteAll(build(tt.in...), tt.v))
		if !slices.Equal(got, tt.want) {
			t.Errorf("DeleteAll(%v, %d) = %v, want %v", tt.in, tt.v, got, tt.want)
		}
	}
}

func TestDeleteAllKeepsNodes(t *testing.T) {
	head := build(4, 1, 4, 2)
	ns := nodes(head)
	got := nodes(DeleteAll(head, 4))
	if len(got) != 2 || got[0] != ns[1] || got[1] != ns[3] {
		t.Errorf("DeleteAll([4 1 4 2], 4) must keep the original nodes holding 1 and 2")
	}
}

package list

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

// values reads a list without using any of the functions under test.
// It stops after 100 nodes so a broken list can't hang the test.
func values(head *Node) []int {
	var out []int
	for n := head; n != nil && len(out) < 100; n = n.Next {
		out = append(out, n.Val)
	}
	return out
}

func TestPush(t *testing.T) {
	head := Push(nil, 2)
	if head == nil {
		t.Fatalf("Push(nil, 2) = nil, want a node")
	}
	if head.Val != 2 || head.Next != nil {
		t.Errorf("Push(nil, 2) = {%d %p}, want {2 nil}", head.Val, head.Next)
	}
	old := head
	head = Push(head, 1)
	if head == nil {
		t.Fatalf("Push(list, 1) = nil, want a node")
	}
	if head.Val != 1 || head.Next != old {
		t.Errorf("Push(list, 1): new head must hold 1 and link to the old head")
	}
	if got := values(head); !slices.Equal(got, []int{1, 2}) {
		t.Errorf("Push(Push(nil, 2), 1) = %v, want [1 2]", got)
	}
}

func TestFromSlice(t *testing.T) {
	tests := [][]int{{1, 2, 3}, {7}, {5, 5, 0, -1}}
	for _, vals := range tests {
		if got := values(FromSlice(vals)); !slices.Equal(got, vals) {
			t.Errorf("FromSlice(%v) holds %v, want %v", vals, got, vals)
		}
	}
	if got := FromSlice(nil); got != nil {
		t.Errorf("FromSlice(nil) = %v, want nil", values(got))
	}
	if got := FromSlice([]int{}); got != nil {
		t.Errorf("FromSlice([]) = %v, want nil", values(got))
	}
}

func TestLen(t *testing.T) {
	tests := []struct {
		head *Node
		want int
	}{
		{nil, 0},
		{build(1), 1},
		{build(1, 2), 2},
		{build(4, 3, 2, 1, 0), 5},
	}
	for _, tt := range tests {
		if got := Len(tt.head); got != tt.want {
			t.Errorf("Len(%v) = %d, want %d", values(tt.head), got, tt.want)
		}
	}
}

func TestToSlice(t *testing.T) {
	tests := [][]int{{1}, {1, 2, 3}, {3, 2, 1, 0}}
	for _, vals := range tests {
		if got := ToSlice(build(vals...)); !slices.Equal(got, vals) {
			t.Errorf("ToSlice(list of %v) = %v, want %v", vals, got, vals)
		}
	}
	if got := ToSlice(nil); got != nil {
		t.Errorf("ToSlice(nil) = %#v, want nil", got)
	}
}

func TestReverse(t *testing.T) {
	tests := []struct {
		in, want []int
	}{
		{[]int{1}, []int{1}},
		{[]int{1, 2}, []int{2, 1}},
		{[]int{1, 2, 3}, []int{3, 2, 1}},
		{[]int{1, 2, 3, 4, 5}, []int{5, 4, 3, 2, 1}},
	}
	for _, tt := range tests {
		if got := values(Reverse(build(tt.in...))); !slices.Equal(got, tt.want) {
			t.Errorf("Reverse(%v) = %v, want %v", tt.in, got, tt.want)
		}
	}
	if got := Reverse(nil); got != nil {
		t.Errorf("Reverse(nil) = %v, want nil", values(got))
	}
}

func TestReverseInPlace(t *testing.T) {
	head := build(1, 2, 3)
	tail := head.Next.Next
	got := Reverse(head)
	if got != tail {
		t.Errorf("Reverse must reuse the nodes: the old tail node should be the new head")
	}
	if head.Next != nil {
		t.Errorf("after Reverse, the old head is the tail and its Next must be nil")
	}
}

func TestRoundTrip(t *testing.T) {
	vals := []int{9, 8, 7, 6}
	head := FromSlice(vals)
	if Len(head) != 4 {
		t.Errorf("Len(FromSlice(%v)) = %d, want 4", vals, Len(head))
	}
	head = Reverse(Reverse(head))
	if got := ToSlice(head); !slices.Equal(got, vals) {
		t.Errorf("ToSlice(Reverse(Reverse(FromSlice(%v)))) = %v, want %v", vals, got, vals)
	}
}

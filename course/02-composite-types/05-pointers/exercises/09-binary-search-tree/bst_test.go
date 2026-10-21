package bst

import (
	"slices"
	"testing"
)

// sample builds this tree by hand, without using Insert:
//
//	    5
//	   / \
//	  3   8
//	 / \   \
//	1   4   9
//	         \
//	         12
func sample() *Tree {
	return &Tree{Key: 5,
		Left: &Tree{Key: 3,
			Left:  &Tree{Key: 1},
			Right: &Tree{Key: 4},
		},
		Right: &Tree{Key: 8,
			Right: &Tree{Key: 9,
				Right: &Tree{Key: 12},
			},
		},
	}
}

// keys reads a tree in order without using the functions under test.
func keys(t *Tree) []int {
	if t == nil {
		return nil
	}
	return append(append(keys(t.Left), t.Key), keys(t.Right)...)
}

func TestInsertIntoEmpty(t *testing.T) {
	root := Insert(nil, 7)
	if root == nil {
		t.Fatalf("Insert(nil, 7) = nil, want a root node")
	}
	if root.Key != 7 || root.Left != nil || root.Right != nil {
		t.Errorf("Insert(nil, 7) = {%d %p %p}, want a leaf holding 7", root.Key, root.Left, root.Right)
	}
}

func TestInsertShape(t *testing.T) {
	var root *Tree
	for _, k := range []int{5, 3, 8, 1, 4} {
		root = Insert(root, k)
		if root == nil {
			t.Fatalf("Insert returned nil root after inserting %d", k)
		}
	}
	if root.Key != 5 {
		t.Fatalf("root.Key = %d, want 5 (the first key inserted stays the root)", root.Key)
	}
	if root.Left == nil || root.Left.Key != 3 || root.Right == nil || root.Right.Key != 8 {
		t.Fatalf("children of 5 should be 3 (left) and 8 (right)")
	}
	l := root.Left
	if l.Left == nil || l.Left.Key != 1 || l.Right == nil || l.Right.Key != 4 {
		t.Errorf("children of 3 should be 1 (left) and 4 (right)")
	}
	if root.Right.Left != nil || root.Right.Right != nil {
		t.Errorf("8 should be a leaf")
	}
}

func TestInsertKeepsRootAndNodes(t *testing.T) {
	root := sample()
	left := root.Left
	if got := Insert(root, 6); got != root {
		t.Fatalf("Insert into a non-empty tree must return the same root pointer")
	}
	if root.Left != left {
		t.Errorf("Insert(6) replaced the left subtree; existing nodes must stay")
	}
	if root.Right.Left == nil || root.Right.Left.Key != 6 {
		t.Errorf("Insert(6) should add 6 as the left child of 8")
	}
}

func TestInsertDuplicates(t *testing.T) {
	var root *Tree
	for _, k := range []int{2, 1, 2, 3, 1, 3, 2} {
		root = Insert(root, k)
	}
	if got := keys(root); !slices.Equal(got, []int{1, 2, 3}) {
		t.Errorf("inserting 2 1 2 3 1 3 2: in-order keys = %v, want [1 2 3]", got)
	}
}

func TestInsertSorted(t *testing.T) {
	want := []int{50, 20, 70, 10, 30, 60, 80, 25, 65}
	var root *Tree
	for _, k := range want {
		root = Insert(root, k)
	}
	slices.Sort(want)
	if got := keys(root); !slices.Equal(got, want) {
		t.Errorf("in-order keys after inserts = %v, want %v", got, want)
	}
}

func TestContains(t *testing.T) {
	root := sample()
	for _, k := range []int{5, 3, 8, 1, 4, 9, 12} {
		if !Contains(root, k) {
			t.Errorf("Contains(sample, %d) = false, want true", k)
		}
	}
	for _, k := range []int{0, 2, 6, 7, 10, 13, -1} {
		if Contains(root, k) {
			t.Errorf("Contains(sample, %d) = true, want false", k)
		}
	}
	if Contains(nil, 0) {
		t.Errorf("Contains(nil, 0) = true, want false")
	}
}

func TestInOrder(t *testing.T) {
	if got, want := InOrder(sample()), []int{1, 3, 4, 5, 8, 9, 12}; !slices.Equal(got, want) {
		t.Errorf("InOrder(sample) = %v, want %v", got, want)
	}
	if got := InOrder(&Tree{Key: 42}); !slices.Equal(got, []int{42}) {
		t.Errorf("InOrder(single 42) = %v, want [42]", got)
	}
	if got := InOrder(nil); got != nil {
		t.Errorf("InOrder(nil) = %#v, want nil", got)
	}
}

func TestHeight(t *testing.T) {
	tests := []struct {
		name string
		t    *Tree
		want int
	}{
		{"empty", nil, 0},
		{"single", &Tree{Key: 1}, 1},
		{"left only", &Tree{Key: 2, Left: &Tree{Key: 1}}, 2},
		{"sample", sample(), 4},
		{"chain", &Tree{Key: 1, Right: &Tree{Key: 2, Right: &Tree{Key: 3, Right: &Tree{Key: 4, Right: &Tree{Key: 5}}}}}, 5},
	}
	for _, tt := range tests {
		if got := Height(tt.t); got != tt.want {
			t.Errorf("Height(%s) = %d, want %d", tt.name, got, tt.want)
		}
	}
}

func TestMin(t *testing.T) {
	if got, ok := Min(sample()); !ok || got != 1 {
		t.Errorf("Min(sample) = (%d, %v), want (1, true)", got, ok)
	}
	if got, ok := Min(&Tree{Key: 7, Right: &Tree{Key: 9}}); !ok || got != 7 {
		t.Errorf("Min(7 with right child 9) = (%d, %v), want (7, true)", got, ok)
	}
	if got, ok := Min(&Tree{Key: -3}); !ok || got != -3 {
		t.Errorf("Min(single -3) = (%d, %v), want (-3, true)", got, ok)
	}
	if got, ok := Min(nil); ok || got != 0 {
		t.Errorf("Min(nil) = (%d, %v), want (0, false)", got, ok)
	}
}

func TestEverythingTogether(t *testing.T) {
	var root *Tree
	for _, k := range []int{40, 20, 60, 10, 30, 50, 70, 5} {
		root = Insert(root, k)
	}
	if got, want := InOrder(root), []int{5, 10, 20, 30, 40, 50, 60, 70}; !slices.Equal(got, want) {
		t.Errorf("InOrder = %v, want %v", got, want)
	}
	if got := Height(root); got != 4 {
		t.Errorf("Height = %d, want 4", got)
	}
	if got, ok := Min(root); !ok || got != 5 {
		t.Errorf("Min = (%d, %v), want (5, true)", got, ok)
	}
	if !Contains(root, 30) || Contains(root, 35) {
		t.Errorf("Contains(30) should be true and Contains(35) false")
	}
}

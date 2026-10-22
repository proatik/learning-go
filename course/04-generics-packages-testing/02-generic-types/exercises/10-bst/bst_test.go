package bst

import (
	"cmp"
	"fmt"
	"slices"
	"testing"
)

// build makes a tree without relying on Insert or From.
func build[T cmp.Ordered](vals ...T) *Tree[T] {
	t := &Tree[T]{}
	for _, v := range vals {
		link := &t.root
		for *link != nil && (*link).val != v {
			if v < (*link).val {
				link = &(*link).left
			} else {
				link = &(*link).right
			}
		}
		if *link == nil {
			*link = &node[T]{val: v}
			t.size++
		}
	}
	return t
}

// shape draws a subtree as "v(left,right)", "v" for a leaf and "-" for nil.
func shape[T cmp.Ordered](n *node[T]) string {
	if n == nil {
		return "-"
	}
	if n.left == nil && n.right == nil {
		return fmt.Sprint(n.val)
	}
	return fmt.Sprintf("%v(%s,%s)", n.val, shape(n.left), shape(n.right))
}

// valid reports whether every node respects the search-tree ordering.
func valid[T cmp.Ordered](n *node[T], lo, hi *T) bool {
	if n == nil {
		return true
	}
	if (lo != nil && n.val <= *lo) || (hi != nil && n.val >= *hi) {
		return false
	}
	return valid(n.left, lo, &n.val) && valid(n.right, &n.val, hi)
}

func TestInsertShape(t *testing.T) {
	tests := []struct {
		vals []int
		want string
	}{
		{[]int{5, 3, 8, 1, 4}, "5(3(1,4),8)"},
		{[]int{1, 2, 3}, "1(-,2(-,3))"},
		{[]int{3, 2, 1}, "3(2(1,-),-)"},
		{[]int{2, 1, 3, 2, 1}, "2(1,3)"},
		{[]int{7}, "7"},
	}
	for _, tt := range tests {
		var tr Tree[int]
		for _, v := range tt.vals {
			tr.Insert(v)
		}
		if got := shape(tr.root); got != tt.want {
			t.Errorf("Insert %v: tree = %s, want %s", tt.vals, got, tt.want)
		}
	}
}

func TestInsertResult(t *testing.T) {
	var tr Tree[string]
	for _, tt := range []struct {
		v    string
		want bool
	}{{"m", true}, {"c", true}, {"m", false}, {"x", true}, {"c", false}, {"", true}} {
		if got := tr.Insert(tt.v); got != tt.want {
			t.Errorf("Insert(%q) = %v, want %v", tt.v, got, tt.want)
		}
	}
	if got := tr.Len(); got != 4 {
		t.Errorf("after 4 distinct inserts, Len() = %d, want 4", got)
	}
}

func TestFrom(t *testing.T) {
	tr := From(5, 3, 8, 1, 4, 3)
	if tr == nil {
		t.Fatalf("From(5, 3, 8, 1, 4, 3) = nil, want a tree")
	}
	if got, want := shape(tr.root), "5(3(1,4),8)"; got != want {
		t.Errorf("From(5, 3, 8, 1, 4, 3) = %s, want %s", got, want)
	}
	if got := tr.Len(); got != 5 {
		t.Errorf("From(5, 3, 8, 1, 4, 3).Len() = %d, want 5", got)
	}
	if empty := From[int](); empty == nil || empty.Len() != 0 {
		t.Errorf("From() should be an empty, non-nil tree")
	}
}

func TestContains(t *testing.T) {
	tr := build(5, 3, 8, 1, 4)
	for _, tt := range []struct {
		v    int
		want bool
	}{{5, true}, {1, true}, {4, true}, {8, true}, {7, false}, {0, false}, {9, false}, {2, false}} {
		if got := tr.Contains(tt.v); got != tt.want {
			t.Errorf("[1 3 4 5 8].Contains(%d) = %v, want %v", tt.v, got, tt.want)
		}
	}
	var empty Tree[int]
	if empty.Contains(0) {
		t.Errorf("empty.Contains(0) = true, want false")
	}
}

func TestInOrder(t *testing.T) {
	tests := []struct {
		tree *Tree[string]
		want []string
	}{
		{build("pear", "fig", "apple", "été", "Zebra"), []string{"Zebra", "apple", "fig", "pear", "été"}},
		{build("b"), []string{"b"}},
		{build[string](), []string{}},
	}
	for _, tt := range tests {
		if got := tt.tree.InOrder(); !slices.Equal(got, tt.want) {
			t.Errorf("InOrder() of %s = %q, want %q", shape(tt.tree.root), got, tt.want)
		}
	}
}

func TestHeight(t *testing.T) {
	tests := []struct {
		vals []int
		want int
	}{
		{nil, 0},
		{[]int{1}, 1},
		{[]int{5, 3, 8, 1, 4}, 3},
		{[]int{1, 2, 3, 4, 5}, 5},
		{[]int{4, 2, 6, 1, 3, 5, 7}, 3},
	}
	for _, tt := range tests {
		tr := build(tt.vals...)
		if got := tr.Height(); got != tt.want {
			t.Errorf("Height() of %s = %d, want %d", shape(tr.root), got, tt.want)
		}
	}
}

func TestMinMax(t *testing.T) {
	tr := build(5.5, -3, 8, 1, 4)
	if v, ok := tr.Min(); v != -3 || !ok {
		t.Errorf("Min() = (%v, %v), want (-3, true)", v, ok)
	}
	if v, ok := tr.Max(); v != 8 || !ok {
		t.Errorf("Max() = (%v, %v), want (8, true)", v, ok)
	}
	var empty Tree[float64]
	if v, ok := empty.Min(); v != 0 || ok {
		t.Errorf("empty.Min() = (%v, %v), want (0, false)", v, ok)
	}
	if v, ok := empty.Max(); v != 0 || ok {
		t.Errorf("empty.Max() = (%v, %v), want (0, false)", v, ok)
	}
}

func TestFloorCeil(t *testing.T) {
	tr := build(50, 30, 70, 20, 40, 60, 80)
	tests := []struct {
		v       int
		floor   int
		floorOK bool
		ceil    int
		ceilOK  bool
	}{
		{50, 50, true, 50, true},
		{45, 40, true, 50, true},
		{55, 50, true, 60, true},
		{35, 30, true, 40, true},
		{65, 60, true, 70, true},
		{19, 0, false, 20, true},
		{81, 80, true, 0, false},
		{20, 20, true, 20, true},
		{80, 80, true, 80, true},
	}
	for _, tt := range tests {
		if v, ok := tr.Floor(tt.v); v != tt.floor || ok != tt.floorOK {
			t.Errorf("Floor(%d) = (%d, %v), want (%d, %v)", tt.v, v, ok, tt.floor, tt.floorOK)
		}
		if v, ok := tr.Ceil(tt.v); v != tt.ceil || ok != tt.ceilOK {
			t.Errorf("Ceil(%d) = (%d, %v), want (%d, %v)", tt.v, v, ok, tt.ceil, tt.ceilOK)
		}
	}
	var empty Tree[int]
	if v, ok := empty.Floor(1); v != 0 || ok {
		t.Errorf("empty.Floor(1) = (%d, %v), want (0, false)", v, ok)
	}
	if v, ok := empty.Ceil(1); v != 0 || ok {
		t.Errorf("empty.Ceil(1) = (%d, %v), want (0, false)", v, ok)
	}
}

func TestDeleteShape(t *testing.T) {
	tests := []struct {
		vals []int
		del  int
		want string
	}{
		{[]int{5, 3, 8, 1, 4}, 1, "5(3(-,4),8)"},                // leaf
		{[]int{5, 3, 8, 1, 4}, 8, "5(3(1,4),-)"},                // leaf on the right
		{[]int{5, 3, 1}, 3, "5(1,-)"},                           // one child (left)
		{[]int{5, 3, 4}, 3, "5(4,-)"},                           // one child (right)
		{[]int{5, 3, 8, 1, 4}, 3, "5(4(1,-),8)"},                // two children
		{[]int{5, 3, 8, 1, 4}, 5, "8(3(1,4),-)"},                // root, two children
		{[]int{50, 30, 70, 60, 80, 65}, 50, "60(30,70(65,80))"}, // successor has a right child
		{[]int{2, 1}, 2, "1"},                                   // root, one child
		{[]int{7}, 7, "-"},                                      // only node
	}
	for _, tt := range tests {
		tr := build(tt.vals...)
		before := shape(tr.root)
		if !tr.Delete(tt.del) {
			t.Errorf("Delete(%d) on %s = false, want true", tt.del, before)
			continue
		}
		if got := shape(tr.root); got != tt.want {
			t.Errorf("Delete(%d) on %s: tree = %s, want %s", tt.del, before, got, tt.want)
		}
		if got := tr.Len(); got != len(tt.vals)-1 {
			t.Errorf("Delete(%d) on %s: Len() = %d, want %d", tt.del, before, got, len(tt.vals)-1)
		}
	}
}

func TestDeleteMissing(t *testing.T) {
	tr := build(5, 3, 8)
	if tr.Delete(4) {
		t.Errorf("Delete(4) on 5(3,8) = true, want false")
	}
	if got := tr.Len(); got != 3 {
		t.Errorf("after missing Delete, Len() = %d, want 3", got)
	}
	var empty Tree[int]
	if empty.Delete(1) {
		t.Errorf("empty.Delete(1) = true, want false")
	}
}

func TestMixedOperations(t *testing.T) {
	var tr Tree[int]
	want := map[int]bool{}
	seed := uint32(42)
	next := func(n int) int {
		seed = seed*1664525 + 1013904223
		return int(seed>>16) % n
	}
	for step := range 3000 {
		v := next(60)
		if next(3) == 0 {
			if got := tr.Delete(v); got != want[v] {
				t.Fatalf("step %d: Delete(%d) = %v, want %v", step, v, got, want[v])
			}
			delete(want, v)
		} else {
			if got := tr.Insert(v); got == want[v] {
				t.Fatalf("step %d: Insert(%d) = %v, want %v", step, v, got, !want[v])
			}
			want[v] = true
		}
		if !valid(tr.root, nil, nil) {
			t.Fatalf("step %d: tree breaks the search-tree ordering: %s", step, shape(tr.root))
		}
		if tr.Len() != len(want) {
			t.Fatalf("step %d: Len() = %d, want %d", step, tr.Len(), len(want))
		}
	}
	var keys []int
	for k := range want {
		keys = append(keys, k)
	}
	slices.Sort(keys)
	if got := tr.InOrder(); !slices.Equal(got, keys) {
		t.Errorf("after 3000 random operations, InOrder() = %v, want %v", got, keys)
	}
}

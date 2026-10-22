package bst

import (
	"iter"
	"slices"
	"testing"
)

// collect ranges over seq and returns its values, breaking after max values
// when max > 0 (and always after 10000, so a runaway sequence can't hang the
// test). A panic during the loop is returned instead of crashing.
func collect[V any](seq iter.Seq[V], max int) (got []V, panicked any) {
	defer func() { panicked = recover() }()
	for v := range seq {
		got = append(got, v)
		if len(got) == max || len(got) == 10000 {
			break
		}
	}
	return got, nil
}

// yieldCalls calls seq directly with a yield that returns false on call
// number stopAt, and reports how many times yield was called in total.
func yieldCalls[V any](seq iter.Seq[V], stopAt int) int {
	calls := 0
	seq(func(V) bool {
		calls++
		return calls < stopAt
	})
	return calls
}

func build(keys ...int) *Tree[int] {
	t := &Tree[int]{}
	for _, k := range keys {
		t.Insert(k)
	}
	return t
}

var sample = []int{5, 3, 8, 1, 4, 9}

func TestAll(t *testing.T) {
	tests := []struct {
		keys []int
		want []int
	}{
		{sample, []int{1, 3, 4, 5, 8, 9}},
		{[]int{1, 2, 3, 4}, []int{1, 2, 3, 4}},
		{[]int{4, 3, 2, 1}, []int{1, 2, 3, 4}},
		{[]int{7}, []int{7}},
		{[]int{2, 2, 1, 1}, []int{1, 2}},
		{nil, nil},
	}
	for _, tt := range tests {
		tr := build(tt.keys...)
		seq := tr.All()
		if seq == nil {
			t.Fatalf("All() after inserting %v = nil, want a non-nil iterator", tt.keys)
		}
		got, p := collect(seq, -1)
		if p != nil {
			t.Fatalf("range over All() after inserting %v panicked: %v", tt.keys, p)
		}
		if !slices.Equal(got, tt.want) {
			t.Errorf("All() after inserting %v = %v, want %v", tt.keys, got, tt.want)
		}
	}
}

func TestBackward(t *testing.T) {
	tests := []struct {
		keys []int
		want []int
	}{
		{sample, []int{9, 8, 5, 4, 3, 1}},
		{[]int{1, 2, 3}, []int{3, 2, 1}},
		{[]int{7}, []int{7}},
		{nil, nil},
	}
	for _, tt := range tests {
		got, p := collect(build(tt.keys...).Backward(), -1)
		if p != nil {
			t.Fatalf("range over Backward() after inserting %v panicked: %v", tt.keys, p)
		}
		if !slices.Equal(got, tt.want) {
			t.Errorf("Backward() after inserting %v = %v, want %v", tt.keys, got, tt.want)
		}
	}
}

func TestStrings(t *testing.T) {
	tr := &Tree[string]{}
	for _, w := range []string{"pear", "apple", "zebra", "école", "fig"} {
		tr.Insert(w)
	}
	got, _ := collect(tr.All(), -1)
	if want := []string{"apple", "fig", "pear", "zebra", "école"}; !slices.Equal(got, want) {
		t.Errorf("All() on string keys = %q, want %q", got, want)
	}
}

func TestRange(t *testing.T) {
	tr := build(sample...)
	tests := []struct {
		lo, hi int
		want   []int
	}{
		{3, 8, []int{3, 4, 5, 8}},
		{0, 100, []int{1, 3, 4, 5, 8, 9}},
		{4, 4, []int{4}},
		{6, 7, nil},
		{2, 3, []int{3}},
		{9, 20, []int{9}},
		{-5, 1, []int{1}},
		{8, 3, nil},
	}
	for _, tt := range tests {
		got, p := collect(tr.Range(tt.lo, tt.hi), -1)
		if p != nil {
			t.Fatalf("range over Range(%d, %d) panicked: %v", tt.lo, tt.hi, p)
		}
		if !slices.Equal(got, tt.want) {
			t.Errorf("Range(%d, %d) on %v = %v, want %v", tt.lo, tt.hi, sample, got, tt.want)
		}
	}
	empty, _ := collect((&Tree[int]{}).Range(0, 10), -1)
	if len(empty) != 0 {
		t.Errorf("Range(0, 10) on an empty tree = %v, want []", empty)
	}
}

// TestRangeSkipsSubtrees plants keys in places a correct walk never looks.
// They break the search-tree order, so they only show up if Range walks
// into a subtree that can't hold keys in [lo, hi].
func TestRangeSkipsSubtrees(t *testing.T) {
	tr := build(50, 20, 80, 10, 30, 70, 90)
	tr.root.left.left.left = &node[int]{key: 75} // under 10, inside [60, 80]
	got, p := collect(tr.Range(60, 80), -1)
	if p != nil {
		t.Fatalf("range over Range(60, 80) panicked: %v", p)
	}
	if want := []int{70, 80}; !slices.Equal(got, want) {
		t.Errorf("Range(60, 80) = %v, want %v (don't go left when n.key < lo)", got, want)
	}

	tr = build(50, 20, 80, 10, 30, 70, 90)
	tr.root.right.right.right = &node[int]{key: 15} // under 90, inside [10, 30]
	got, _ = collect(tr.Range(10, 30), -1)
	if want := []int{10, 20, 30}; !slices.Equal(got, want) {
		t.Errorf("Range(10, 30) = %v, want %v (don't go right when n.key > hi)", got, want)
	}
}

func TestBreak(t *testing.T) {
	keys := make([]int, 0, 200)
	for i := range 200 {
		keys = append(keys, (i*73)%200) // a bushy, unbalanced tree
	}
	tr := build(keys...)
	for _, stopAt := range []int{1, 2, 7, 100} {
		got, p := collect(tr.All(), stopAt)
		if p != nil {
			t.Fatalf("break after %d keys of All() panicked: %v (every level must stop when yield returns false)", stopAt, p)
		}
		if len(got) != stopAt || got[0] != 0 || got[len(got)-1] != stopAt-1 {
			t.Errorf("first %d keys of All() = %v, want 0..%d", stopAt, got, stopAt-1)
		}
		if n := yieldCalls(tr.All(), stopAt); n != stopAt {
			t.Errorf("All() with yield returning false on call %d: yield calls = %d, want %d", stopAt, n, stopAt)
		}
		if n := yieldCalls(tr.Backward(), stopAt); n != stopAt {
			t.Errorf("Backward() with yield returning false on call %d: yield calls = %d, want %d", stopAt, n, stopAt)
		}
		if n := yieldCalls(tr.Range(10, 150), stopAt); n != stopAt {
			t.Errorf("Range(10, 150) with yield returning false on call %d: yield calls = %d, want %d", stopAt, n, stopAt)
		}
	}
	got, p := collect(tr.Backward(), 3)
	if p != nil {
		t.Fatalf("break after 3 keys of Backward() panicked: %v", p)
	}
	if want := []int{199, 198, 197}; !slices.Equal(got, want) {
		t.Errorf("first 3 keys of Backward() = %v, want %v", got, want)
	}
	got, p = collect(tr.Range(50, 60), 2)
	if p != nil {
		t.Fatalf("break after 2 keys of Range(50, 60) panicked: %v", p)
	}
	if want := []int{50, 51}; !slices.Equal(got, want) {
		t.Errorf("first 2 keys of Range(50, 60) = %v, want %v", got, want)
	}
}

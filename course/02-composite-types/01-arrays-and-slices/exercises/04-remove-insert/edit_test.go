package edit

import (
	"slices"
	"testing"
)

func TestRemoveAt(t *testing.T) {
	tests := []struct {
		in   []int
		i    int
		want []int
	}{
		{[]int{10, 20, 30, 40}, 1, []int{10, 30, 40}},
		{[]int{10, 20, 30, 40}, 0, []int{20, 30, 40}},
		{[]int{10, 20, 30, 40}, 3, []int{10, 20, 30}},
		{[]int{10}, 0, []int{}},
		{[]int{1, 1, 1}, 1, []int{1, 1}},
		{[]int{5, 6}, 1, []int{5}},
	}
	for _, tt := range tests {
		got := RemoveAt(slices.Clone(tt.in), tt.i)
		if !slices.Equal(got, tt.want) {
			t.Errorf("RemoveAt(%v, %d) = %v, want %v", tt.in, tt.i, got, tt.want)
		}
	}
}

func TestRemoveAtInPlace(t *testing.T) {
	s := make([]int, 4, 10)
	copy(s, []int{1, 2, 3, 4})
	got := RemoveAt(s, 1)
	if len(got) != 3 {
		t.Fatalf("RemoveAt([1 2 3 4], 1) = %v, want 3 elements", got)
	}
	if &got[0] != &s[0] || cap(got) != cap(s) {
		t.Errorf("RemoveAt allocated a new slice; it should reuse the input's backing array")
	}
}

func TestRemoveUnordered(t *testing.T) {
	tests := []struct {
		in   []int
		i    int
		want []int
	}{
		{[]int{10, 20, 30, 40}, 0, []int{40, 20, 30}},
		{[]int{10, 20, 30, 40}, 1, []int{10, 40, 30}},
		{[]int{10, 20, 30, 40}, 3, []int{10, 20, 30}},
		{[]int{7}, 0, []int{}},
	}
	for _, tt := range tests {
		got := RemoveUnordered(slices.Clone(tt.in), tt.i)
		if !slices.Equal(got, tt.want) {
			t.Errorf("RemoveUnordered(%v, %d) = %v, want %v", tt.in, tt.i, got, tt.want)
		}
	}
}

func TestInsertAt(t *testing.T) {
	tests := []struct {
		in   []int
		i, v int
		want []int
	}{
		{[]int{1, 2, 4}, 2, 3, []int{1, 2, 3, 4}},
		{[]int{1, 2}, 0, 0, []int{0, 1, 2}},
		{[]int{1, 2}, 2, 3, []int{1, 2, 3}},
		{[]int{1, 2}, 1, 9, []int{1, 9, 2}},
		{nil, 0, 7, []int{7}},
		{[]int{}, 0, 7, []int{7}},
		{[]int{5, 5, 5}, 1, 5, []int{5, 5, 5, 5}},
	}
	for _, tt := range tests {
		got := InsertAt(slices.Clone(tt.in), tt.i, tt.v)
		if !slices.Equal(got, tt.want) {
			t.Errorf("InsertAt(%v, %d, %d) = %v, want %v", tt.in, tt.i, tt.v, got, tt.want)
		}
	}
}

func TestInsertAtSpareCapacity(t *testing.T) {
	// With spare capacity, InsertAt must not allocate.
	s := make([]int, 3, 8)
	copy(s, []int{1, 2, 4})
	got := InsertAt(s, 2, 3)
	if !slices.Equal(got, []int{1, 2, 3, 4}) {
		t.Fatalf("InsertAt([1 2 4], 2, 3) = %v, want [1 2 3 4]", got)
	}
	if &got[0] != &s[0] {
		t.Errorf("InsertAt allocated although the input had spare capacity")
	}
}

func TestInsertThenRemove(t *testing.T) {
	var s []int
	for i := range 5 {
		s = InsertAt(s, 0, i) // build [4 3 2 1 0]
	}
	if !slices.Equal(s, []int{4, 3, 2, 1, 0}) {
		t.Fatalf("five InsertAt(s, 0, i) calls gave %v, want [4 3 2 1 0]", s)
	}
	s = RemoveAt(s, 2)
	s = RemoveAt(s, 0)
	if !slices.Equal(s, []int{3, 1, 0}) {
		t.Errorf("after removing indices 2 then 0: got %v, want [3 1 0]", s)
	}
}

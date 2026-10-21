package rotate

import (
	"slices"
	"testing"
)

func TestRotateLeft(t *testing.T) {
	tests := []struct {
		in   []int
		k    int
		want []int
	}{
		{[]int{1, 2, 3, 4, 5}, 2, []int{3, 4, 5, 1, 2}},
		{[]int{1, 2, 3, 4, 5}, 1, []int{2, 3, 4, 5, 1}},
		{[]int{1, 2, 3, 4, 5}, 4, []int{5, 1, 2, 3, 4}},
		{[]int{1, 2, 3, 4, 5}, 0, []int{1, 2, 3, 4, 5}},
		{[]int{1, 2, 3, 4, 5}, 5, []int{1, 2, 3, 4, 5}},
		{[]int{1, 2, 3, 4, 5}, 7, []int{3, 4, 5, 1, 2}},
		{[]int{1, 2, 3, 4, 5}, -1, []int{5, 1, 2, 3, 4}},
		{[]int{1, 2, 3, 4, 5}, -7, []int{4, 5, 1, 2, 3}},
		{[]int{1, 2, 3}, 3, []int{1, 2, 3}},
		{[]int{1, 2}, 1, []int{2, 1}},
		{[]int{9}, 100, []int{9}},
		{[]int{1, 2, 3, 4, 5, 6}, 3, []int{4, 5, 6, 1, 2, 3}},
		{[]int{}, 4, []int{}},
		{nil, 4, nil},
	}
	for _, tt := range tests {
		s := slices.Clone(tt.in)
		RotateLeft(s, tt.k)
		if !slices.Equal(s, tt.want) {
			t.Errorf("after RotateLeft(%v, %d): got %v, want %v", tt.in, tt.k, s, tt.want)
		}
	}
}

func TestRotateRight(t *testing.T) {
	tests := []struct {
		in   []int
		k    int
		want []int
	}{
		{[]int{1, 2, 3, 4, 5}, 2, []int{4, 5, 1, 2, 3}},
		{[]int{1, 2, 3, 4, 5}, 1, []int{5, 1, 2, 3, 4}},
		{[]int{1, 2, 3, 4, 5}, 6, []int{5, 1, 2, 3, 4}},
		{[]int{1, 2, 3, 4, 5}, -2, []int{3, 4, 5, 1, 2}},
		{[]int{1, 2, 3, 4, 5}, 0, []int{1, 2, 3, 4, 5}},
		{nil, -3, nil},
	}
	for _, tt := range tests {
		s := slices.Clone(tt.in)
		RotateRight(s, tt.k)
		if !slices.Equal(s, tt.want) {
			t.Errorf("after RotateRight(%v, %d): got %v, want %v", tt.in, tt.k, s, tt.want)
		}
	}
}

func TestRotateRoundTrip(t *testing.T) {
	orig := []int{10, 20, 30, 40, 50, 60, 70}
	for k := -10; k <= 10; k++ {
		s := slices.Clone(orig)
		RotateLeft(s, k)
		RotateRight(s, k)
		if !slices.Equal(s, orig) {
			t.Errorf("RotateLeft(s, %d) then RotateRight(s, %d) gave %v, want %v", k, k, s, orig)
		}
	}
}

func TestRotateNoAllocations(t *testing.T) {
	s := make([]int, 1000)
	for i := range s {
		s[i] = i
	}
	allocs := testing.AllocsPerRun(50, func() {
		RotateLeft(s, 337)
		RotateRight(s, 123)
	})
	if allocs != 0 {
		t.Errorf("RotateLeft + RotateRight allocated %.0f times per run, want 0", allocs)
	}
}

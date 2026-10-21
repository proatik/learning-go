package filter

import (
	"slices"
	"testing"
)

func TestKeepEven(t *testing.T) {
	tests := []struct {
		in   []int
		want []int
	}{
		{[]int{1, 2, 3, 4, 5, 6}, []int{2, 4, 6}},
		{[]int{-3, -2, 0, 7}, []int{-2, 0}},
		{[]int{1, 3}, []int{}},
		{[]int{2, 4}, []int{2, 4}},
		{[]int{8}, []int{8}},
		{[]int{9}, []int{}},
		{[]int{}, []int{}},
		{nil, nil},
	}
	for _, tt := range tests {
		got := KeepEven(slices.Clone(tt.in))
		if !slices.Equal(got, tt.want) {
			t.Errorf("KeepEven(%v) = %v, want %v", tt.in, got, tt.want)
		}
	}
}

func TestRemoveValue(t *testing.T) {
	tests := []struct {
		in   []int
		v    int
		want []int
	}{
		{[]int{3, 1, 3, 3, 2}, 3, []int{1, 2}},
		{[]int{1, 2}, 9, []int{1, 2}},
		{[]int{5, 5, 5}, 5, []int{}},
		{[]int{0, 1, 0}, 0, []int{1}},
		{[]int{1, 2, 3, 4}, 1, []int{2, 3, 4}},
		{[]int{1, 2, 3, 4}, 4, []int{1, 2, 3}},
		{nil, 1, nil},
	}
	for _, tt := range tests {
		got := RemoveValue(slices.Clone(tt.in), tt.v)
		if !slices.Equal(got, tt.want) {
			t.Errorf("RemoveValue(%v, %d) = %v, want %v", tt.in, tt.v, got, tt.want)
		}
	}
}

func TestKeepEvenReusesArrayAndClearsTail(t *testing.T) {
	s := []int{1, 2, 3, 4, 5, 6}
	got := KeepEven(s)
	if len(got) != 3 {
		t.Fatalf("KeepEven([1 2 3 4 5 6]) = %v, want [2 4 6]", got)
	}
	if &got[0] != &s[0] {
		t.Errorf("KeepEven returned a new array; it must reuse the input's backing array")
	}
	if want := []int{2, 4, 6, 0, 0, 0}; !slices.Equal(s, want) {
		t.Errorf("original slice after KeepEven = %v, want %v (tail zeroed with clear)", s, want)
	}
}

func TestRemoveValueReusesArrayAndClearsTail(t *testing.T) {
	s := []int{7, 1, 7, 2, 7}
	got := RemoveValue(s, 7)
	if len(got) != 2 {
		t.Fatalf("RemoveValue([7 1 7 2 7], 7) = %v, want [1 2]", got)
	}
	if &got[0] != &s[0] {
		t.Errorf("RemoveValue returned a new array; it must reuse the input's backing array")
	}
	if want := []int{1, 2, 0, 0, 0}; !slices.Equal(s, want) {
		t.Errorf("original slice after RemoveValue = %v, want %v (tail zeroed with clear)", s, want)
	}
}

func TestNoAllocations(t *testing.T) {
	data := make([]int, 100)
	s := make([]int, 100)
	allocs := testing.AllocsPerRun(100, func() {
		for i := range data {
			data[i] = i
		}
		copy(s, data)
		_ = KeepEven(s)
		copy(s, data)
		_ = RemoveValue(s, 50)
	})
	if allocs != 0 {
		t.Errorf("KeepEven + RemoveValue allocated %.0f times per run, want 0", allocs)
	}
}

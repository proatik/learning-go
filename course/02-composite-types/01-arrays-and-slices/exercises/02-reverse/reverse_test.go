package reverse

import (
	"slices"
	"testing"
)

func TestReverse(t *testing.T) {
	tests := []struct {
		in   []int
		want []int
	}{
		{nil, nil},
		{[]int{}, []int{}},
		{[]int{1}, []int{1}},
		{[]int{1, 2}, []int{2, 1}},
		{[]int{1, 2, 3}, []int{3, 2, 1}},
		{[]int{1, 2, 3, 4}, []int{4, 3, 2, 1}},
		{[]int{5, 5, 1, 5}, []int{5, 1, 5, 5}},
		{[]int{-1, 0, 1, 2, 3, 4, 5}, []int{5, 4, 3, 2, 1, 0, -1}},
	}
	for _, tt := range tests {
		s := slices.Clone(tt.in)
		Reverse(s)
		if !slices.Equal(s, tt.want) {
			t.Errorf("after Reverse(%v): got %v, want %v", tt.in, s, tt.want)
		}
	}
}

func TestReverseSubslice(t *testing.T) {
	s := []int{0, 1, 2, 3, 4, 5}
	Reverse(s[1:5]) // only the middle part
	want := []int{0, 4, 3, 2, 1, 5}
	if !slices.Equal(s, want) {
		t.Errorf("after Reverse(s[1:5]) on [0 1 2 3 4 5]: got %v, want %v", s, want)
	}
}

func TestReversedArray(t *testing.T) {
	tests := []struct {
		in, want [5]int
	}{
		{[5]int{1, 2, 3, 4, 5}, [5]int{5, 4, 3, 2, 1}},
		{[5]int{}, [5]int{}},
		{[5]int{9, 0, 0, 0, 0}, [5]int{0, 0, 0, 0, 9}},
		{[5]int{1, 2, 3, 2, 1}, [5]int{1, 2, 3, 2, 1}},
	}
	for _, tt := range tests {
		orig := tt.in
		got := ReversedArray(tt.in)
		if got != tt.want {
			t.Errorf("ReversedArray(%v) = %v, want %v", orig, got, tt.want)
		}
		if tt.in != orig {
			t.Errorf("ReversedArray modified the caller's array: %v", tt.in)
		}
	}
}

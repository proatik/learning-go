package sums

import "testing"

func TestSumArray(t *testing.T) {
	tests := []struct {
		in   [5]int
		want int
	}{
		{[5]int{1, 2, 3, 4, 5}, 15},
		{[5]int{}, 0},
		{[5]int{-1, -2, -3, -4, -5}, -15},
		{[5]int{10, -10, 7, 0, 3}, 10},
		{[5]int{4: 100}, 100},
	}
	for _, tt := range tests {
		if got := SumArray(tt.in); got != tt.want {
			t.Errorf("SumArray(%v) = %d, want %d", tt.in, got, tt.want)
		}
	}
}

func TestSum(t *testing.T) {
	tests := []struct {
		in   []int
		want int
	}{
		{nil, 0},
		{[]int{}, 0},
		{[]int{42}, 42},
		{[]int{4, -1, 7}, 10},
		{[]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}, 55},
		{[]int{-3, 3}, 0},
	}
	for _, tt := range tests {
		if got := Sum(tt.in); got != tt.want {
			t.Errorf("Sum(%v) = %d, want %d", tt.in, got, tt.want)
		}
	}
}

func TestLargest(t *testing.T) {
	tests := []struct {
		in     []int
		want   int
		wantOK bool
	}{
		{nil, 0, false},
		{[]int{}, 0, false},
		{[]int{7}, 7, true},
		{[]int{0}, 0, true},
		{[]int{3, 9, 2}, 9, true},
		{[]int{9, 3, 2}, 9, true},
		{[]int{2, 3, 9}, 9, true},
		{[]int{-5, -2, -8}, -2, true},
		{[]int{5, 5, 5}, 5, true},
	}
	for _, tt := range tests {
		got, ok := Largest(tt.in)
		if got != tt.want || ok != tt.wantOK {
			t.Errorf("Largest(%v) = (%d, %t), want (%d, %t)", tt.in, got, ok, tt.want, tt.wantOK)
		}
	}
}

func TestSumArrayDoesNotModify(t *testing.T) {
	a := [5]int{1, 2, 3, 4, 5}
	SumArray(a)
	if a != [5]int{1, 2, 3, 4, 5} {
		t.Errorf("SumArray modified the caller's array: %v", a)
	}
}

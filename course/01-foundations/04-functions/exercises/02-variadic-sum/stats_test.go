package stats

import "testing"

func TestSumArgs(t *testing.T) {
	if got := Sum(); got != 0 {
		t.Errorf("Sum() = %d, want 0", got)
	}
	if got := Sum(4); got != 4 {
		t.Errorf("Sum(4) = %d, want 4", got)
	}
	if got := Sum(1, 2, 3); got != 6 {
		t.Errorf("Sum(1, 2, 3) = %d, want 6", got)
	}
	if got := Sum(-1, -2, 3); got != 0 {
		t.Errorf("Sum(-1, -2, 3) = %d, want 0", got)
	}
}

func TestSumSpread(t *testing.T) {
	tests := []struct {
		nums []int
		want int
	}{
		{nil, 0},
		{[]int{}, 0},
		{[]int{10, -5}, 5},
		{[]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}, 55},
		{[]int{-100}, -100},
	}
	for _, tt := range tests {
		if got := Sum(tt.nums...); got != tt.want {
			t.Errorf("Sum(%v...) = %d, want %d", tt.nums, got, tt.want)
		}
	}
}

func TestMean(t *testing.T) {
	tests := []struct {
		nums []int
		want float64
	}{
		{[]int{1, 2, 3, 4}, 2.5},
		{[]int{7}, 7},
		{[]int{7, 0}, 3.5},
		{[]int{-3, 3}, 0},
		{[]int{-1, -2}, -1.5},
		{[]int{1, 1, 2}, 4.0 / 3.0},
	}
	for _, tt := range tests {
		got, err := Mean(tt.nums...)
		if err != nil {
			t.Errorf("Mean(%v...) returned error %v, want nil", tt.nums, err)
			continue
		}
		if got != tt.want {
			t.Errorf("Mean(%v...) = %v, want %v", tt.nums, got, tt.want)
		}
	}
}

func TestMeanEmpty(t *testing.T) {
	got, err := Mean()
	if err == nil {
		t.Errorf("Mean(): err = nil, want an error")
	}
	if got != 0 {
		t.Errorf("Mean() = %v, want 0 alongside the error", got)
	}
}

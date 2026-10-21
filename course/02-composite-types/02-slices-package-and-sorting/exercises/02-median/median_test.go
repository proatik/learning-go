package median

import (
	"slices"
	"testing"
)

func TestMedian(t *testing.T) {
	tests := []struct {
		nums []float64
		want float64
	}{
		{[]float64{3, 1, 2}, 2},
		{[]float64{4, 1, 3, 2}, 2.5},
		{[]float64{7}, 7},
		{[]float64{1, 2}, 1.5},
		{[]float64{5, 5, 5, 5}, 5},
		{[]float64{-3, -1, -2}, -2},
		{[]float64{10, 1, 9, 2, 8}, 8},
		{[]float64{0.5, 0.25}, 0.375},
		{[]float64{9, 8, 7, 6, 5, 4, 3, 2, 1, 0}, 4.5},
	}
	for _, tt := range tests {
		orig := slices.Clone(tt.nums)
		got, err := Median(tt.nums)
		if err != nil {
			t.Errorf("Median(%v) returned error %v, want nil", orig, err)
			continue
		}
		if got != tt.want {
			t.Errorf("Median(%v) = %v, want %v", orig, got, tt.want)
		}
		if !slices.Equal(tt.nums, orig) {
			t.Errorf("Median modified its input: %v became %v", orig, tt.nums)
		}
	}
}

func TestMedianEmpty(t *testing.T) {
	for _, nums := range [][]float64{nil, {}} {
		got, err := Median(nums)
		if err == nil {
			t.Errorf("Median(%v): err = nil, want an error", nums)
		}
		if got != 0 {
			t.Errorf("Median(%v) = %v, want 0 alongside the error", nums, got)
		}
	}
}

func TestSpread(t *testing.T) {
	tests := []struct {
		nums []int
		want int
	}{
		{[]int{4, -2, 9}, 11},
		{[]int{5}, 0},
		{[]int{1, 2, 3}, 2},
		{[]int{3, 3, 3}, 0},
		{[]int{-10, -20}, 10},
		{[]int{100, 0, 50}, 100},
	}
	for _, tt := range tests {
		orig := slices.Clone(tt.nums)
		got, err := Spread(tt.nums)
		if err != nil {
			t.Errorf("Spread(%v) returned error %v, want nil", orig, err)
			continue
		}
		if got != tt.want {
			t.Errorf("Spread(%v) = %d, want %d", orig, got, tt.want)
		}
		if !slices.Equal(tt.nums, orig) {
			t.Errorf("Spread modified its input: %v became %v", orig, tt.nums)
		}
	}
}

func TestSpreadEmpty(t *testing.T) {
	for _, nums := range [][]int{nil, {}} {
		got, err := Spread(nums)
		if err == nil {
			t.Errorf("Spread(%v): err = nil, want an error", nums)
		}
		if got != 0 {
			t.Errorf("Spread(%v) = %d, want 0 alongside the error", nums, got)
		}
	}
}

package topk

import (
	"slices"
	"testing"
)

func TestTopK(t *testing.T) {
	tests := []struct {
		nums []int
		k    int
		want []int
	}{
		{[]int{5, 1, 9, 3, 9}, 3, []int{9, 9, 5}},
		{[]int{5, 1, 9, 3, 9}, 1, []int{9}},
		{[]int{5, 1, 9, 3, 9}, 5, []int{9, 9, 5, 3, 1}},
		{[]int{2, 1}, 10, []int{2, 1}},
		{[]int{2, 1}, 0, []int{}},
		{[]int{2, 1}, -1, []int{}},
		{[]int{-4, -1, -9}, 2, []int{-1, -4}},
		{[]int{1, 2, 3, 4}, 2, []int{4, 3}},
		{nil, 3, []int{}},
	}
	for _, tt := range tests {
		orig := slices.Clone(tt.nums)
		got := TopK(tt.nums, tt.k)
		if !slices.Equal(got, tt.want) {
			t.Errorf("TopK(%v, %d) = %v, want %v", orig, tt.k, got, tt.want)
		}
		if !slices.Equal(tt.nums, orig) {
			t.Errorf("TopK modified its input: %v became %v", orig, tt.nums)
		}
	}
}

func TestTopKDoesNotAlias(t *testing.T) {
	nums := []int{3, 2, 1}
	got := TopK(nums, 3)
	if len(got) != 3 {
		t.Fatalf("TopK(%v, 3) = %v, want 3 elements", nums, got)
	}
	got[0] = 100
	if !slices.Equal(nums, []int{3, 2, 1}) {
		t.Errorf("writing to TopK's result changed the input: %v", nums)
	}
}

func TestKthSmallest(t *testing.T) {
	tests := []struct {
		nums []int
		k    int
		want int
	}{
		{[]int{7, 2, 9, 2}, 1, 2},
		{[]int{7, 2, 9, 2}, 2, 2},
		{[]int{7, 2, 9, 2}, 3, 7},
		{[]int{7, 2, 9, 2}, 4, 9},
		{[]int{-5}, 1, -5},
		{[]int{10, 30, 20}, 2, 20},
	}
	for _, tt := range tests {
		orig := slices.Clone(tt.nums)
		got, err := KthSmallest(tt.nums, tt.k)
		if err != nil {
			t.Errorf("KthSmallest(%v, %d) returned error %v, want nil", orig, tt.k, err)
			continue
		}
		if got != tt.want {
			t.Errorf("KthSmallest(%v, %d) = %d, want %d", orig, tt.k, got, tt.want)
		}
		if !slices.Equal(tt.nums, orig) {
			t.Errorf("KthSmallest modified its input: %v became %v", orig, tt.nums)
		}
	}
}

func TestKthSmallestErrors(t *testing.T) {
	tests := []struct {
		nums []int
		k    int
	}{
		{[]int{7, 2, 9, 2}, 5},
		{[]int{7, 2, 9, 2}, 0},
		{[]int{7, 2, 9, 2}, -1},
		{nil, 1},
	}
	for _, tt := range tests {
		got, err := KthSmallest(tt.nums, tt.k)
		if err == nil {
			t.Errorf("KthSmallest(%v, %d): err = nil, want an error", tt.nums, tt.k)
		}
		if got != 0 {
			t.Errorf("KthSmallest(%v, %d) = %d, want 0 alongside the error", tt.nums, tt.k, got)
		}
	}
}

func TestKthLargestDistinct(t *testing.T) {
	tests := []struct {
		nums []int
		k    int
		want int
	}{
		{[]int{7, 2, 9, 2, 9}, 1, 9},
		{[]int{7, 2, 9, 2, 9}, 2, 7},
		{[]int{7, 2, 9, 2, 9}, 3, 2},
		{[]int{4, 4, 4}, 1, 4},
		{[]int{1, 2, 3, 4, 5}, 5, 1},
		{[]int{-1, -1, -3}, 2, -3},
	}
	for _, tt := range tests {
		orig := slices.Clone(tt.nums)
		got, err := KthLargestDistinct(tt.nums, tt.k)
		if err != nil {
			t.Errorf("KthLargestDistinct(%v, %d) returned error %v, want nil", orig, tt.k, err)
			continue
		}
		if got != tt.want {
			t.Errorf("KthLargestDistinct(%v, %d) = %d, want %d", orig, tt.k, got, tt.want)
		}
		if !slices.Equal(tt.nums, orig) {
			t.Errorf("KthLargestDistinct modified its input: %v became %v", orig, tt.nums)
		}
	}
}

func TestKthLargestDistinctErrors(t *testing.T) {
	tests := []struct {
		nums []int
		k    int
	}{
		{[]int{4, 4, 4}, 2},
		{[]int{7, 2, 9, 2, 9}, 4},
		{[]int{1, 2}, 0},
		{nil, 1},
	}
	for _, tt := range tests {
		got, err := KthLargestDistinct(tt.nums, tt.k)
		if err == nil {
			t.Errorf("KthLargestDistinct(%v, %d): err = nil, want an error", tt.nums, tt.k)
		}
		if got != 0 {
			t.Errorf("KthLargestDistinct(%v, %d) = %d, want 0 alongside the error", tt.nums, tt.k, got)
		}
	}
}

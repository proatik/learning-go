package twosum

import (
	"slices"
	"testing"
)

func TestTwoSum(t *testing.T) {
	tests := []struct {
		nums   []int
		target int
		i, j   int
		ok     bool
	}{
		{[]int{2, 7, 11, 15}, 9, 0, 1, true},
		{[]int{3, 2, 4}, 6, 1, 2, true},
		{[]int{3, 3}, 6, 0, 1, true},
		{[]int{1, 5, 1, 5}, 6, 0, 1, true},
		{[]int{-3, 4, 3, 90}, 0, 0, 2, true},
		{[]int{0, 4, 3, 0}, 0, 0, 3, true},
		{[]int{2, 2, 2, 4}, 6, 0, 3, true},
		{[]int{5, 5, 5}, 10, 0, 1, true},
		{[]int{1, 9, 4, 6, 5, 5}, 10, 0, 1, true},
		{[]int{4, 1, 6, 9}, 10, 0, 2, true},
		{[]int{8, 1, 1, 2}, 3, 1, 3, true},
		{[]int{1, 2, 3}, 7, 0, 0, false},
		{[]int{3}, 6, 0, 0, false},
		{[]int{3, 4}, 6, 0, 0, false},
		{[]int{}, 0, 0, 0, false},
		{nil, 5, 0, 0, false},
	}
	for _, tt := range tests {
		nums := slices.Clone(tt.nums)
		i, j, ok := TwoSum(nums, tt.target)
		if i != tt.i || j != tt.j || ok != tt.ok {
			t.Errorf("TwoSum(%v, %d) = (%d, %d, %t), want (%d, %d, %t)",
				tt.nums, tt.target, i, j, ok, tt.i, tt.j, tt.ok)
		}
		if !slices.Equal(nums, tt.nums) {
			t.Errorf("TwoSum modified its input: %v became %v", tt.nums, nums)
		}
	}
}

func TestTwoSumLarge(t *testing.T) {
	const n = 20000
	nums := make([]int, n)
	for k := range nums {
		nums[k] = k * 2 // all even
	}
	target := nums[n-2] + nums[n-1]
	i, j, ok := TwoSum(nums, target)
	if !ok || i != n-2 || j != n-1 {
		t.Errorf("TwoSum(0, 2, ..., %d; %d) = (%d, %d, %t), want (%d, %d, true)",
			nums[n-1], target, i, j, ok, n-2, n-1)
	}
}

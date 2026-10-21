package topk

import (
	"cmp"
	"fmt"
	"slices"
)

// TopK returns the k largest values of nums, largest first.
// If k > len(nums) it returns all values; if k <= 0 an empty slice.
// nums is not modified.
func TopK(nums []int, k int) []int {
	sorted := slices.Clone(nums) // a new backing array: the result shares nothing with nums
	slices.SortFunc(sorted, func(a, b int) int { return cmp.Compare(b, a) })
	k = max(0, min(k, len(sorted)))
	return sorted[:k]
}

// KthSmallest returns the k-th smallest value of nums (k counts from 1).
// It returns an error if k is out of range. nums is not modified.
func KthSmallest(nums []int, k int) (int, error) {
	if k < 1 || k > len(nums) {
		return 0, fmt.Errorf("k = %d out of range [1, %d]", k, len(nums))
	}
	sorted := slices.Clone(nums)
	slices.Sort(sorted)
	return sorted[k-1], nil // k is 1-based, indices are 0-based
}

// KthLargestDistinct returns the k-th largest distinct value of nums
// (k counts from 1). It returns an error if there is no such value.
// nums is not modified.
func KthLargestDistinct(nums []int, k int) (int, error) {
	distinct := slices.Clone(nums)
	slices.Sort(distinct)
	distinct = slices.Compact(distinct) // ascending, no duplicates
	if k < 1 || k > len(distinct) {
		return 0, fmt.Errorf("k = %d out of range [1, %d]", k, len(distinct))
	}
	return distinct[len(distinct)-k], nil // count from the end
}

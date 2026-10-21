package median

import (
	"errors"
	"slices"
)

// Median returns the median of nums without modifying nums.
// For an even count it is the mean of the two middle values.
// It returns an error if nums is empty.
func Median(nums []float64) (float64, error) {
	if len(nums) == 0 {
		return 0, errors.New("median of an empty slice")
	}
	sorted := slices.Clone(nums) // sort a copy: the caller's order is preserved
	slices.Sort(sorted)

	mid := len(sorted) / 2
	if len(sorted)%2 == 1 {
		return sorted[mid], nil
	}
	return (sorted[mid-1] + sorted[mid]) / 2, nil
}

// Spread returns the largest value minus the smallest value in nums.
// It returns an error if nums is empty.
func Spread(nums []int) (int, error) {
	if len(nums) == 0 {
		return 0, errors.New("spread of an empty slice") // Min/Max would panic
	}
	return slices.Max(nums) - slices.Min(nums), nil
}

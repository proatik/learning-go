package stats

import "errors"

// Sum returns the sum of nums. The sum of no numbers is 0.
func Sum(nums ...int) int {
	total := 0
	for _, n := range nums {
		total += n
	}
	return total
}

// Mean returns the arithmetic mean of nums.
// It returns an error if nums is empty.
func Mean(nums ...int) (float64, error) {
	if len(nums) == 0 {
		return 0, errors.New("mean of no numbers")
	}
	// Spread the slice back into Sum's variadic parameter.
	return float64(Sum(nums...)) / float64(len(nums)), nil
}

package sums

// SumArray returns the sum of the five elements of a.
// a is a copy of the caller's array: arrays are values.
func SumArray(a [5]int) int {
	total := 0
	for _, v := range a {
		total += v
	}
	return total
}

// Sum returns the sum of the elements of s. The sum of an empty slice is 0.
// Ranging over a nil slice runs zero iterations, so no special case is needed.
func Sum(s []int) int {
	total := 0
	for _, v := range s {
		total += v
	}
	return total
}

// Largest returns the largest element of s and ok == true.
// For an empty (or nil) slice it returns 0, false.
func Largest(s []int) (largest int, ok bool) {
	if len(s) == 0 {
		return 0, false
	}
	largest = s[0] // start from a real element, not 0
	for _, v := range s[1:] {
		if v > largest {
			largest = v
		}
	}
	return largest, true
}

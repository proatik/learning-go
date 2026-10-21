package average

// IntAverage returns the average of a, b and c using integer division
// (the result is truncated toward zero).
func IntAverage(a, b, c int) int {
	return (a + b + c) / 3
}

// Average returns the exact average of a, b and c.
func Average(a, b, c int) float64 {
	sum := a + b + c
	// Convert before dividing; the untyped constant 3 becomes float64.
	return float64(sum) / 3
}

// Percent returns part as a percentage of total. total must not be zero.
func Percent(part, total int) float64 {
	return float64(part) / float64(total) * 100
}

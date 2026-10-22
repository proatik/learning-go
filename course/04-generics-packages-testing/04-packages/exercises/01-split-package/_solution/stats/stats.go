// Package stats computes simple statistics over integers.
package stats

// Sum returns the sum of xs.
func Sum(xs []int) int {
	total := 0
	for _, x := range xs {
		total += x
	}
	return total
}

// Mean returns the average of xs, or 0 when xs is empty.
func Mean(xs []int) float64 {
	if len(xs) == 0 {
		return 0
	}
	return float64(Sum(xs)) / float64(len(xs))
}

// MinMax returns the smallest and largest value of xs and true,
// or 0, 0, false when xs is empty.
func MinMax(xs []int) (lo, hi int, ok bool) {
	if len(xs) == 0 {
		return 0, 0, false
	}
	lo, hi = xs[0], xs[0]
	for _, x := range xs[1:] {
		lo = min(lo, x)
		hi = max(hi, x)
	}
	return lo, hi, true
}

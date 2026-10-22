package stats

import (
	"math"
	"slices"
)

// Number is every integer and float type, including named types based on them.
type Number interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 |
		~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 |
		~float32 | ~float64
}

// Summary describes a sample of numbers. All statistics are float64.
type Summary struct {
	Count                          int
	Min, Max, Mean, Median, StdDev float64
}

// Mean returns the arithmetic mean of xs. ok is false when xs is empty.
func Mean[T Number](xs []T) (float64, bool) {
	if len(xs) == 0 {
		return 0, false
	}
	var sum float64
	for _, x := range xs {
		sum += float64(x)
	}
	return sum / float64(len(xs)), true
}

// Variance returns the population variance of xs. ok is false when xs is empty.
func Variance[T Number](xs []T) (float64, bool) {
	mean, ok := Mean(xs)
	if !ok {
		return 0, false
	}
	var sum float64
	for _, x := range xs {
		d := float64(x) - mean
		sum += d * d
	}
	return sum / float64(len(xs)), true
}

// Median returns the median of xs without modifying it. ok is false when xs is empty.
func Median[T Number](xs []T) (float64, bool) {
	n := len(xs)
	if n == 0 {
		return 0, false
	}
	sorted := slices.Clone(xs)
	slices.Sort(sorted)
	if n%2 == 1 {
		return float64(sorted[n/2]), true
	}
	return (float64(sorted[n/2-1]) + float64(sorted[n/2])) / 2, true
}

// Summarize returns the count, min, max, mean, median and standard deviation of xs.
// ok is false when xs is empty.
func Summarize[T Number](xs []T) (Summary, bool) {
	if len(xs) == 0 {
		return Summary{}, false
	}
	mean, _ := Mean(xs)
	median, _ := Median(xs)
	variance, _ := Variance(xs)
	return Summary{
		Count:  len(xs),
		Min:    float64(slices.Min(xs)),
		Max:    float64(slices.Max(xs)),
		Mean:   mean,
		Median: median,
		StdDev: math.Sqrt(variance),
	}, true
}

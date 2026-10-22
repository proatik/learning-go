package stats

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
	// TODO: add float64(x), never x itself
	return 0, false
}

// Variance returns the population variance of xs. ok is false when xs is empty.
func Variance[T Number](xs []T) (float64, bool) {
	// TODO
	return 0, false
}

// Median returns the median of xs without modifying it. ok is false when xs is empty.
func Median[T Number](xs []T) (float64, bool) {
	// TODO: sort a copy
	return 0, false
}

// Summarize returns the count, min, max, mean, median and standard deviation of xs.
// ok is false when xs is empty.
func Summarize[T Number](xs []T) (Summary, bool) {
	// TODO
	return Summary{}, false
}

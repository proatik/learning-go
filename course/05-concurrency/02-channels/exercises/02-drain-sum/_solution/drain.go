package drain

// Sum receives every value from ch until it is closed and returns their
// total.
func Sum(ch <-chan int) int {
	total := 0
	for v := range ch {
		total += v
	}
	return total
}

// Drain receives every value from ch until it is closed and returns them in
// order.
func Drain[T any](ch <-chan T) []T {
	var out []T
	for v := range ch {
		out = append(out, v)
	}
	return out
}

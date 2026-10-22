package drain

// Sum receives every value from ch until it is closed and returns their
// total.
func Sum(ch <-chan int) int {
	// TODO
	return 0
}

// Drain receives every value from ch until it is closed and returns them in
// order.
func Drain[T any](ch <-chan T) []T {
	// TODO
	return nil
}

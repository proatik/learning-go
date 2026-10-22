package pipeline

// Gen returns a channel that yields nums in order and is then closed.
func Gen(nums ...int) <-chan int {
	// TODO
	return nil
}

// Square returns a channel that yields n*n for every n received from in. It
// is closed after in is closed and drained.
func Square(in <-chan int) <-chan int {
	// TODO
	return nil
}

// Sum receives every value from in until it is closed and returns the total.
func Sum(in <-chan int) int {
	// TODO
	return 0
}

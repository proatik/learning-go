package pipeline

// Gen returns a channel that yields nums in order and is then closed.
func Gen(nums ...int) <-chan int {
	out := make(chan int)
	go func() {
		defer close(out)
		for _, n := range nums {
			out <- n
		}
	}()
	return out
}

// Square returns a channel that yields n*n for every n received from in. It
// is closed after in is closed and drained.
func Square(in <-chan int) <-chan int {
	out := make(chan int)
	go func() {
		defer close(out)
		for n := range in {
			out <- n * n
		}
	}()
	return out
}

// Sum receives every value from in until it is closed and returns the total.
func Sum(in <-chan int) int {
	total := 0
	for n := range in {
		total += n
	}
	return total
}

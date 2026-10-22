package stream

// From returns a channel that receives start, start+1, start+2, … until done
// is closed. Then the channel is closed and the goroutine exits.
func From(done <-chan struct{}, start int) <-chan int {
	// TODO
	return nil
}

// Take returns a channel that receives the first n values from in and is then
// closed. It closes early if in is closed or done is closed.
func Take(done <-chan struct{}, in <-chan int, n int) <-chan int {
	// TODO
	return nil
}

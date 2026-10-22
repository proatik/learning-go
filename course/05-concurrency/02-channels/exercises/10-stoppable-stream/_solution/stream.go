package stream

// From returns a channel that receives start, start+1, start+2, … until done
// is closed. Then the channel is closed and the goroutine exits.
func From(done <-chan struct{}, start int) <-chan int {
	out := make(chan int)
	go func() {
		defer close(out)
		for i := start; ; i++ {
			select {
			case out <- i:
			case <-done:
				return
			}
		}
	}()
	return out
}

// Take returns a channel that receives the first n values from in and is then
// closed. It closes early if in is closed or done is closed.
func Take(done <-chan struct{}, in <-chan int, n int) <-chan int {
	out := make(chan int)
	go func() {
		defer close(out)
		for range n {
			var v int
			select {
			case x, ok := <-in:
				if !ok {
					return
				}
				v = x
			case <-done:
				return
			}
			select {
			case out <- v:
			case <-done:
				return
			}
		}
	}()
	return out
}

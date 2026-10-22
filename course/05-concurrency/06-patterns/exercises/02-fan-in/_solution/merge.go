package fanin

import "sync"

// Merge returns one channel that yields every value received from every
// input. It is closed once all inputs are closed and drained. Nil inputs are
// ignored.
func Merge[T any](inputs ...<-chan T) <-chan T {
	out := make(chan T)
	var wg sync.WaitGroup
	for _, in := range inputs {
		if in == nil {
			continue
		}
		wg.Go(func() {
			for v := range in {
				out <- v
			}
		})
	}
	go func() {
		wg.Wait()
		close(out)
	}()
	return out
}

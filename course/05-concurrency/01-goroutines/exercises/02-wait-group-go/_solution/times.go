package times

import "sync"

// Times calls f(0) … f(n-1), each in its own goroutine, waits for all of
// them and returns how many goroutines it started.
func Times(n int, f func(i int)) int {
	var wg sync.WaitGroup
	started := 0
	for i := range n {
		wg.Go(func() { f(i) })
		started++
	}
	wg.Wait()
	return started
}

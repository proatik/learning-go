package runall

import "sync"

// RunAll runs every task in its own goroutine and waits for all of them to
// finish.
func RunAll(tasks []func()) {
	var wg sync.WaitGroup
	for _, task := range tasks {
		wg.Add(1)
		go func() {
			defer wg.Done()
			task()
		}()
	}
	wg.Wait()
}

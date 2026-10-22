package timed

import (
	"sync"
	"time"
)

// Timed runs every task concurrently and returns the wall time of the whole
// batch and the time each task took, in input order.
func Timed(tasks []func()) (total time.Duration, each []time.Duration) {
	each = make([]time.Duration, len(tasks))
	start := time.Now()
	var wg sync.WaitGroup
	for i, task := range tasks {
		wg.Go(func() {
			t0 := time.Now()
			task()
			each[i] = time.Since(t0)
		})
	}
	wg.Wait()
	return time.Since(start), each
}

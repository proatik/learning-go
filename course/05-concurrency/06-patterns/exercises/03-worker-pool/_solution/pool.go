package pool

import "sync"

// Result pairs a job with the value work produced for it.
type Result struct {
	Job   int
	Value int
}

// Run starts n workers (at least 1) that take jobs from jobs until it is
// closed and send Result{job, work(job)} on the returned channel. The
// returned channel is closed after every worker has finished.
func Run(n int, jobs <-chan int, work func(int) int) <-chan Result {
	n = max(n, 1)
	out := make(chan Result)
	var wg sync.WaitGroup
	for range n {
		wg.Go(func() {
			for j := range jobs {
				out <- Result{Job: j, Value: work(j)}
			}
		})
	}
	go func() {
		wg.Wait()
		close(out)
	}()
	return out
}

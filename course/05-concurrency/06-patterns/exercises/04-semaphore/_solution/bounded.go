package bounded

import "sync"

// Map calls f on every element of in, running at most limit calls at once
// (limit < 1 means 1), and returns the results in input order.
func Map[T, R any](limit int, in []T, f func(T) R) []R {
	out := make([]R, len(in))
	sem := make(chan struct{}, max(limit, 1))
	var wg sync.WaitGroup
	for i, v := range in {
		sem <- struct{}{} // acquire before starting: at most limit goroutines exist
		wg.Go(func() {
			defer func() { <-sem }()
			out[i] = f(v)
		})
	}
	wg.Wait()
	return out
}

package pmap

import "sync"

// Map returns f applied to every element of in, computed concurrently, in
// input order.
func Map[T, U any](in []T, f func(T) U) []U {
	out := make([]U, len(in))
	var wg sync.WaitGroup
	for i, v := range in {
		wg.Go(func() { out[i] = f(v) })
	}
	wg.Wait()
	return out
}

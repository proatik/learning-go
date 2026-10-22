package gen

import "context"

// Count returns a channel that yields start, start+step, start+2*step, …
// until ctx is done, then closes.
func Count(ctx context.Context, start, step int) <-chan int {
	out := make(chan int)
	go func() {
		defer close(out)
		for n := start; ; n += step {
			select {
			case out <- n:
			case <-ctx.Done():
				return
			}
		}
	}()
	return out
}

// Map returns a channel that yields f(v) for every v received from in. It
// closes when in is closed or ctx is done.
func Map(ctx context.Context, in <-chan int, f func(int) int) <-chan int {
	out := make(chan int)
	go func() {
		defer close(out)
		for {
			var v int
			select {
			case x, ok := <-in:
				if !ok {
					return
				}
				v = x
			case <-ctx.Done():
				return
			}
			select {
			case out <- f(v):
			case <-ctx.Done():
				return
			}
		}
	}()
	return out
}

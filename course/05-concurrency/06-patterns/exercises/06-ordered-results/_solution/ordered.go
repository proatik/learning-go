package ordered

import "context"

// Process calls f on every value received from in, running up to workers
// calls at once (workers < 1 means 1), and sends the results on the returned
// channel in the same order as their inputs. The channel is closed after the
// last result, or as soon as ctx is done.
func Process[T, R any](ctx context.Context, workers int, in <-chan T, f func(T) R) <-chan R {
	workers = max(workers, 1)
	out := make(chan R)
	// Each input gets a one-slot "promise" channel; queue holds the promises
	// in input order, so the emitter can wait for them one by one.
	queue := make(chan chan R, workers)
	sem := make(chan struct{}, workers)

	go func() {
		defer close(queue)
		for {
			var v T
			select {
			case <-ctx.Done():
				return
			case x, ok := <-in:
				if !ok {
					return
				}
				v = x
			}
			select {
			case <-ctx.Done():
				return
			case sem <- struct{}{}:
			}
			p := make(chan R, 1) // buffered: the worker never blocks, even after cancel
			select {
			case <-ctx.Done():
				<-sem
				return
			case queue <- p:
			}
			go func() {
				defer func() { <-sem }()
				p <- f(v)
			}()
		}
	}()

	go func() {
		defer close(out)
		for p := range queue {
			var r R
			select {
			case <-ctx.Done():
				return
			case r = <-p:
			}
			select {
			case <-ctx.Done():
				return
			case out <- r:
			}
		}
	}()
	return out
}

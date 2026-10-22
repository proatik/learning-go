package ordered

import "context"

// Process calls f on every value received from in, running up to workers
// calls at once (workers < 1 means 1), and sends the results on the returned
// channel in the same order as their inputs. The channel is closed after the
// last result, or as soon as ctx is done.
func Process[T, R any](ctx context.Context, workers int, in <-chan T, f func(T) R) <-chan R {
	// TODO
	return nil
}

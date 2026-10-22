package tee

import "context"

// Tee copies every value from in to both returned channels, in order. Both
// outputs are closed when in is closed or ctx is done.
func Tee[T any](ctx context.Context, in <-chan T) (<-chan T, <-chan T) {
	// TODO
	return nil, nil
}

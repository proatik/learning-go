package drain

import "context"

// Drain calls handle for every job received from in, until in is closed or
// ctx is done. It returns the number of jobs handled and nil, or ctx.Err() if
// the context ended the loop.
func Drain(ctx context.Context, in <-chan string, handle func(string)) (int, error) {
	// TODO
	return 0, nil
}

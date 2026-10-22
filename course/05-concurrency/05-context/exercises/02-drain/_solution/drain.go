package drain

import "context"

// Drain calls handle for every job received from in, until in is closed or
// ctx is done. It returns the number of jobs handled and nil, or ctx.Err() if
// the context ended the loop.
func Drain(ctx context.Context, in <-chan string, handle func(string)) (int, error) {
	n := 0
	for {
		if err := ctx.Err(); err != nil {
			return n, err
		}
		select {
		case <-ctx.Done():
			return n, ctx.Err()
		case job, ok := <-in:
			if !ok {
				return n, nil
			}
			handle(job)
			n++
		}
	}
}

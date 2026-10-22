package retry

import (
	"context"
	"errors"
	"fmt"
	"time"
)

// ErrNoTime reports that the context's deadline would pass before the next
// attempt could start.
var ErrNoTime = errors.New("retry: not enough time for another attempt")

// Retry calls op up to attempts times, waiting base, 2*base, 4*base, …
// between failures. It stops early when ctx is done, or when the next wait
// would end at or after ctx's deadline.
func Retry(ctx context.Context, attempts int, base time.Duration, op func(context.Context) error) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	if attempts < 1 {
		attempts = 1
	}
	for n := 1; ; n++ {
		err := op(ctx)
		if err == nil {
			return nil
		}
		if n == attempts {
			return fmt.Errorf("retry: giving up after attempt %d: %w", n, err)
		}
		if ctx.Err() != nil {
			return fmt.Errorf("retry: stopped after attempt %d: %w: %w", n, context.Cause(ctx), err)
		}
		wait := base << (n - 1)
		if dl, ok := ctx.Deadline(); ok && !time.Now().Add(wait).Before(dl) {
			return fmt.Errorf("%w (after attempt %d): %w", ErrNoTime, n, err)
		}
		if !sleep(ctx, wait) {
			return fmt.Errorf("retry: stopped after attempt %d: %w: %w", n, context.Cause(ctx), err)
		}
	}
}

// sleep waits for d and reports true, or returns false as soon as ctx is
// done.
func sleep(ctx context.Context, d time.Duration) bool {
	t := time.NewTimer(d)
	defer t.Stop()
	select {
	case <-t.C:
		return true
	case <-ctx.Done():
		return false
	}
}

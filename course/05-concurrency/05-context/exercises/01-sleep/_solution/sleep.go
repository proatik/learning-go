package nap

import (
	"context"
	"time"
)

// Sleep pauses for d, or until ctx is done. It returns nil if the full
// duration passed, or ctx.Err() if the context ended first.
func Sleep(ctx context.Context, d time.Duration) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	if d <= 0 {
		return nil
	}
	t := time.NewTimer(d)
	defer t.Stop()
	select {
	case <-t.C:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

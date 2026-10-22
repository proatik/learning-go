package nap

import (
	"context"
	"time"
)

// Sleep pauses for d, or until ctx is done. It returns nil if the full
// duration passed, or ctx.Err() if the context ended first.
func Sleep(ctx context.Context, d time.Duration) error {
	// TODO
	return nil
}

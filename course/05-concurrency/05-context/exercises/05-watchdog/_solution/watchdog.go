package watchdog

import (
	"context"
	"errors"
	"time"
)

var (
	// ErrStalled is the cause when no heartbeat arrives within the timeout.
	ErrStalled = errors.New("watchdog: no heartbeat in time")
	// ErrStopped is the cause when the heartbeat channel is closed.
	ErrStopped = errors.New("watchdog: heartbeats stopped")
)

// Watch returns a child of parent that is canceled with cause ErrStalled if
// no value arrives on beats for timeout, or ErrStopped if beats is closed.
// The returned CancelFunc cancels it with cause context.Canceled.
func Watch(parent context.Context, beats <-chan struct{}, timeout time.Duration) (context.Context, context.CancelFunc) {
	ctx, cancel := context.WithCancelCause(parent)
	go func() {
		t := time.NewTimer(timeout)
		defer t.Stop()
		for {
			select {
			case <-ctx.Done():
				return
			case _, ok := <-beats:
				if !ok {
					cancel(ErrStopped)
					return
				}
				t.Reset(timeout)
			case <-t.C:
				cancel(ErrStalled)
				return
			}
		}
	}()
	return ctx, func() { cancel(nil) }
}

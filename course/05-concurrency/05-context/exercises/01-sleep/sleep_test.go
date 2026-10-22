package nap

import (
	"context"
	"errors"
	"testing"
	"testing/synctest"
	"time"
)

func TestSleep(t *testing.T) {
	tests := []struct {
		name      string
		timeout   time.Duration // > 0: ctx from WithTimeout
		cancelAt  time.Duration // > 0: cancel() after this long; < 0: cancel before calling
		d         time.Duration
		wantErr   error
		wantAfter time.Duration
	}{
		{name: "Background, 1s", d: time.Second, wantAfter: time.Second},
		{name: "Background, 90m", d: 90 * time.Minute, wantAfter: 90 * time.Minute},
		{name: "timeout 2s, 1h", timeout: 2 * time.Second, d: time.Hour, wantErr: context.DeadlineExceeded, wantAfter: 2 * time.Second},
		{name: "timeout 2s, 1s", timeout: 2 * time.Second, d: time.Second, wantAfter: time.Second},
		{name: "cancel after 300ms, 1h", cancelAt: 300 * time.Millisecond, d: time.Hour, wantErr: context.Canceled, wantAfter: 300 * time.Millisecond},
		{name: "already canceled, 1s", cancelAt: -1, d: time.Second, wantErr: context.Canceled},
		{name: "Background, 0", d: 0},
		{name: "Background, -1s", d: -time.Second},
		{name: "already canceled, 0", cancelAt: -1, d: 0, wantErr: context.Canceled},
		{name: "already canceled, -1s", cancelAt: -1, d: -time.Second, wantErr: context.Canceled},
	}
	for _, tt := range tests {
		synctest.Test(t, func(t *testing.T) {
			ctx, cancel := context.WithCancel(t.Context())
			defer cancel()
			if tt.timeout > 0 {
				ctx, cancel = context.WithTimeout(ctx, tt.timeout)
				defer cancel()
			}
			switch {
			case tt.cancelAt < 0:
				cancel()
			case tt.cancelAt > 0:
				time.AfterFunc(tt.cancelAt, cancel)
			}
			start := time.Now()
			err := Sleep(ctx, tt.d)
			elapsed := time.Since(start)
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("Sleep(%s) error = %v, want %v", tt.name, err, tt.wantErr)
			}
			if elapsed != tt.wantAfter {
				t.Errorf("Sleep(%s) returned after %v, want %v", tt.name, elapsed, tt.wantAfter)
			}
		})
	}
}

func TestSleepReturnsCtxErr(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		ErrShutdown := errors.New("shutdown")
		ctx, cancel := context.WithCancelCause(t.Context())
		cancel(ErrShutdown)
		if err := Sleep(ctx, time.Second); err != context.Canceled {
			t.Errorf("Sleep(canceled with a cause, 1s) error = %v, want %v (ctx.Err(), not the cause)", err, context.Canceled)
		}
	})
}

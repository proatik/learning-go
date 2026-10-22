package retry

import (
	"context"
	"errors"
	"slices"
	"testing"
	"testing/synctest"
	"time"
)

var errBoom = errors.New("boom")

// flaky fails its first `fails` calls (-1: always), taking `took` per call
// unless its ctx ends first. It records when each call started.
type flaky struct {
	fails int
	took  time.Duration
	start time.Time
	at    []time.Duration
}

func (f *flaky) op(ctx context.Context) error {
	f.at = append(f.at, time.Since(f.start))
	if f.took > 0 {
		select {
		case <-time.After(f.took):
		case <-ctx.Done():
			return ctx.Err()
		}
	}
	if f.fails < 0 || len(f.at) <= f.fails {
		return errBoom
	}
	return nil
}

type tcase struct {
	name      string
	attempts  int
	base      time.Duration
	fails     int
	took      time.Duration
	deadline  time.Duration // > 0: ctx deadline
	cancelAt  time.Duration // > 0: cancel then; < 0: canceled before the call
	cause     error         // cancel with this cause
	wantAt    []time.Duration
	wantMsg   string  // "" means nil
	wantIs    []error // errors.Is targets
	wantAfter time.Duration
}

func run(t *testing.T, tt tcase) {
	synctest.Test(t, func(t *testing.T) {
		ctx, cancel := context.WithCancelCause(t.Context())
		defer cancel(nil)
		if tt.deadline > 0 {
			var c context.CancelFunc
			ctx, c = context.WithTimeout(ctx, tt.deadline)
			defer c()
		}
		switch {
		case tt.cancelAt < 0:
			cancel(tt.cause)
		case tt.cancelAt > 0:
			time.AfterFunc(tt.cancelAt, func() { cancel(tt.cause) })
		}
		f := &flaky{fails: tt.fails, took: tt.took, start: time.Now()}
		err := Retry(ctx, tt.attempts, tt.base, f.op)
		elapsed := time.Since(f.start)

		if !slices.Equal(f.at, tt.wantAt) {
			t.Errorf("Retry(%s) attempts started at %v, want %v", tt.name, f.at, tt.wantAt)
		}
		switch {
		case tt.wantMsg == "" && err != nil:
			t.Errorf("Retry(%s) error = %v, want <nil>", tt.name, err)
		case tt.wantMsg != "" && (err == nil || err.Error() != tt.wantMsg):
			t.Errorf("Retry(%s) error = %v, want %q", tt.name, err, tt.wantMsg)
		}
		for _, target := range tt.wantIs {
			if !errors.Is(err, target) {
				t.Errorf("Retry(%s) error = %v, want errors.Is(err, %v)", tt.name, err, target)
			}
		}
		if elapsed != tt.wantAfter {
			t.Errorf("Retry(%s) returned after %v, want %v", tt.name, elapsed, tt.wantAfter)
		}
	})
}

const ms = time.Millisecond

func TestRetry(t *testing.T) {
	tests := []tcase{
		{
			name: "succeeds at once", attempts: 3, base: time.Second, fails: 0,
			wantAt: []time.Duration{0},
		},
		{
			name: "fails twice then succeeds, base 1s", attempts: 5, base: time.Second, fails: 2,
			wantAt: []time.Duration{0, time.Second, 3 * time.Second}, wantAfter: 3 * time.Second,
		},
		{
			name: "always fails, 3 attempts, base 100ms", attempts: 3, base: 100 * ms, fails: -1,
			wantAt:  []time.Duration{0, 100 * ms, 300 * ms},
			wantMsg: "retry: giving up after attempt 3: boom", wantIs: []error{errBoom}, wantAfter: 300 * ms,
		},
		{
			name: "always fails, 5 attempts, base 10ms", attempts: 5, base: 10 * ms, fails: -1,
			wantAt:  []time.Duration{0, 10 * ms, 30 * ms, 70 * ms, 150 * ms},
			wantMsg: "retry: giving up after attempt 5: boom", wantIs: []error{errBoom}, wantAfter: 150 * ms,
		},
		{
			name: "attempts 0 counts as 1", attempts: 0, base: time.Second, fails: -1,
			wantAt:  []time.Duration{0},
			wantMsg: "retry: giving up after attempt 1: boom", wantIs: []error{errBoom},
		},
		{
			name: "attempts -4 counts as 1", attempts: -4, base: time.Second, fails: 0,
			wantAt: []time.Duration{0},
		},
		{
			name: "each call takes 1s, fails once, base 500ms", attempts: 3, base: 500 * ms, fails: 1, took: time.Second,
			wantAt: []time.Duration{0, 1500 * ms}, wantAfter: 2500 * ms,
		},
	}
	for _, tt := range tests {
		run(t, tt)
	}
}

func TestRetryCanceled(t *testing.T) {
	ErrShutdown := errors.New("shutdown")
	tests := []tcase{
		{
			name: "cancel at 1.5s, base 1s", attempts: 10, base: time.Second, fails: -1, cancelAt: 1500 * ms,
			wantAt:  []time.Duration{0, time.Second},
			wantMsg: "retry: stopped after attempt 2: context canceled: boom", wantIs: []error{errBoom, context.Canceled}, wantAfter: 1500 * ms,
		},
		{
			name: "cancel(ErrShutdown) at 500ms, base 1s", attempts: 10, base: time.Second, fails: -1, cancelAt: 500 * ms, cause: ErrShutdown,
			wantAt:  []time.Duration{0},
			wantMsg: "retry: stopped after attempt 1: shutdown: boom", wantIs: []error{errBoom, ErrShutdown}, wantAfter: 500 * ms,
		},
		{
			name: "already canceled", attempts: 3, base: time.Second, fails: -1, cancelAt: -1,
			wantMsg: "context canceled", wantIs: []error{context.Canceled},
		},
		{
			name: "canceled during the call (takes 2s, cancel at 1s)", attempts: 3, base: time.Second, fails: -1, took: 2 * time.Second, cancelAt: time.Second,
			wantAt:  []time.Duration{0},
			wantMsg: "retry: stopped after attempt 1: context canceled: context canceled", wantIs: []error{context.Canceled}, wantAfter: time.Second,
		},
	}
	for _, tt := range tests {
		run(t, tt)
	}
}

func TestRetryDeadline(t *testing.T) {
	tests := []tcase{
		{
			name: "deadline 5s, base 1s", attempts: 10, base: time.Second, fails: -1, deadline: 5 * time.Second,
			wantAt:  []time.Duration{0, time.Second, 3 * time.Second},
			wantMsg: "retry: not enough time for another attempt (after attempt 3): boom", wantIs: []error{ErrNoTime, errBoom}, wantAfter: 3 * time.Second,
		},
		{
			name: "deadline exactly 7s, base 1s", attempts: 10, base: time.Second, fails: -1, deadline: 7 * time.Second,
			wantAt:  []time.Duration{0, time.Second, 3 * time.Second},
			wantMsg: "retry: not enough time for another attempt (after attempt 3): boom", wantIs: []error{ErrNoTime, errBoom}, wantAfter: 3 * time.Second,
		},
		{
			name: "deadline 7s+1ns, base 1s", attempts: 10, base: time.Second, fails: -1, deadline: 7*time.Second + 1,
			wantAt:  []time.Duration{0, time.Second, 3 * time.Second, 7 * time.Second},
			wantMsg: "retry: not enough time for another attempt (after attempt 4): boom", wantIs: []error{ErrNoTime, errBoom}, wantAfter: 7 * time.Second,
		},
		{
			name: "deadline 500ms, base 1s", attempts: 10, base: time.Second, fails: -1, deadline: 500 * ms,
			wantAt:  []time.Duration{0},
			wantMsg: "retry: not enough time for another attempt (after attempt 1): boom", wantIs: []error{ErrNoTime, errBoom},
		},
		{
			name: "deadline 10s, base 1s, succeeds on attempt 3", attempts: 10, base: time.Second, fails: 2, deadline: 10 * time.Second,
			wantAt: []time.Duration{0, time.Second, 3 * time.Second}, wantAfter: 3 * time.Second,
		},
		{
			name: "deadline 3s, calls take 2s, base 100ms", attempts: 10, base: 100 * ms, fails: -1, took: 2 * time.Second, deadline: 3 * time.Second,
			wantAt:  []time.Duration{0, 2100 * ms},
			wantMsg: "retry: stopped after attempt 2: context deadline exceeded: context deadline exceeded", wantIs: []error{context.DeadlineExceeded}, wantAfter: 3 * time.Second,
		},
		{
			name: "deadline 3s, last attempt allowed", attempts: 2, base: time.Second, fails: -1, deadline: 3 * time.Second,
			wantAt:  []time.Duration{0, time.Second},
			wantMsg: "retry: giving up after attempt 2: boom", wantIs: []error{errBoom}, wantAfter: time.Second,
		},
	}
	for _, tt := range tests {
		run(t, tt)
	}
}

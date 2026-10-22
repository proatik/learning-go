package retry

import (
	"errors"
	"fmt"
	"slices"
	"testing"
	"testing/synctest"
	"time"
)

const ms = time.Millisecond

func TestDo(t *testing.T) {
	tests := []struct {
		attempts       int
		base, maxDelay time.Duration
		succeedOn      int           // the call that returns nil; 0 = never
		stopAt         time.Duration // when stop is closed; 0 = never (stop is nil)
		wantCalls      []time.Duration
		wantMsg        string // "" means nil
		wantSentinel   error
		wantAt         time.Duration
	}{
		{attempts: 4, base: 100 * ms, maxDelay: time.Second,
			wantCalls: []time.Duration{0, 100 * ms, 300 * ms, 700 * ms},
			wantMsg:   "retry: 4 attempts failed: fail 4", wantAt: 700 * ms},
		{attempts: 5, base: 100 * ms, maxDelay: time.Second, succeedOn: 3,
			wantCalls: []time.Duration{0, 100 * ms, 300 * ms}, wantAt: 300 * ms},
		{attempts: 5, base: 100 * ms, maxDelay: time.Second, succeedOn: 1,
			wantCalls: []time.Duration{0}, wantAt: 0},
		{attempts: 6, base: 100 * ms, maxDelay: 300 * ms,
			wantCalls: []time.Duration{0, 100 * ms, 300 * ms, 600 * ms, 900 * ms, 1200 * ms},
			wantMsg:   "retry: 6 attempts failed: fail 6", wantAt: 1200 * ms},
		{attempts: 3, base: time.Second, maxDelay: time.Hour, succeedOn: 3,
			wantCalls: []time.Duration{0, time.Second, 3 * time.Second}, wantAt: 3 * time.Second},
		{attempts: 1, base: 100 * ms, maxDelay: time.Second,
			wantCalls: []time.Duration{0},
			wantMsg:   "retry: 1 attempts failed: fail 1", wantAt: 0},
		{attempts: 0, base: 100 * ms, maxDelay: time.Second,
			wantMsg: "retry: no attempts", wantSentinel: ErrNoAttempts},
		{attempts: -3, base: 100 * ms, maxDelay: time.Second,
			wantMsg: "retry: no attempts", wantSentinel: ErrNoAttempts},
		{attempts: 5, base: 100 * ms, maxDelay: time.Second, stopAt: 250 * ms,
			wantCalls: []time.Duration{0, 100 * ms},
			wantMsg:   "retry: stopped after 2 attempts: fail 2", wantSentinel: ErrStopped, wantAt: 250 * ms},
		{attempts: 3, base: 100 * ms, maxDelay: time.Second, stopAt: 50 * ms,
			wantCalls: []time.Duration{0},
			wantMsg:   "retry: stopped after 1 attempts: fail 1", wantSentinel: ErrStopped, wantAt: 50 * ms},
		{attempts: 3, base: 100 * ms, maxDelay: time.Second, stopAt: time.Hour,
			wantCalls: []time.Duration{0, 100 * ms, 300 * ms},
			wantMsg:   "retry: 3 attempts failed: fail 3", wantAt: 300 * ms},
	}
	for _, tt := range tests {
		synctest.Test(t, func(t *testing.T) {
			var stop chan struct{}
			stopDesc := "nil"
			if tt.stopAt > 0 {
				stop = make(chan struct{})
				time.AfterFunc(tt.stopAt, func() { close(stop) })
				stopDesc = fmt.Sprintf("closed at %v", tt.stopAt)
			}
			succeed := "never succeeds"
			if tt.succeedOn > 0 {
				succeed = fmt.Sprintf("succeeds on call %d", tt.succeedOn)
			}
			call := fmt.Sprintf("Do(stop %s, %d, %v, %v, f %s)", stopDesc, tt.attempts, tt.base, tt.maxDelay, succeed)

			start := time.Now()
			var calls []time.Duration
			var lastErr error
			err := Do(stop, tt.attempts, tt.base, tt.maxDelay, func() error {
				calls = append(calls, time.Since(start))
				if len(calls) == tt.succeedOn {
					return nil
				}
				lastErr = errors.New(fmt.Sprintf("fail %d", len(calls)))
				return lastErr
			})
			at := time.Since(start)

			switch {
			case tt.wantMsg == "" && err != nil:
				t.Errorf("%s error = %v, want <nil>", call, err)
			case tt.wantMsg != "" && err == nil:
				t.Errorf("%s error = <nil>, want %q", call, tt.wantMsg)
			case tt.wantMsg != "" && err.Error() != tt.wantMsg:
				t.Errorf("%s error = %q, want %q", call, err.Error(), tt.wantMsg)
			case tt.wantSentinel != nil && !errors.Is(err, tt.wantSentinel):
				t.Errorf("%s error = %v, want one that errors.Is %v", call, err, tt.wantSentinel)
			case tt.wantMsg != "" && tt.wantSentinel != ErrNoAttempts && !errors.Is(err, lastErr):
				t.Errorf("%s error = %v, want one that wraps f's last error (use %%w)", call, err)
			}
			if !slices.Equal(calls, tt.wantCalls) {
				t.Errorf("%s called f at %v, want %v", call, calls, tt.wantCalls)
			}
			if at != tt.wantAt {
				t.Errorf("%s returned after %v, want %v", call, at, tt.wantAt)
			}
		})
	}
}

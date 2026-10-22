package heartbeat

import (
	"fmt"
	"slices"
	"testing"
	"testing/synctest"
	"time"
)

const ms = time.Millisecond

func TestEvery(t *testing.T) {
	tests := []struct {
		interval time.Duration
		stopAt   time.Duration // -1: stop is closed before the call
		wantAt   []time.Duration
	}{
		{100 * ms, 1050 * ms, []time.Duration{100 * ms, 200 * ms, 300 * ms, 400 * ms, 500 * ms, 600 * ms, 700 * ms, 800 * ms, 900 * ms, 1000 * ms}},
		{time.Second, 2500 * ms, []time.Duration{time.Second, 2 * time.Second}},
		{time.Second, 999 * ms, nil},
		{250 * ms, 260 * ms, []time.Duration{250 * ms}},
		{time.Second, -1, nil},
	}
	for _, tt := range tests {
		synctest.Test(t, func(t *testing.T) {
			stop := make(chan struct{})
			returnAt := tt.stopAt
			if tt.stopAt < 0 {
				close(stop)
				returnAt = 0
			} else {
				time.AfterFunc(tt.stopAt, func() { close(stop) })
			}
			call := fmt.Sprintf("Every(%v, stop closed at %v, f)", tt.interval, returnAt)

			start := time.Now()
			var calls []int
			var at []time.Duration
			n := Every(tt.interval, stop, func(n int) {
				calls = append(calls, n)
				at = append(at, time.Since(start))
			})
			elapsed := time.Since(start)

			var wantCalls []int
			for i := range tt.wantAt {
				wantCalls = append(wantCalls, i+1)
			}
			if n != len(tt.wantAt) {
				t.Errorf("%s = %d, want %d", call, n, len(tt.wantAt))
			}
			if !slices.Equal(calls, wantCalls) {
				t.Errorf("%s called f with %v, want %v", call, calls, wantCalls)
			}
			if !slices.Equal(at, tt.wantAt) {
				t.Errorf("%s called f at %v, want %v", call, at, tt.wantAt)
			}
			if elapsed != returnAt {
				t.Errorf("%s returned after %v, want %v", call, elapsed, returnAt)
			}
		})
	}
}

func TestEveryNonPositive(t *testing.T) {
	for _, interval := range []time.Duration{0, -time.Second} {
		synctest.Test(t, func(t *testing.T) {
			stop := make(chan struct{})
			time.AfterFunc(time.Hour, func() { close(stop) })
			calls := 0
			var n int
			start := time.Now()
			func() {
				defer func() {
					if r := recover(); r != nil {
						t.Errorf("Every(%v, stop, f) panicked: %v, want 0 returned at once", interval, r)
					}
				}()
				n = Every(interval, stop, func(int) { calls++ })
			}()
			if n != 0 || calls != 0 || time.Since(start) != 0 {
				t.Errorf("Every(%v, stop, f) = %d after %v with %d calls, want 0 after 0s with 0 calls", interval, n, time.Since(start), calls)
			}
		})
	}
}

func TestEverySlowF(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		stop := make(chan struct{})
		time.AfterFunc(550*ms, func() { close(stop) })
		start := time.Now()
		var at []time.Duration
		n := Every(100*ms, stop, func(int) {
			at = append(at, time.Since(start))
			time.Sleep(30 * ms) // f takes a while, but less than the interval
		})
		want := []time.Duration{100 * ms, 200 * ms, 300 * ms, 400 * ms, 500 * ms}
		if n != 5 || !slices.Equal(at, want) {
			t.Errorf("Every(100ms, stop closed at 550ms, f taking 30ms) = %d with calls at %v, want 5 with calls at %v", n, at, want)
		}
		if got := time.Since(start); got != 550*ms {
			t.Errorf("Every(100ms, stop closed at 550ms, f taking 30ms) returned after %v, want 550ms", got)
		}
	})
}

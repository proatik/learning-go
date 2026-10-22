package ordered

import (
	"context"
	"fmt"
	"slices"
	"strings"
	"sync/atomic"
	"testing"
	"testing/synctest"
	"time"
)

// closed returns a closed channel holding vals.
func closed[T any](vals ...T) <-chan T {
	ch := make(chan T, len(vals))
	for _, v := range vals {
		ch <- v
	}
	close(ch)
	return ch
}

// received is one result and the (fake) time it arrived.
type received[R any] struct {
	v  R
	at time.Duration
}

// collect reads out until it is closed, inside a bubble (fake time bounds
// the wait).
func collect[R any](t *testing.T, name string, out <-chan R) []received[R] {
	t.Helper()
	if out == nil {
		t.Fatalf("%s returned a nil channel", name)
	}
	start := time.Now()
	var got []received[R]
	for {
		select {
		case v, ok := <-out:
			if !ok {
				return got
			}
			got = append(got, received[R]{v, time.Since(start)})
		case <-time.After(time.Hour):
			t.Fatalf("%s: output not closed (received %d results so far)", name, len(got))
		}
	}
}

func values[R any](rs []received[R]) []R {
	var out []R
	for _, r := range rs {
		out = append(out, r.v)
	}
	return out
}

func times[R any](rs []received[R]) []time.Duration {
	var out []time.Duration
	for _, r := range rs {
		out = append(out, r.at)
	}
	return out
}

func TestProcessResults(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		in := []string{"go", "héllo", "", "世界"}
		got := values(collect(t, "Process", Process(t.Context(), 2, closed(in...), strings.ToUpper)))
		if want := []string{"GO", "HÉLLO", "", "世界"}; !slices.Equal(got, want) {
			t.Errorf("Process(2, %q, ToUpper) = %q, want %q", in, got, want)
		}
		if got := collect(t, "Process", Process(t.Context(), 3, closed[int](), func(n int) int { return n })); len(got) != 0 {
			t.Errorf("Process(3, [], id) = %v, want nothing", values(got))
		}
	})
}

// TestProcessKeepsOrder: item 0 is the slowest and item 5 the fastest, all
// running at once. Results still come out in input order, the first one
// as soon as item 0 is done.
func TestProcessKeepsOrder(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		in := []int{6, 5, 4, 3, 2, 1} // item i sleeps in[i] seconds
		got := collect(t, "Process", Process(t.Context(), 6, closed(in...), func(n int) int {
			time.Sleep(time.Duration(n) * time.Second)
			return n * 10
		}))
		if want := []int{60, 50, 40, 30, 20, 10}; !slices.Equal(values(got), want) {
			t.Fatalf("Process(6, sleeps %v) = %v, want %v", in, values(got), want)
		}
		for i, at := range times(got) {
			if at != 6*time.Second {
				t.Errorf("Process(6, sleeps %v): result %d arrived at %v, want 6s (all run at once)", in, i, at)
			}
		}
	})
}

// TestProcessStreams: results are sent as soon as they (and all earlier
// ones) are ready, not after the whole input is done.
func TestProcessStreams(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		got := collect(t, "Process", Process(t.Context(), 2, closed(1, 1, 1, 1, 1), func(n int) int {
			time.Sleep(time.Second)
			return n
		}))
		want := []time.Duration{1 * time.Second, 1 * time.Second, 2 * time.Second, 2 * time.Second, 3 * time.Second}
		if !slices.Equal(times(got), want) {
			t.Errorf("Process(2, 5 items of 1s): results arrived at %v, want %v", times(got), want)
		}
	})
}

func TestProcessLimit(t *testing.T) {
	tests := []struct {
		workers, items int
		wantPeak       int64
		wantElapsed    time.Duration
	}{
		{3, 7, 3, 3 * time.Second},
		{1, 3, 1, 3 * time.Second},
		{8, 4, 4, 1 * time.Second},
		{0, 2, 1, 2 * time.Second},
		{-5, 2, 1, 2 * time.Second},
	}
	for _, tt := range tests {
		synctest.Test(t, func(t *testing.T) {
			name := fmt.Sprintf("Process(%d workers, %d items of 1s)", tt.workers, tt.items)
			in := make([]int, tt.items)
			for i := range in {
				in[i] = i
			}
			var running, peak atomic.Int64
			start := time.Now()
			got := collect(t, name, Process(t.Context(), tt.workers, closed(in...), func(n int) int {
				now := running.Add(1)
				for p := peak.Load(); now > p && !peak.CompareAndSwap(p, now); p = peak.Load() {
				}
				time.Sleep(time.Second)
				running.Add(-1)
				return -n
			}))
			elapsed := time.Since(start)
			want := make([]int, tt.items)
			for i := range want {
				want[i] = -i
			}
			if !slices.Equal(values(got), want) {
				t.Errorf("%s = %v, want %v", name, values(got), want)
			}
			if p := peak.Load(); p != tt.wantPeak {
				t.Errorf("%s max concurrent calls = %d, want %d", name, p, tt.wantPeak)
			}
			if elapsed != tt.wantElapsed {
				t.Errorf("%s took %v, want %v", name, elapsed, tt.wantElapsed)
			}
		})
	}
}

// release closes in and drains out, so a Process that ignores ctx fails
// with a message instead of leaving a blocked goroutine in the bubble.
func release(in chan int, out <-chan int) {
	close(in)
	go func() {
		for range out {
		}
	}()
}

// TestProcessCancel: in stays open and the reader stops reading. Cancelling
// ctx must close the output and stop every goroutine (synctest fails the
// test if one is left blocked).
func TestProcessCancel(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		ctx, cancel := context.WithCancel(t.Context())
		in := make(chan int, 10)
		for i := range 10 {
			in <- i
		}
		var calls atomic.Int64
		out := Process(ctx, 2, in, func(n int) int {
			calls.Add(1)
			time.Sleep(time.Second)
			return n
		})
		if out == nil {
			t.Fatalf("Process returned a nil channel")
		}
		defer release(in, out)
		select {
		case v := <-out:
			if v != 0 {
				t.Errorf("Process: first result = %d, want 0", v)
			}
		case <-time.After(time.Hour):
			t.Fatalf("Process: no first result")
		}
		time.Sleep(5 * time.Second) // nobody reads: the pipeline must stall, not run ahead
		synctest.Wait()
		if c := calls.Load(); c > 6 {
			t.Errorf("Process(2 workers) with a stalled reader made %d calls, want at most 6 (bound how far it runs ahead)", c)
		}
		cancel()
		synctest.Wait()
		select {
		case v, ok := <-out:
			if ok {
				t.Fatalf("after cancel: output received %d, want it closed", v)
			}
		default:
			t.Fatalf("after cancel: output still open, want it closed (watch ctx.Done() at every blocking step)")
		}
	})
}

// TestProcessCancelWhileWaitingForInput: cancel while in is open and empty.
func TestProcessCancelWhileWaitingForInput(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		ctx, cancel := context.WithCancel(t.Context())
		in := make(chan int)
		out := Process(ctx, 3, in, func(n int) int { return n })
		if out == nil {
			t.Fatalf("Process returned a nil channel")
		}
		defer release(in, out)
		synctest.Wait()
		cancel()
		synctest.Wait()
		select {
		case v, ok := <-out:
			if ok {
				t.Fatalf("after cancel: output received %d, want it closed", v)
			}
		default:
			t.Fatalf("after cancel: output still open, want it closed")
		}
	})
}

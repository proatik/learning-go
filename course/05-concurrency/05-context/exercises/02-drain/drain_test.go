package drain

import (
	"context"
	"errors"
	"slices"
	"testing"
	"testing/synctest"
	"time"
)

// filled returns a buffered channel holding jobs, closed if closeIt is true.
func filled(closeIt bool, jobs ...string) chan string {
	ch := make(chan string, len(jobs))
	for _, j := range jobs {
		ch <- j
	}
	if closeIt {
		close(ch)
	}
	return ch
}

func TestDrainClosed(t *testing.T) {
	tests := []struct {
		jobs []string
	}{
		{[]string{"a", "b", "c"}},
		{[]string{"héllo", "", "世界", "x"}},
		{nil},
	}
	for _, tt := range tests {
		var got []string
		n, err := Drain(t.Context(), filled(true, tt.jobs...), func(s string) { got = append(got, s) })
		if n != len(tt.jobs) || err != nil {
			t.Errorf("Drain(closed %q) = (%d, %v), want (%d, <nil>)", tt.jobs, n, err, len(tt.jobs))
		}
		if !slices.Equal(got, tt.jobs) {
			t.Errorf("Drain(closed %q) handled %q, want %q", tt.jobs, got, tt.jobs)
		}
	}
}

func TestDrainAlreadyCanceled(t *testing.T) {
	// Repeated because select picks randomly: a loop without the up-front
	// ctx check handles some jobs on some runs.
	for range 30 {
		ctx, cancel := context.WithCancel(t.Context())
		cancel()
		handled := 0
		n, err := Drain(ctx, filled(false, "a", "b", "c", "d", "e"), func(string) { handled++ })
		if n != 0 || !errors.Is(err, context.Canceled) || handled != 0 {
			t.Fatalf("Drain(canceled ctx, 5 waiting jobs) = (%d, %v) with %d handled, want (0, %v) with 0 handled", n, err, handled, context.Canceled)
		}
	}
}

func TestDrainCancelDuringJob(t *testing.T) {
	for range 30 {
		ctx, cancel := context.WithCancel(t.Context())
		var got []string
		handle := func(s string) {
			got = append(got, s)
			if len(got) == 2 {
				cancel()
			}
		}
		n, err := Drain(ctx, filled(true, "a", "b", "c", "d", "e"), handle)
		cancel()
		if n != 2 || !errors.Is(err, context.Canceled) {
			t.Fatalf("Drain(cancel while handling job 2 of 5) = (%d, %v), want (2, %v)", n, err, context.Canceled)
		}
		if !slices.Equal(got, []string{"a", "b"}) {
			t.Fatalf("Drain(cancel while handling job 2 of 5) handled %q, want [\"a\" \"b\"]", got)
		}
	}
}

func TestDrainTimeoutWhileWaiting(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		ctx, cancel := context.WithTimeout(t.Context(), 2500*time.Millisecond)
		defer cancel()
		in := make(chan string, 2) // buffered, so these sends never block
		time.AfterFunc(time.Second, func() { in <- "x" })
		time.AfterFunc(2*time.Second, func() { in <- "y" })
		var got []string
		start := time.Now()
		n, err := Drain(ctx, in, func(s string) { got = append(got, s) })
		elapsed := time.Since(start)
		if n != 2 || !errors.Is(err, context.DeadlineExceeded) {
			t.Errorf("Drain(jobs at 1s and 2s, timeout 2.5s) = (%d, %v), want (2, %v)", n, err, context.DeadlineExceeded)
		}
		if !slices.Equal(got, []string{"x", "y"}) {
			t.Errorf("Drain(jobs at 1s and 2s, timeout 2.5s) handled %q, want [\"x\" \"y\"]", got)
		}
		if elapsed != 2500*time.Millisecond {
			t.Errorf("Drain(jobs at 1s and 2s, timeout 2.5s) returned after %v, want 2.5s", elapsed)
		}
	})
}

func TestDrainNilChannel(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		ctx, cancel := context.WithCancel(t.Context())
		defer cancel()
		time.AfterFunc(time.Second, cancel)
		start := time.Now()
		n, err := Drain(ctx, nil, func(string) {})
		elapsed := time.Since(start)
		if n != 0 || !errors.Is(err, context.Canceled) {
			t.Errorf("Drain(nil channel, cancel at 1s) = (%d, %v), want (0, %v)", n, err, context.Canceled)
		}
		if elapsed != time.Second {
			t.Errorf("Drain(nil channel, cancel at 1s) returned after %v, want 1s", elapsed)
		}
	})
}

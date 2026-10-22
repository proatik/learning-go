package stream

import (
	"fmt"
	"runtime"
	"slices"
	"testing"
	"time"
)

// start calls f and fails unless it returns a non-nil channel within 2s.
func start(t *testing.T, call string, f func() <-chan int) <-chan int {
	t.Helper()
	got := make(chan (<-chan int), 1)
	go func() { got <- f() }()
	select {
	case ch := <-got:
		if ch == nil {
			t.Fatalf("%s returned a nil channel", call)
		}
		return ch
	case <-time.After(2 * time.Second):
		t.Fatalf("%s did not return within 2s, want it to return at once and work in a goroutine", call)
		return nil
	}
}

// take receives exactly n values from ch, failing after 2s or if ch closes.
func take(t *testing.T, call string, ch <-chan int, n int) []int {
	t.Helper()
	var got []int
	timeout := time.After(2 * time.Second)
	for len(got) < n {
		select {
		case v, ok := <-ch:
			if !ok {
				t.Fatalf("%s closed after %v, want %d values", call, got, n)
			}
			got = append(got, v)
		case <-timeout:
			t.Fatalf("%s: got %v within 2s, want %d values", call, got, n)
		}
	}
	return got
}

// drain reads ch until it is closed and returns the values, failing if it
// takes over 2s or more than limit values arrive.
func drain(t *testing.T, call string, ch <-chan int, limit int) []int {
	t.Helper()
	var got []int
	timeout := time.After(2 * time.Second)
	for {
		select {
		case v, ok := <-ch:
			if !ok {
				return got
			}
			got = append(got, v)
			if len(got) > limit {
				t.Fatalf("%s sent %d+ values, want it to close (is every send inside a select with <-done?)", call, len(got))
			}
		case <-timeout:
			t.Fatalf("%s: channel not closed within 2s (received %v)", call, got)
		}
	}
}

// settle waits up to 1s for the goroutine count to drop back to before and
// returns how many extra goroutines are left.
func settle(before int) int {
	deadline := time.Now().Add(time.Second)
	for runtime.NumGoroutine() > before && time.Now().Before(deadline) {
		time.Sleep(5 * time.Millisecond)
	}
	return runtime.NumGoroutine() - before
}

// buffered returns a channel holding vals; it is closed only if closeIt.
func buffered(closeIt bool, vals ...int) chan int {
	ch := make(chan int, len(vals))
	for _, v := range vals {
		ch <- v
	}
	if closeIt {
		close(ch)
	}
	return ch
}

func TestFrom(t *testing.T) {
	before := runtime.NumGoroutine()
	done := make(chan struct{})
	ch := start(t, "From(done, 10)", func() <-chan int { return From(done, 10) })
	if got := take(t, "From(done, 10)", ch, 5); !slices.Equal(got, []int{10, 11, 12, 13, 14}) {
		t.Fatalf("From(done, 10) = %v, want [10 11 12 13 14 …]", got)
	}
	close(done)
	rest := drain(t, "From(done, 10) after close(done)", ch, 100)
	for i, v := range rest {
		if v != 15+i {
			t.Fatalf("From(done, 10) after close(done) sent %v, want it to continue the sequence (15, 16, …) or stop", rest)
		}
	}
	if n := settle(before); n > 0 {
		t.Errorf("From(done, 10) after close(done): %d goroutine(s) still running, want 0", n)
	}
}

func TestFromNoReader(t *testing.T) {
	for _, closeFirst := range []bool{false, true} {
		before := runtime.NumGoroutine()
		done := make(chan struct{})
		if closeFirst {
			close(done)
		}
		call := "From(done, 0) with nobody reading"
		if closeFirst {
			call = "From(done, 0) with done already closed"
		}
		ch := start(t, call, func() <-chan int { return From(done, 0) })
		if !closeFirst {
			close(done)
		}
		if n := settle(before); n > 0 {
			t.Fatalf("%s, then close(done): %d goroutine(s) still running after 1s, want the generator to exit", call, n)
		}
		if got := drain(t, call, ch, 0); len(got) != 0 {
			t.Errorf("%s: received %v, want the channel closed with no values", call, got)
		}
	}
}

func TestTake(t *testing.T) {
	tests := []struct {
		in       []int
		closeIn  bool
		n        int
		want     []int
		leftInIn int
	}{
		{[]int{7, 8, 9, 10}, false, 3, []int{7, 8, 9}, 1},
		{[]int{7, 8, 9}, false, 3, []int{7, 8, 9}, 0},
		{[]int{7, 8}, true, 10, []int{7, 8}, 0},
		{[]int{0, 0, 5}, true, 2, []int{0, 0}, 1},
		{[]int{1, 2, 3, 4, 5}, false, 0, nil, 5},
		{[]int{1, 2}, false, -1, nil, 2},
		{nil, true, 3, nil, 0},
	}
	for _, tt := range tests {
		before := runtime.NumGoroutine()
		done := make(chan struct{}) // never closed: Take must finish on its own
		in := buffered(tt.closeIn, tt.in...)
		call := fmt.Sprintf("Take(done, %v, %d)", tt.in, tt.n)
		got := drain(t, call, start(t, call, func() <-chan int { return Take(done, in, tt.n) }), 1000)
		if !slices.Equal(got, tt.want) {
			t.Errorf("%s = %v, want %v", call, got, tt.want)
			continue
		}
		if len(in) != tt.leftInIn {
			t.Errorf("%s left %d values in in, want %d (don't receive more than n)", call, len(in), tt.leftInIn)
		}
		if n := settle(before); n > 0 {
			t.Errorf("%s: %d goroutine(s) still running after the output closed, want 0", call, n)
		}
	}
}

func TestTakeStopsOnDone(t *testing.T) {
	// Blocked receiving: in never delivers anything.
	before := runtime.NumGoroutine()
	done := make(chan struct{})
	idle := make(chan int)
	out := start(t, "Take(done, idle, 5)", func() <-chan int { return Take(done, idle, 5) })
	close(done)
	if got := drain(t, "Take(done, idle, 5) after close(done)", out, 0); len(got) != 0 {
		t.Errorf("Take(done, idle, 5) after close(done) = %v, want no values", got)
	}
	if n := settle(before); n > 0 {
		t.Errorf("Take(done, idle, 5) after close(done): %d goroutine(s) still running, want 0 (watch done while receiving from in)", n)
	}

	// Blocked sending: in has values but nobody reads the output.
	before = runtime.NumGoroutine()
	done = make(chan struct{})
	in := buffered(false, 1, 2, 3)
	start(t, "Take(done, [1 2 3], 3)", func() <-chan int { return Take(done, in, 3) })
	close(done)
	if n := settle(before); n > 0 {
		t.Errorf("Take(done, [1 2 3], 3) with no reader, then close(done): %d goroutine(s) still running, want 0 (watch done while sending)", n)
	}
}

func TestPipeline(t *testing.T) {
	tests := []struct {
		from, n int
		want    []int
	}{
		{1, 3, []int{1, 2, 3}},
		{-2, 4, []int{-2, -1, 0, 1}},
		{100, 1, []int{100}},
		{5, 0, nil},
	}
	for _, tt := range tests {
		before := runtime.NumGoroutine()
		done := make(chan struct{})
		call := fmt.Sprintf("Take(done, From(done, %d), %d)", tt.from, tt.n)
		out := start(t, call, func() <-chan int { return Take(done, From(done, tt.from), tt.n) })
		if got := drain(t, call, out, 1000); !slices.Equal(got, tt.want) {
			t.Errorf("%s = %v, want %v", call, got, tt.want)
		}
		close(done) // From is still blocked on its next send until now
		if n := settle(before); n > 0 {
			t.Errorf("%s, then close(done): %d goroutine(s) still running, want 0", call, n)
		}
	}
}

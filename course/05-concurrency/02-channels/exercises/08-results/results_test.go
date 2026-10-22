package results

import (
	"errors"
	"fmt"
	"runtime"
	"slices"
	"testing"
	"time"
)

var errBad = errors.New("bad")

func double(n int) (int, error) { return 2 * n, nil }

func failOn3(n int) (int, error) {
	if n == 3 {
		return 0, fmt.Errorf("input %d: %w", n, errBad)
	}
	return n, nil
}

// within runs f and fails if it doesn't return within 2s.
func within[R any](t *testing.T, call string, f func() R) R {
	t.Helper()
	done := make(chan R, 1)
	go func() { done <- f() }()
	select {
	case r := <-done:
		return r
	case <-time.After(2 * time.Second):
		t.Fatalf("%s did not return within 2s", call)
		panic("unreachable")
	}
}

// run calls Run and reads every result until the channel is closed.
func run(t *testing.T, call string, inputs []int, f func(int) (int, error)) []Result {
	t.Helper()
	ch := within(t, call, func() <-chan Result { return Run(inputs, f) })
	if ch == nil {
		t.Fatalf("%s returned a nil channel", call)
	}
	if cap(ch) != len(inputs) {
		t.Fatalf("cap(%s) = %d, want %d (one slot per input, so the worker never blocks)", call, cap(ch), len(inputs))
	}
	var got []Result
	timeout := time.After(2 * time.Second)
	for {
		select {
		case r, ok := <-ch:
			if !ok {
				return got
			}
			got = append(got, r)
		case <-timeout:
			t.Fatalf("%s: channel not closed within 2s (received %v so far)", call, got)
		}
	}
}

func TestRun(t *testing.T) {
	got := run(t, "Run([1 2 3], double)", []int{1, 2, 3}, double)
	if want := []Result{{2, nil}, {4, nil}, {6, nil}}; !slices.Equal(got, want) {
		t.Errorf("Run([1 2 3], double) = %v, want %v", got, want)
	}

	got = run(t, "Run([], double)", nil, double)
	if len(got) != 0 {
		t.Errorf("Run([], double) = %v, want no results", got)
	}
}

func TestRunStopsAtFirstError(t *testing.T) {
	var calls []int // only the worker goroutine appends; read after the channel is closed
	f := func(n int) (int, error) {
		calls = append(calls, n)
		return failOn3(n)
	}
	got := run(t, "Run([1 3 5 3], failOn3)", []int{1, 3, 5, 3}, f)
	if len(got) != 2 || got[0] != (Result{1, nil}) || !errors.Is(got[1].Err, errBad) {
		t.Fatalf("Run([1 3 5 3], failOn3) = %v, want [{1 <nil>} {0 input 3: bad}]", got)
	}
	if !slices.Equal(calls, []int{1, 3}) {
		t.Errorf("Run([1 3 5 3], failOn3) called f with %v, want [1 3] (stop after the first error)", calls)
	}
}

func TestRunNeverBlocks(t *testing.T) {
	before := runtime.NumGoroutine()
	inputs := make([]int, 100)
	ch := within(t, "Run(100 inputs, double)", func() <-chan Result { return Run(inputs, double) })
	if ch == nil {
		t.Fatalf("Run(100 inputs, double) returned a nil channel")
	}
	// Nobody reads: the worker must still finish.
	deadline := time.Now().Add(time.Second)
	for runtime.NumGoroutine() > before && time.Now().Before(deadline) {
		time.Sleep(5 * time.Millisecond)
	}
	if n := runtime.NumGoroutine() - before; n > 0 {
		t.Fatalf("Run(100 inputs, double) with no reader: %d goroutine(s) still running after 1s, want the worker to exit", n)
	}
	if len(ch) != 100 {
		t.Errorf("len(Run(100 inputs, double)) with no reader = %d, want 100 buffered results", len(ch))
	}
}

// closed returns a buffered channel holding rs, already closed.
func closed(rs ...Result) chan Result {
	ch := make(chan Result, len(rs))
	for _, r := range rs {
		ch <- r
	}
	close(ch)
	return ch
}

type collected struct {
	vals []int
	err  error
}

func collect(t *testing.T, call string, ch <-chan Result) collected {
	t.Helper()
	return within(t, call, func() collected {
		v, err := Collect(ch)
		return collected{v, err}
	})
}

func TestCollect(t *testing.T) {
	got := collect(t, "Collect([{2} {4} {6}])", closed(Result{2, nil}, Result{4, nil}, Result{6, nil}))
	if !slices.Equal(got.vals, []int{2, 4, 6}) || got.err != nil {
		t.Errorf("Collect([{2} {4} {6}]) = (%v, %v), want ([2 4 6], <nil>)", got.vals, got.err)
	}

	got = collect(t, "Collect([{0} {0}])", closed(Result{0, nil}, Result{0, nil}))
	if !slices.Equal(got.vals, []int{0, 0}) || got.err != nil {
		t.Errorf("Collect([{0} {0}]) = (%v, %v), want ([0 0], <nil>)", got.vals, got.err)
	}

	got = collect(t, "Collect([])", closed())
	if len(got.vals) != 0 || got.err != nil {
		t.Errorf("Collect([]) = (%v, %v), want ([], <nil>)", got.vals, got.err)
	}
}

func TestCollectStopsAtFirstError(t *testing.T) {
	wrapped := fmt.Errorf("job 2: %w", errBad)
	ch := closed(Result{1, nil}, Result{0, wrapped}, Result{3, nil}, Result{0, errors.New("later")})
	got := collect(t, "Collect([{1} {err} {3} {err}])", ch)
	if !slices.Equal(got.vals, []int{1}) || !errors.Is(got.err, errBad) {
		t.Errorf("Collect([{1} {job 2: bad} {3} {later}]) = (%v, %v), want ([1], job 2: bad)", got.vals, got.err)
	}
	if len(ch) != 2 {
		t.Errorf("Collect([{1} {err} {3} {err}]) left %d results unread, want 2 (stop reading at the first error)", len(ch))
	}

	got = collect(t, "Collect([{err}])", closed(Result{7, errBad}))
	if len(got.vals) != 0 || !errors.Is(got.err, errBad) {
		t.Errorf("Collect([{7 bad}]) = (%v, %v), want ([], bad)", got.vals, got.err)
	}
}

func TestCollectNil(t *testing.T) {
	got := collect(t, "Collect(nil)", nil)
	if got.vals != nil || got.err != nil {
		t.Errorf("Collect(nil) = (%v, %v), want ([], <nil>)", got.vals, got.err)
	}
}

func TestRunThenCollect(t *testing.T) {
	ch := within(t, "Run", func() <-chan Result { return Run([]int{1, 3, 5}, failOn3) })
	if ch == nil {
		t.Fatalf("Run([1 3 5], failOn3) returned a nil channel")
	}
	got := collect(t, "Collect(Run([1 3 5], failOn3))", ch)
	if !slices.Equal(got.vals, []int{1}) || !errors.Is(got.err, errBad) {
		t.Errorf("Collect(Run([1 3 5], failOn3)) = (%v, %v), want ([1], input 3: bad)", got.vals, got.err)
	}
}

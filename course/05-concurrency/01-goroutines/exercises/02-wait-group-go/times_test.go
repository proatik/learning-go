package times_test

import (
	"sync/atomic"
	"testing"
	"testing/synctest"
	"time"

	times "learngo/course/05-concurrency/01-goroutines/exercises/02-wait-group-go"
)

// within runs f in a goroutine and fails the test if it doesn't return in 2s.
func within(t *testing.T, name string, f func()) {
	t.Helper()
	done := make(chan struct{})
	go func() {
		defer close(done)
		f()
	}()
	select {
	case <-done:
	case <-time.After(2 * time.Second):
		t.Fatalf("%s did not return within 2s", name)
	}
}

// recorder counts how often f is called with each i in [0, n), and how often
// with an i outside that range.
type recorder struct {
	calls []atomic.Int32
	bad   atomic.Int32
}

func newRecorder(n int) *recorder {
	return &recorder{calls: make([]atomic.Int32, max(n, 0))}
}

func (r *recorder) record(i int) {
	if i < 0 || i >= len(r.calls) {
		r.bad.Add(1)
		return
	}
	r.calls[i].Add(1)
}

// check reports every i not called exactly once.
func (r *recorder) check(t *testing.T, n int) {
	t.Helper()
	if b := r.bad.Load(); b != 0 {
		t.Errorf("Times(%d, f): f called with an out-of-range i %d times, want 0", n, b)
	}
	for i := range r.calls {
		if c := r.calls[i].Load(); c != 1 {
			t.Errorf("Times(%d, f): f(%d) called %d times, want 1", n, i, c)
		}
	}
}

func TestConcurrent(t *testing.T) {
	for _, n := range []int{1, 2, 4, 25} {
		synctest.Test(t, func(t *testing.T) {
			r := newRecorder(n)
			start := time.Now()
			got := times.Times(n, func(i int) {
				time.Sleep(time.Second)
				r.record(i)
			})
			elapsed := time.Since(start)
			if got != n {
				t.Errorf("Times(%d, f) = %d, want %d", n, got, n)
			}
			r.check(t, n)
			if elapsed != time.Second {
				t.Errorf("Times(%d, f of 1s) took %v, want 1s", n, elapsed)
			}
		})
	}
}

func TestNonPositive(t *testing.T) {
	for _, n := range []int{0, -1, -3} {
		var calls atomic.Int32
		var got int
		within(t, "Times", func() {
			got = times.Times(n, func(int) { calls.Add(1) })
		})
		if got != 0 {
			t.Errorf("Times(%d, f) = %d, want 0", n, got)
		}
		if c := calls.Load(); c != 0 {
			t.Errorf("Times(%d, f): f called %d times, want 0", n, c)
		}
	}
}

func TestMany(t *testing.T) {
	const n = 1000
	r := newRecorder(n)
	var got int
	within(t, "Times(1000, f)", func() { got = times.Times(n, r.record) })
	if got != n {
		t.Errorf("Times(%d, f) = %d, want %d", n, got, n)
	}
	r.check(t, n)
}

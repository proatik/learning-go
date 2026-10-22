package breaker

import (
	"errors"
	"sync"
	"sync/atomic"
	"testing"
	"testing/synctest"
	"time"
)

var errBoom = errors.New("boom")

// counter returns a function that counts its calls and returns err.
func counter(err error) (func() error, *atomic.Int64) {
	var n atomic.Int64
	return func() error { n.Add(1); return err }, &n
}

func wantState(t *testing.T, step string, b *Breaker, want State) {
	t.Helper()
	if got := b.State(); got != want {
		t.Fatalf("%s: State() = %v, want %v", step, got, want)
	}
}

// open drives a fresh breaker into the Open state.
func open(t *testing.T, threshold int, cooldown time.Duration) *Breaker {
	t.Helper()
	b := New(threshold, cooldown)
	fail, _ := counter(errBoom)
	for i := range threshold {
		if err := b.Call(fail); !errors.Is(err, errBoom) {
			t.Fatalf("Call(fail) #%d error = %v, want %v", i+1, err, errBoom)
		}
	}
	wantState(t, "after threshold failures", b, Open)
	return b
}

func TestStartsClosed(t *testing.T) {
	b := New(3, time.Second)
	wantState(t, "New(3, 1s)", b, Closed)
	ok, calls := counter(nil)
	if err := b.Call(ok); err != nil {
		t.Errorf("Call(ok) error = %v, want <nil>", err)
	}
	if calls.Load() != 1 {
		t.Errorf("Call(ok) ran fn %d times, want 1", calls.Load())
	}
	wantState(t, "Call(ok)", b, Closed)
}

func TestOpensAfterThreshold(t *testing.T) {
	b := New(3, 10*time.Second)
	fail, _ := counter(errBoom)
	for i := range 2 {
		if err := b.Call(fail); !errors.Is(err, errBoom) {
			t.Fatalf("Call(fail) #%d error = %v, want %v", i+1, err, errBoom)
		}
	}
	wantState(t, "2 failures, threshold 3", b, Closed)
	b.Call(fail)
	wantState(t, "3 failures, threshold 3", b, Open)
	ok, calls := counter(nil)
	if err := b.Call(ok); !errors.Is(err, ErrOpen) {
		t.Errorf("Call(ok) while open: error = %v, want %v", err, ErrOpen)
	}
	if calls.Load() != 0 {
		t.Errorf("Call(ok) while open ran fn %d times, want 0", calls.Load())
	}
}

func TestThresholdBelowOne(t *testing.T) {
	b := New(0, time.Second)
	fail, _ := counter(errBoom)
	b.Call(fail)
	wantState(t, "New(0, 1s), 1 failure", b, Open)
}

func TestSuccessResetsCount(t *testing.T) {
	b := New(3, time.Second)
	fail, _ := counter(errBoom)
	ok, _ := counter(nil)
	b.Call(fail)
	b.Call(fail)
	b.Call(ok)
	b.Call(fail)
	b.Call(fail)
	wantState(t, "fail, fail, ok, fail, fail (threshold 3)", b, Closed)
	b.Call(fail)
	wantState(t, "fail, fail, ok, fail, fail, fail (threshold 3)", b, Open)
}

func TestCooldown(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		b := open(t, 2, 10*time.Second)
		ok, calls := counter(nil)
		time.Sleep(10*time.Second - time.Millisecond)
		if err := b.Call(ok); !errors.Is(err, ErrOpen) {
			t.Fatalf("Call(ok) 9.999s after opening (cooldown 10s): error = %v, want %v", err, ErrOpen)
		}
		wantState(t, "9.999s after opening", b, Open)
		time.Sleep(time.Millisecond)
		if err := b.Call(ok); err != nil {
			t.Fatalf("Call(ok) 10s after opening: error = %v, want <nil> (trial call)", err)
		}
		if calls.Load() != 1 {
			t.Fatalf("fn ran %d times, want 1 (only the trial call)", calls.Load())
		}
		wantState(t, "successful trial call", b, Closed)
	})
}

func TestTrialFailureReopens(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		b := open(t, 1, 10*time.Second)
		time.Sleep(10 * time.Second)
		fail, _ := counter(errBoom)
		if err := b.Call(fail); !errors.Is(err, errBoom) {
			t.Fatalf("trial Call(fail): error = %v, want %v", err, errBoom)
		}
		wantState(t, "failed trial call", b, Open)
		ok, calls := counter(nil)
		time.Sleep(9 * time.Second)
		if err := b.Call(ok); !errors.Is(err, ErrOpen) {
			t.Fatalf("Call(ok) 9s after the failed trial (cooldown restarts): error = %v, want %v", err, ErrOpen)
		}
		time.Sleep(time.Second)
		if err := b.Call(ok); err != nil {
			t.Fatalf("Call(ok) 10s after the failed trial: error = %v, want <nil>", err)
		}
		if calls.Load() != 1 {
			t.Fatalf("fn ran %d times, want 1", calls.Load())
		}
		wantState(t, "second trial succeeded", b, Closed)
	})
}

func TestOnlyOneTrial(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		b := open(t, 1, time.Second)
		time.Sleep(time.Second)
		trialDone := make(chan error, 1)
		go func() {
			trialDone <- b.Call(func() error { time.Sleep(time.Second); return nil })
		}()
		synctest.Wait()
		wantState(t, "while the trial call runs", b, HalfOpen)
		ok, calls := counter(nil)
		if err := b.Call(ok); !errors.Is(err, ErrOpen) {
			t.Errorf("Call(ok) while the trial call runs: error = %v, want %v", err, ErrOpen)
		}
		if calls.Load() != 0 {
			t.Errorf("Call(ok) while the trial call runs ran fn %d times, want 0", calls.Load())
		}
		time.Sleep(time.Second)
		synctest.Wait()
		select {
		case err := <-trialDone:
			if err != nil {
				t.Errorf("trial call error = %v, want <nil>", err)
			}
		default:
			t.Fatalf("trial call did not finish")
		}
		wantState(t, "after the trial call", b, Closed)
	})
}

func TestClosedCallsRunConcurrently(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		b := New(3, time.Second)
		var calls atomic.Int64
		start := time.Now()
		var wg sync.WaitGroup
		for range 10 {
			wg.Go(func() {
				b.Call(func() error { calls.Add(1); time.Sleep(time.Second); return nil })
			})
		}
		wg.Wait()
		if n := calls.Load(); n != 10 {
			t.Fatalf("10 concurrent Calls ran fn %d times, want 10", n)
		}
		if el := time.Since(start); el != time.Second {
			t.Errorf("10 concurrent Calls of 1s each took %v, want 1s (don't hold the lock while fn runs)", el)
		}
	})
}

func TestConcurrentFailuresCounted(t *testing.T) {
	const workers, each = 50, 200
	b := New(workers*each+1, time.Second)
	fail, calls := counter(errBoom)
	var wg sync.WaitGroup
	for range workers {
		wg.Go(func() {
			for range each {
				b.Call(fail)
			}
		})
	}
	wg.Wait()
	if n := calls.Load(); n != workers*each {
		t.Fatalf("fn ran %d times, want %d", n, workers*each)
	}
	wantState(t, "10000 concurrent failures, threshold 10001", b, Closed)
	b.Call(fail)
	wantState(t, "10001 failures, threshold 10001 (lost updates?)", b, Open)
}

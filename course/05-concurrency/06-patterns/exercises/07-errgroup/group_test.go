package group

import (
	"context"
	"errors"
	"fmt"
	"sync/atomic"
	"testing"
	"testing/synctest"
	"time"
)

var (
	errA = errors.New("task a failed")
	errB = errors.New("task b failed")
)

// sleepThen returns a task that sleeps d (fake time) and returns err.
func sleepThen(d time.Duration, err error, ran *atomic.Int64) func() error {
	return func() error {
		time.Sleep(d)
		ran.Add(1)
		return err
	}
}

func TestZeroGroup(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		var g Group
		var ran atomic.Int64
		for range 5 {
			g.Go(sleepThen(time.Second, nil, &ran))
		}
		start := time.Now()
		if err := g.Wait(); err != nil {
			t.Errorf("Wait() with 5 nil tasks = %v, want <nil>", err)
		}
		if n := ran.Load(); n != 5 {
			t.Errorf("Wait() returned after %d of 5 tasks finished, want 5", n)
		}
		if d := time.Since(start); d != time.Second {
			t.Errorf("Wait() for 5 tasks of 1s took %v, want 1s (run each task in its own goroutine)", d)
		}

		var g2 Group
		g2.Go(sleepThen(0, errA, &ran))
		if err := g2.Wait(); err != errA {
			t.Errorf("zero Group: Wait() = %v, want %v", err, errA)
		}
		var empty Group
		if err := empty.Wait(); err != nil {
			t.Errorf("Wait() with no tasks = %v, want <nil>", err)
		}
	})
}

// TestFirstErrorWins: b fails at 1s, a at 2s, c succeeds at 3s. Wait still
// waits for c and returns b's error, unchanged.
func TestFirstErrorWins(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		g, _ := WithContext(t.Context())
		var ran atomic.Int64
		g.Go(sleepThen(2*time.Second, errA, &ran))
		g.Go(sleepThen(1*time.Second, errB, &ran))
		g.Go(sleepThen(3*time.Second, nil, &ran))
		start := time.Now()
		err := g.Wait()
		if err != errB {
			t.Errorf("Wait() = %v, want %v (the first error)", err, errB)
		}
		if n := ran.Load(); n != 3 {
			t.Errorf("Wait() returned after %d of 3 tasks finished, want 3", n)
		}
		if d := time.Since(start); d != 3*time.Second {
			t.Errorf("Wait() took %v, want 3s (wait for every task, even after an error)", d)
		}
	})
}

func TestWrappedErrorKept(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		g, _ := WithContext(t.Context())
		g.Go(func() error { return fmt.Errorf("fetch user 7: %w", errA) })
		err := g.Wait()
		if !errors.Is(err, errA) || err == nil || err.Error() != "fetch user 7: task a failed" {
			t.Errorf("Wait() = %v, want %q (wrapping errA)", err, "fetch user 7: task a failed")
		}
	})
}

// TestErrorCancelsContext: the slow task watches ctx and stops as soon as
// the fast one fails.
func TestErrorCancelsContext(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		g, ctx := WithContext(t.Context())
		var sawCancel atomic.Bool
		g.Go(func() error {
			select {
			case <-ctx.Done():
				sawCancel.Store(true)
				return ctx.Err()
			case <-time.After(time.Hour):
				return nil
			}
		})
		g.Go(func() error {
			time.Sleep(time.Second)
			return errA
		})
		start := time.Now()
		err := g.Wait()
		if err != errA {
			t.Errorf("Wait() = %v, want %v (the first error, not the cancelled task's)", err, errA)
		}
		if !sawCancel.Load() {
			t.Errorf("the slow task never saw ctx.Done(): the first error must cancel the context")
		}
		if d := time.Since(start); d != time.Second {
			t.Errorf("Wait() took %v, want 1s (the error should cancel the slow task at once)", d)
		}
		if got := ctx.Err(); got != context.Canceled {
			t.Errorf("ctx.Err() after a failure = %v, want %v", got, context.Canceled)
		}
		if got := context.Cause(ctx); got != errA {
			t.Errorf("context.Cause(ctx) after a failure = %v, want %v", got, errA)
		}
	})
}

func TestContextCancelledByWait(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		g, ctx := WithContext(t.Context())
		var ran atomic.Int64
		g.Go(sleepThen(time.Second, nil, &ran))
		time.Sleep(2 * time.Second)
		synctest.Wait()
		if ran.Load() != 1 {
			t.Fatalf("Go(task): task did not run")
		}
		if err := ctx.Err(); err != nil {
			t.Errorf("ctx.Err() before Wait, no failures = %v, want <nil>", err)
		}
		if err := g.Wait(); err != nil {
			t.Errorf("Wait() = %v, want <nil>", err)
		}
		if err := ctx.Err(); err != context.Canceled {
			t.Errorf("ctx.Err() after Wait = %v, want %v (Wait releases the context)", err, context.Canceled)
		}
		if got := context.Cause(ctx); got != context.Canceled {
			t.Errorf("context.Cause(ctx) after a clean Wait = %v, want %v", got, context.Canceled)
		}
	})
}

func TestParentCancel(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		parent, cancel := context.WithCancel(t.Context())
		g, ctx := WithContext(parent)
		g.Go(func() error {
			<-ctx.Done()
			return ctx.Err()
		})
		cancel()
		if err := g.Wait(); err != context.Canceled {
			t.Errorf("Wait() after cancelling the parent = %v, want %v", err, context.Canceled)
		}
	})
}

// TestManyGoroutines: 1000 tasks, 100 of them failing at the same moment.
// Exactly one error is kept and every task runs once.
func TestManyGoroutines(t *testing.T) {
	for range 20 {
		g, _ := WithContext(context.Background())
		var ran atomic.Int64
		errs := make(map[error]bool)
		for i := range 1000 {
			var err error
			if i%10 == 0 {
				err = fmt.Errorf("task %d", i)
				errs[err] = true
			}
			g.Go(func() error {
				ran.Add(1)
				return err
			})
		}
		done := make(chan error, 1)
		go func() { done <- g.Wait() }()
		select {
		case err := <-done:
			if !errs[err] {
				t.Fatalf("Wait() with 100 failing tasks = %v, want one of their errors", err)
			}
		case <-time.After(2 * time.Second):
			t.Fatalf("Wait() did not return within 2s")
		}
		if n := ran.Load(); n != 1000 {
			t.Fatalf("Go called 1000 times, %d tasks ran, want 1000", n)
		}
	}
}

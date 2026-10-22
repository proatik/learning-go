package server

import (
	"context"
	"errors"
	"fmt"
	"slices"
	"sync"
	"sync/atomic"
	"testing"
	"testing/synctest"
	"time"
)

// recorder is a handle func that sleeps d (fake time in a bubble), then
// records the job. It also tracks how many jobs run at once.
type recorder struct {
	d             time.Duration
	mu            sync.Mutex
	jobs          []int
	running, peak atomic.Int64
}

func (r *recorder) handle(job int) {
	now := r.running.Add(1)
	for p := r.peak.Load(); now > p && !r.peak.CompareAndSwap(p, now); p = r.peak.Load() {
	}
	time.Sleep(r.d)
	r.running.Add(-1)
	r.mu.Lock()
	r.jobs = append(r.jobs, job)
	r.mu.Unlock()
}

// handled returns the recorded jobs, sorted.
func (r *recorder) handled() []int {
	r.mu.Lock()
	defer r.mu.Unlock()
	out := slices.Clone(r.jobs)
	slices.Sort(out)
	return out
}

func submitAll(t *testing.T, s *Server, jobs ...int) {
	t.Helper()
	for _, j := range jobs {
		if err := s.Submit(j); err != nil {
			t.Fatalf("Submit(%d) = %v, want <nil>", j, err)
		}
	}
}

// TestShutdownDrainsQueue: 6 queued jobs of 1s on 2 workers. Shutdown waits
// for all of them (3s), then the server refuses new jobs.
func TestShutdownDrainsQueue(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		r := &recorder{d: time.Second}
		s := New(2, 10, r.handle)
		submitAll(t, s, 1, 2, 3, 4, 5, 6)
		ctx, cancel := context.WithTimeout(t.Context(), 10*time.Second)
		defer cancel()
		start := time.Now()
		if err := s.Shutdown(ctx); err != nil {
			t.Errorf("Shutdown(10s timeout) = %v, want <nil>", err)
		}
		if d := time.Since(start); d != 3*time.Second {
			t.Errorf("Shutdown with 6 jobs of 1s on 2 workers took %v, want 3s", d)
		}
		if got, want := r.handled(), []int{1, 2, 3, 4, 5, 6}; !slices.Equal(got, want) {
			t.Errorf("jobs handled before Shutdown returned = %v, want %v", got, want)
		}
		if p := r.peak.Load(); p != 2 {
			t.Errorf("max jobs running at once = %d, want 2 (one per worker)", p)
		}
		if err := s.Submit(7); !errors.Is(err, ErrClosed) {
			t.Errorf("Submit(7) after Shutdown = %v, want %v", err, ErrClosed)
		}
		if err := s.Shutdown(t.Context()); err != nil {
			t.Errorf("second Shutdown = %v, want <nil>", err)
		}
		if got := r.handled(); len(got) != 6 {
			t.Errorf("jobs handled after a rejected Submit = %v, want 6 jobs", got)
		}
	})
}

// TestShutdownTimeout: the context ends before the queue is drained.
// Shutdown returns ctx.Err() then, and the workers carry on.
func TestShutdownTimeout(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		r := &recorder{d: time.Second}
		s := New(2, 10, r.handle)
		submitAll(t, s, 1, 2, 3, 4, 5, 6)
		ctx, cancel := context.WithTimeout(t.Context(), 1500*time.Millisecond)
		defer cancel()
		start := time.Now()
		if err := s.Shutdown(ctx); !errors.Is(err, context.DeadlineExceeded) {
			t.Errorf("Shutdown(1.5s timeout) = %v, want %v", err, context.DeadlineExceeded)
		}
		if d := time.Since(start); d != 1500*time.Millisecond {
			t.Errorf("Shutdown(1.5s timeout) returned after %v, want 1.5s", d)
		}
		if n := len(r.handled()); n != 2 {
			t.Errorf("jobs handled when Shutdown timed out = %d, want 2", n)
		}
		if err := s.Submit(7); !errors.Is(err, ErrClosed) {
			t.Errorf("Submit(7) after a timed-out Shutdown = %v, want %v", err, ErrClosed)
		}
		if err := s.Shutdown(t.Context()); err != nil {
			t.Errorf("second Shutdown(no timeout) = %v, want <nil>", err)
		}
		if d := time.Since(start); d != 3*time.Second {
			t.Errorf("second Shutdown returned at %v, want 3s (when the queue is drained)", d)
		}
		if got, want := r.handled(), []int{1, 2, 3, 4, 5, 6}; !slices.Equal(got, want) {
			t.Errorf("jobs handled = %v, want %v", got, want)
		}
	})
}

func TestShutdownCancelledContext(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		gate := make(chan struct{})
		var ran atomic.Int64
		s := New(1, 1, func(int) { <-gate; ran.Add(1) })
		submitAll(t, s, 1)
		ctx, cancel := context.WithCancel(t.Context())
		cancel()
		if err := s.Shutdown(ctx); !errors.Is(err, context.Canceled) {
			t.Errorf("Shutdown(cancelled ctx) with a job running = %v, want %v", err, context.Canceled)
		}
		close(gate)
		if err := s.Shutdown(t.Context()); err != nil {
			t.Errorf("Shutdown after the job finished = %v, want <nil>", err)
		}
		if n := ran.Load(); n != 1 {
			t.Errorf("jobs handled = %d, want 1", n)
		}
	})
}

// TestQueueFull: 1 worker busy + 1 queue slot taken means the third job is
// refused at once.
func TestQueueFull(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		gate := make(chan struct{})
		var mu sync.Mutex
		var got []int
		s := New(1, 1, func(j int) {
			<-gate
			mu.Lock()
			got = append(got, j)
			mu.Unlock()
		})
		synctest.Wait() // the worker is idle, waiting for a job
		submitAll(t, s, 1)
		synctest.Wait() // the worker took job 1 and is stuck in handle
		submitAll(t, s, 2)
		if err := s.Submit(3); !errors.Is(err, ErrQueueFull) {
			t.Errorf("Submit(3) with the worker busy and the queue full = %v, want %v", err, ErrQueueFull)
		}
		close(gate)
		if err := s.Shutdown(t.Context()); err != nil {
			t.Errorf("Shutdown = %v, want <nil>", err)
		}
		if want := []int{1, 2}; !slices.Equal(got, want) {
			t.Errorf("jobs handled = %v, want %v", got, want)
		}
	})
}

// TestClamping: New(0, -1, h) has one worker and no queue: a job is accepted
// only while that worker is idle.
func TestClamping(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		r := &recorder{d: time.Second}
		s := New(0, -1, r.handle)
		synctest.Wait()
		submitAll(t, s, 1)
		synctest.Wait()
		if err := s.Submit(2); !errors.Is(err, ErrQueueFull) {
			t.Errorf("New(0, -1): Submit(2) while the only worker is busy = %v, want %v", err, ErrQueueFull)
		}
		time.Sleep(2 * time.Second)
		synctest.Wait()
		submitAll(t, s, 3)
		if err := s.Shutdown(t.Context()); err != nil {
			t.Errorf("Shutdown = %v, want <nil>", err)
		}
		if got, want := r.handled(), []int{1, 3}; !slices.Equal(got, want) {
			t.Errorf("New(0, -1): jobs handled = %v, want %v", got, want)
		}
	})
}

func TestShutdownIdle(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		r := &recorder{}
		s := New(4, 4, r.handle)
		start := time.Now()
		if err := s.Shutdown(t.Context()); err != nil {
			t.Errorf("Shutdown of an idle server = %v, want <nil>", err)
		}
		if d := time.Since(start); d != 0 {
			t.Errorf("Shutdown of an idle server took %v, want 0s", d)
		}
		if err := s.Submit(1); !errors.Is(err, ErrClosed) {
			t.Errorf("Submit(1) after Shutdown = %v, want %v", err, ErrClosed)
		}
	})
}

// TestSubmitDuringShutdown: 20 goroutines submit while Shutdown runs. No
// panic, and exactly the accepted jobs are handled.
func TestSubmitDuringShutdown(t *testing.T) {
	for round := range 20 {
		var handledN atomic.Int64
		s := New(4, 50, func(int) { handledN.Add(1) })
		var accepted atomic.Int64
		var wg sync.WaitGroup
		var badErr atomic.Value
		for g := range 20 {
			wg.Go(func() {
				for i := range 200 {
					switch err := s.Submit(g*1000 + i); {
					case err == nil:
						accepted.Add(1)
					case errors.Is(err, ErrClosed), errors.Is(err, ErrQueueFull):
					default:
						badErr.Store(fmt.Sprint(err))
					}
				}
			})
		}
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		err := s.Shutdown(ctx)
		cancel()
		wg.Wait()
		if err != nil {
			t.Fatalf("round %d: Shutdown(2s timeout) = %v, want <nil>", round, err)
		}
		if e := badErr.Load(); e != nil {
			t.Fatalf("round %d: Submit returned %v, want nil, ErrClosed or ErrQueueFull", round, e)
		}
		if a, h := accepted.Load(), handledN.Load(); a != h {
			t.Fatalf("round %d: %d jobs accepted but %d handled, want equal (Shutdown must wait for every accepted job)", round, a, h)
		}
		if err := s.Submit(-1); !errors.Is(err, ErrClosed) {
			t.Fatalf("round %d: Submit after Shutdown = %v, want %v", round, err, ErrClosed)
		}
	}
}

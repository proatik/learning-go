package pool

import (
	"fmt"
	"sync/atomic"
	"testing"
	"testing/synctest"
	"time"
)

// jobsChan returns a closed channel holding jobs 1..n.
func jobsChan(n int) <-chan int {
	ch := make(chan int, n)
	for i := 1; i <= n; i++ {
		ch <- i
	}
	close(ch)
	return ch
}

// collect reads out until it is closed. limit bounds the wait: real time
// outside synctest, fake time inside a bubble.
func collect(t *testing.T, name string, out <-chan Result, limit time.Duration) []Result {
	t.Helper()
	if out == nil {
		t.Fatalf("%s returned a nil channel", name)
	}
	var got []Result
	timeout := time.After(limit)
	for {
		select {
		case r, ok := <-out:
			if !ok {
				return got
			}
			got = append(got, r)
		case <-timeout:
			t.Fatalf("%s: results channel not closed within %v (received %d results)", name, limit, len(got))
		}
	}
}

// tracker records how many work calls run at the same time.
type tracker struct {
	running, peak atomic.Int64
}

// work sleeps for d (fake time in a bubble) and returns j*10.
func (tr *tracker) work(d time.Duration) func(int) int {
	return func(j int) int {
		now := tr.running.Add(1)
		for {
			p := tr.peak.Load()
			if now <= p || tr.peak.CompareAndSwap(p, now) {
				break
			}
		}
		time.Sleep(d)
		tr.running.Add(-1)
		return j * 10
	}
}

func TestRunResults(t *testing.T) {
	for _, tt := range []struct{ workers, jobs int }{{4, 20}, {1, 5}, {8, 3}, {3, 0}} {
		name := fmt.Sprintf("Run(%d, %d jobs)", tt.workers, tt.jobs)
		got := collect(t, name, Run(tt.workers, jobsChan(tt.jobs), func(j int) int { return j * 10 }), 2*time.Second)
		if len(got) != tt.jobs {
			t.Errorf("%s returned %d results, want %d", name, len(got), tt.jobs)
			continue
		}
		seen := make(map[int]bool)
		for _, r := range got {
			if r.Value != r.Job*10 {
				t.Errorf("%s: result %+v, want Value = Job*10 = %d", name, r, r.Job*10)
			}
			if seen[r.Job] {
				t.Errorf("%s: job %d reported twice", name, r.Job)
			}
			seen[r.Job] = true
		}
		for j := 1; j <= tt.jobs; j++ {
			if !seen[j] {
				t.Errorf("%s: no result for job %d", name, j)
			}
		}
	}
}

func TestRunConcurrency(t *testing.T) {
	tests := []struct {
		workers, jobs int
		wantPeak      int64
		wantElapsed   time.Duration
	}{
		{3, 10, 3, 4 * time.Second}, // 3+3+3+1
		{1, 5, 1, 5 * time.Second},
		{5, 5, 5, 1 * time.Second},
		{10, 4, 4, 1 * time.Second}, // more workers than jobs
		{0, 3, 1, 3 * time.Second},  // n < 1 means 1 worker
		{-2, 2, 1, 2 * time.Second},
	}
	for _, tt := range tests {
		synctest.Test(t, func(t *testing.T) {
			name := fmt.Sprintf("Run(%d, %d jobs of 1s)", tt.workers, tt.jobs)
			var tr tracker
			start := time.Now()
			got := collect(t, name, Run(tt.workers, jobsChan(tt.jobs), tr.work(time.Second)), time.Hour)
			elapsed := time.Since(start)
			if len(got) != tt.jobs {
				t.Errorf("%s returned %d results, want %d", name, len(got), tt.jobs)
			}
			if p := tr.peak.Load(); p != tt.wantPeak {
				t.Errorf("%s max concurrent jobs = %d, want %d", name, p, tt.wantPeak)
			}
			if elapsed != tt.wantElapsed {
				t.Errorf("%s took %v, want %v", name, elapsed, tt.wantElapsed)
			}
		})
	}
}

// TestRunWaitsForOpenJobs: workers keep waiting while jobs is open, and the
// results channel closes only after jobs is closed.
func TestRunWaitsForOpenJobs(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		jobs := make(chan int)
		out := Run(2, jobs, func(j int) int { return -j })
		if out == nil {
			t.Fatalf("Run returned a nil channel")
		}
		for _, j := range []int{7, 8} {
			select {
			case jobs <- j:
			case <-time.After(time.Hour):
				t.Fatalf("Run: no worker took job %d", j)
			}
			select {
			case r := <-out:
				if r != (Result{j, -j}) {
					t.Errorf("Run: result %+v, want %+v", r, Result{j, -j})
				}
			case <-time.After(time.Hour):
				t.Fatalf("Run: no result for job %d", j)
			}
		}
		synctest.Wait()
		select {
		case r, ok := <-out:
			t.Fatalf("Run with jobs still open: received (%+v, %v), want to keep waiting", r, ok)
		default:
		}
		close(jobs)
		if got := collect(t, "Run", out, time.Hour); len(got) != 0 {
			t.Errorf("Run after close(jobs) = %v, want no more results", got)
		}
	})
}

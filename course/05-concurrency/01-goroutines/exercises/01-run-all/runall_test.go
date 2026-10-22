package runall_test

import (
	"slices"
	"sync/atomic"
	"testing"
	"testing/synctest"
	"time"

	runall "learngo/course/05-concurrency/01-goroutines/exercises/01-run-all"
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

// sleepers returns one task per duration; task i sleeps, then counts itself
// in runs[i].
func sleepers(durs []time.Duration, runs []atomic.Int32) []func() {
	tasks := make([]func(), len(durs))
	for i, d := range durs {
		tasks[i] = func() {
			time.Sleep(d)
			runs[i].Add(1)
		}
	}
	return tasks
}

func TestRunsConcurrently(t *testing.T) {
	tests := []struct {
		name string
		durs []time.Duration
		want time.Duration
	}{
		{"3 tasks of 1s", []time.Duration{time.Second, time.Second, time.Second}, time.Second},
		{"tasks of 1s, 3s, 2s", []time.Duration{time.Second, 3 * time.Second, 2 * time.Second}, 3 * time.Second},
		{"1 task of 5s", []time.Duration{5 * time.Second}, 5 * time.Second},
		{"10 tasks of 1s", slices.Repeat([]time.Duration{time.Second}, 10), time.Second},
	}
	for _, tt := range tests {
		synctest.Test(t, func(t *testing.T) {
			runs := make([]atomic.Int32, len(tt.durs))
			start := time.Now()
			runall.RunAll(sleepers(tt.durs, runs))
			elapsed := time.Since(start)
			for i := range runs {
				if n := runs[i].Load(); n != 1 {
					t.Errorf("RunAll(%s): task %d ran %d times when RunAll returned, want 1", tt.name, i, n)
				}
			}
			if elapsed != tt.want {
				t.Errorf("RunAll(%s) took %v, want %v", tt.name, elapsed, tt.want)
			}
		})
	}
}

func TestEmpty(t *testing.T) {
	within(t, "RunAll(nil)", func() { runall.RunAll(nil) })
	within(t, "RunAll([]func(){})", func() { runall.RunAll([]func(){}) })
}

func TestManyTasks(t *testing.T) {
	const n = 1000
	runs := make([]atomic.Int32, n)
	tasks := make([]func(), n)
	for i := range tasks {
		tasks[i] = func() { runs[i].Add(1) }
	}
	within(t, "RunAll(1000 tasks)", func() { runall.RunAll(tasks) })
	missing := 0
	for i := range runs {
		if runs[i].Load() != 1 {
			missing++
		}
	}
	if missing != 0 {
		t.Errorf("RunAll(1000 tasks): tasks not run exactly once = %d, want 0", missing)
	}
}

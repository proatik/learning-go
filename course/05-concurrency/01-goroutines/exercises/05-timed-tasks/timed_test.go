package timed_test

import (
	"slices"
	"testing"
	"testing/synctest"
	"time"

	timed "learngo/course/05-concurrency/01-goroutines/exercises/05-timed-tasks"
)

const s = time.Second

func sleeper(d time.Duration) func() {
	return func() { time.Sleep(d) }
}

func TestTimed(t *testing.T) {
	tests := []struct {
		durs  []time.Duration
		total time.Duration
	}{
		{[]time.Duration{1 * s, 3 * s, 2 * s}, 3 * s},
		{[]time.Duration{s, s, s, s, s}, s},
		{[]time.Duration{2 * s, 0}, 2 * s},
		{[]time.Duration{0, 0, 0}, 0},
		{[]time.Duration{10 * s}, 10 * s},
		{[]time.Duration{5 * s, 4 * s, 3 * s, 2 * s, 1 * s, 500 * time.Millisecond}, 5 * s},
		{[]time.Duration{90 * time.Minute, time.Hour}, 90 * time.Minute},
	}
	for _, tt := range tests {
		synctest.Test(t, func(t *testing.T) {
			tasks := make([]func(), len(tt.durs))
			for i, d := range tt.durs {
				tasks[i] = sleeper(d)
			}
			before := time.Now()
			total, each := timed.Timed(tasks)
			outside := time.Since(before)
			if !slices.Equal(each, tt.durs) {
				t.Errorf("Timed(sleeps %v) each = %v, want %v", tt.durs, each, tt.durs)
			}
			if total != tt.total {
				t.Errorf("Timed(sleeps %v) total = %v, want %v", tt.durs, total, tt.total)
			}
			if outside != tt.total {
				t.Errorf("Timed(sleeps %v) returned after %v, want %v", tt.durs, outside, tt.total)
			}
		})
	}
}

func TestEmpty(t *testing.T) {
	for _, tasks := range [][]func(){nil, {}} {
		synctest.Test(t, func(t *testing.T) {
			total, each := timed.Timed(tasks)
			if total != 0 || len(each) != 0 {
				t.Errorf("Timed(%d tasks) = (%v, %v), want (0s, [])", len(tasks), total, each)
			}
		})
	}
}

func TestStaggeredStart(t *testing.T) {
	// Task 1 waits until 1s after the batch started, then works for 3s. It
	// takes 4s only if it starts together with task 0 (run one after the
	// other, it would start at 2s and take 3s).
	synctest.Test(t, func(t *testing.T) {
		start := time.Now()
		tasks := []func(){
			func() { time.Sleep(2 * s) },
			func() { time.Sleep(time.Until(start.Add(1 * s))); time.Sleep(3 * s) },
		}
		total, each := timed.Timed(tasks)
		want := []time.Duration{2 * s, 4 * s}
		if !slices.Equal(each, want) {
			t.Errorf("Timed(sleep 2s; wait until 1s then sleep 3s) each = %v, want %v", each, want)
		}
		if total != 4*s {
			t.Errorf("Timed(sleep 2s; wait until 1s then sleep 3s) total = %v, want 4s", total)
		}
	})
}

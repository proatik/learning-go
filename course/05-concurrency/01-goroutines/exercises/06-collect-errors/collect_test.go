package collect_test

import (
	"errors"
	"fmt"
	"sync/atomic"
	"testing"
	"testing/synctest"
	"time"

	collect "learngo/course/05-concurrency/01-goroutines/exercises/06-collect-errors"
)

var (
	errA = errors.New("a failed")
	errB = errors.New("b failed")
	errC = errors.New("c failed")
)

// step describes one task: how long it sleeps and what it returns.
type step struct {
	sleep time.Duration
	err   error
}

func (s step) String() string {
	if s.err == nil {
		return fmt.Sprintf("ok after %v", s.sleep)
	}
	return fmt.Sprintf("%q after %v", s.err, s.sleep)
}

func TestRunAll(t *testing.T) {
	const s = time.Second
	tests := []struct {
		steps   []step
		want    string // "" means a nil error
		elapsed time.Duration
	}{
		{[]step{{s, nil}, {s, nil}, {s, nil}}, "", s},
		{[]step{{s, nil}, {s, errA}, {s, nil}, {s, errB}}, "a failed\nb failed", s},
		{[]step{{3 * s, errB}, {s, errA}}, "b failed\na failed", 3 * s},
		{[]step{{s, errC}, {2 * s, errB}, {3 * s, errA}}, "c failed\nb failed\na failed", 3 * s},
		{[]step{{2 * s, errA}}, "a failed", 2 * s},
		{[]step{{s, errA}, {5 * s, nil}}, "a failed", 5 * s},
		{[]step{{s, nil}, {0, errA}, {0, errA}}, "a failed\na failed", s},
	}
	for _, tt := range tests {
		synctest.Test(t, func(t *testing.T) {
			var ran atomic.Int32
			tasks := make([]func() error, len(tt.steps))
			for i, st := range tt.steps {
				tasks[i] = func() error {
					time.Sleep(st.sleep)
					ran.Add(1)
					return st.err
				}
			}
			start := time.Now()
			err := collect.RunAll(tasks)
			elapsed := time.Since(start)

			if tt.want == "" {
				if err != nil {
					t.Errorf("RunAll(%v) = %q, want nil", tt.steps, err)
				}
			} else if err == nil {
				t.Errorf("RunAll(%v) = nil, want %q", tt.steps, tt.want)
			} else if err.Error() != tt.want {
				t.Errorf("RunAll(%v) = %q, want %q", tt.steps, err, tt.want)
			}
			for _, st := range tt.steps {
				if st.err != nil && !errors.Is(err, st.err) {
					t.Errorf("RunAll(%v): errors.Is(err, %q) = false, want true", tt.steps, st.err)
				}
			}
			if r := int(ran.Load()); r != len(tt.steps) {
				t.Errorf("RunAll(%v): tasks finished = %d, want %d", tt.steps, r, len(tt.steps))
			}
			if elapsed != tt.elapsed {
				t.Errorf("RunAll(%v) took %v, want %v", tt.steps, elapsed, tt.elapsed)
			}
		})
	}
}

func TestEmpty(t *testing.T) {
	for _, tasks := range [][]func() error{nil, {}} {
		synctest.Test(t, func(t *testing.T) {
			if err := collect.RunAll(tasks); err != nil {
				t.Errorf("RunAll(%d tasks) = %q, want nil", len(tasks), err)
			}
		})
	}
}

func TestManyFailures(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		const n = 500
		tasks := make([]func() error, n)
		want := ""
		for i := range tasks {
			if i%3 == 0 {
				msg := fmt.Sprintf("task %d", i)
				tasks[i] = func() error { time.Sleep(time.Duration(n-i) * time.Millisecond); return errors.New(msg) }
				if want != "" {
					want += "\n"
				}
				want += msg
			} else {
				tasks[i] = func() error { return nil }
			}
		}
		err := collect.RunAll(tasks)
		if err == nil {
			t.Fatalf("RunAll(500 tasks, every 3rd fails) = nil, want %d joined errors", (n+2)/3)
		}
		if err.Error() != want {
			t.Errorf("RunAll(500 tasks, every 3rd fails) = %d bytes of message, want %d bytes listing tasks 0, 3, 6, … in order",
				len(err.Error()), len(want))
		}
	})
}

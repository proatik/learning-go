package dag

import (
	"context"
	"errors"
	"fmt"
	"maps"
	"sync"
	"testing"
	"testing/synctest"
	"time"
)

var errBoom = errors.New("boom")

// log records when each job started and ended (fake time since start).
type log struct {
	mu         sync.Mutex
	start      time.Time
	began, end map[string]time.Duration
}

func newLog() *log {
	return &log{start: time.Now(), began: map[string]time.Duration{}, end: map[string]time.Duration{}}
}

// job returns a Job that sleeps d (or stops early if ctx ends) and then
// returns err.
func (l *log) job(name string, d time.Duration, err error, deps ...string) Job {
	return Job{Name: name, Deps: deps, Run: func(ctx context.Context) error {
		l.mu.Lock()
		l.began[name] = time.Since(l.start)
		l.mu.Unlock()
		defer func() {
			l.mu.Lock()
			l.end[name] = time.Since(l.start)
			l.mu.Unlock()
		}()
		select {
		case <-time.After(d):
			return err
		case <-ctx.Done():
			return ctx.Err()
		}
	}}
}

func (l *log) ran(name string) bool {
	l.mu.Lock()
	defer l.mu.Unlock()
	_, ok := l.began[name]
	return ok
}

func wantReport(t *testing.T, step string, got, want map[string]Status) {
	t.Helper()
	if got == nil {
		t.Fatalf("%s: Run() report = nil, want %v", step, want)
	}
	if !maps.Equal(got, want) {
		t.Fatalf("%s: Run() report = %v, want %v", step, got, want)
	}
}

func TestDiamond(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		l := newLog()
		jobs := []Job{
			l.job("d", time.Second, nil, "b", "c"),
			l.job("b", time.Second, nil, "a"),
			l.job("a", time.Second, nil),
			l.job("c", time.Second, nil, "a"),
		}
		report, err := Run(context.Background(), jobs)
		if err != nil {
			t.Fatalf("Run(diamond) error = %v, want <nil>", err)
		}
		wantReport(t, "diamond", report, map[string]Status{"a": Done, "b": Done, "c": Done, "d": Done})
		wantStart := map[string]time.Duration{"a": 0, "b": time.Second, "c": time.Second, "d": 2 * time.Second}
		for name, want := range wantStart {
			if got := l.began[name]; got != want {
				t.Errorf("diamond: job %q started at %v, want %v", name, got, want)
			}
		}
		if el := time.Since(l.start); el != 3*time.Second {
			t.Errorf("Run(diamond of 1s jobs) took %v, want 3s", el)
		}
	})
}

func TestIndependentJobsRunConcurrently(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		l := newLog()
		var jobs []Job
		want := map[string]Status{}
		for i := range 5 {
			name := fmt.Sprintf("j%d", i)
			jobs = append(jobs, l.job(name, time.Second, nil))
			want[name] = Done
		}
		report, err := Run(context.Background(), jobs)
		if err != nil {
			t.Fatalf("Run(5 independent jobs) error = %v, want <nil>", err)
		}
		wantReport(t, "5 independent jobs", report, want)
		if el := time.Since(l.start); el != time.Second {
			t.Errorf("Run(5 independent 1s jobs) took %v, want 1s", el)
		}
	})
}

func TestFailureSkipsDependents(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		l := newLog()
		jobs := []Job{
			l.job("a", time.Second, nil),
			l.job("b", time.Second, errBoom, "a"),
			l.job("c", time.Second, nil, "b"),
			l.job("c2", time.Second, nil, "c", "a"),
			l.job("d", 3*time.Second, nil, "a"),
			l.job("e", time.Second, nil),
		}
		report, err := Run(context.Background(), jobs)
		if !errors.Is(err, errBoom) || err.Error() != `job "b": boom` {
			t.Fatalf(`Run() error = %v, want job "b": boom (wrapping errBoom)`, err)
		}
		wantReport(t, "b fails", report, map[string]Status{
			"a": Done, "b": Failed, "c": Skipped, "c2": Skipped, "d": Done, "e": Done,
		})
		for _, name := range []string{"c", "c2"} {
			if l.ran(name) {
				t.Errorf("job %q ran, want it skipped (it depends on the failed job)", name)
			}
		}
		if el := time.Since(l.start); el != 4*time.Second {
			t.Errorf("Run() took %v, want 4s (unrelated job d still runs to the end)", el)
		}
	})
}

func TestFirstFailureWins(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		l := newLog()
		errLate := errors.New("late")
		jobs := []Job{
			l.job("late", 2*time.Second, errLate),
			l.job("early", time.Second, errBoom),
		}
		report, err := Run(context.Background(), jobs)
		if !errors.Is(err, errBoom) {
			t.Fatalf("Run() error = %v, want the first failure in time: job \"early\": boom", err)
		}
		wantReport(t, "two failures", report, map[string]Status{"late": Failed, "early": Failed})
	})
}

func TestValidation(t *testing.T) {
	tests := []struct {
		name string
		jobs func(l *log) []Job
		want error
		msg  string
	}{
		{"two-job cycle", func(l *log) []Job {
			return []Job{l.job("a", 0, nil, "b"), l.job("b", 0, nil, "a")}
		}, ErrCycle, ""},
		{"self dependency", func(l *log) []Job {
			return []Job{l.job("ok", 0, nil), l.job("a", 0, nil, "a")}
		}, ErrCycle, ""},
		{"cycle behind a valid start", func(l *log) []Job {
			return []Job{l.job("a", 0, nil), l.job("b", 0, nil, "a", "d"), l.job("c", 0, nil, "b"), l.job("d", 0, nil, "c")}
		}, ErrCycle, ""},
		{"unknown dependency", func(l *log) []Job {
			return []Job{l.job("x", 0, nil), l.job("a", 0, nil, "x", "zz")}
		}, ErrUnknownDep, `dag: unknown dependency: "a" needs "zz"`},
		{"duplicate name", func(l *log) []Job {
			return []Job{l.job("a", 0, nil), l.job("b", 0, nil), l.job("a", 0, nil)}
		}, ErrDuplicate, `dag: duplicate job: "a"`},
	}
	for _, tt := range tests {
		synctest.Test(t, func(t *testing.T) {
			l := newLog()
			report, err := Run(context.Background(), tt.jobs(l))
			if !errors.Is(err, tt.want) {
				t.Fatalf("%s: Run() error = %v, want %v", tt.name, err, tt.want)
			}
			if tt.msg != "" && err.Error() != tt.msg {
				t.Errorf("%s: Run() error = %q, want %q", tt.name, err.Error(), tt.msg)
			}
			if report != nil {
				t.Errorf("%s: Run() report = %v, want nil", tt.name, report)
			}
			if n := len(l.began); n != 0 {
				t.Errorf("%s: %d jobs ran, want 0 (validate before running anything)", tt.name, n)
			}
		})
	}
}

func TestNoJobs(t *testing.T) {
	report, err := Run(context.Background(), nil)
	if err != nil {
		t.Fatalf("Run(nil) error = %v, want <nil>", err)
	}
	wantReport(t, "no jobs", report, map[string]Status{})
}

func TestCancelMidway(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		l := newLog()
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		time.AfterFunc(1500*time.Millisecond, cancel)
		jobs := []Job{
			l.job("a", time.Second, nil),
			l.job("b", time.Second, nil, "a"),
			l.job("c", time.Second, nil, "b"),
			l.job("z", time.Second, nil),
		}
		report, err := Run(ctx, jobs)
		if !errors.Is(err, context.Canceled) || err.Error() != `job "b": context canceled` {
			t.Fatalf(`Run() cancelled at 1.5s: error = %v, want job "b": context canceled`, err)
		}
		wantReport(t, "cancelled at 1.5s", report, map[string]Status{"a": Done, "b": Failed, "c": Skipped, "z": Done})
		if el := time.Since(l.start); el != 1500*time.Millisecond {
			t.Errorf("Run() cancelled at 1.5s returned after %v, want 1.5s", el)
		}
	})
}

func TestCancelledBeforeStart(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		l := newLog()
		ctx, cancel := context.WithCancel(context.Background())
		cancel()
		jobs := []Job{l.job("a", time.Second, nil), l.job("b", time.Second, nil, "a")}
		report, err := Run(ctx, jobs)
		if !errors.Is(err, context.Canceled) {
			t.Fatalf("Run(cancelled ctx) error = %v, want %v", err, context.Canceled)
		}
		wantReport(t, "cancelled ctx", report, map[string]Status{"a": Skipped, "b": Skipped})
		if len(l.began) != 0 {
			t.Errorf("Run(cancelled ctx): %d jobs ran, want 0", len(l.began))
		}
	})
}

func TestLongChain(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		l := newLog()
		var jobs []Job
		want := map[string]Status{}
		for i := range 100 {
			name := fmt.Sprintf("j%03d", i)
			var deps []string
			for _, back := range []int{1, 2, 7} {
				if i-back >= 0 {
					deps = append(deps, fmt.Sprintf("j%03d", i-back))
				}
			}
			jobs = append(jobs, l.job(name, time.Millisecond, nil, deps...))
			want[name] = Done
		}
		report, err := Run(context.Background(), jobs)
		if err != nil {
			t.Fatalf("Run(100-job chain) error = %v, want <nil>", err)
		}
		wantReport(t, "100-job chain", report, want)
		if el := time.Since(l.start); el != 100*time.Millisecond {
			t.Errorf("Run(100-job chain of 1ms jobs) took %v, want 100ms", el)
		}
	})
}

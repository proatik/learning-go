package first

import (
	"context"
	"errors"
	"fmt"
	"sync/atomic"
	"testing"
	"testing/synctest"
	"time"
)

// replica answers val (or fails with err) after delay, unless its ctx is
// canceled first. It records what happened.
type replica struct {
	name     string
	delay    time.Duration
	err      error
	started  atomic.Bool
	canceled atomic.Bool // returned because its ctx was canceled
}

func (r *replica) fetch(ctx context.Context) (string, error) {
	r.started.Store(true)
	select {
	case <-time.After(r.delay):
		if r.err != nil {
			return "", r.err
		}
		return r.name, nil
	case <-ctx.Done():
		r.canceled.Store(true)
		return "", ctx.Err()
	}
}

func funcs(rs []*replica) []func(context.Context) (string, error) {
	out := make([]func(context.Context) (string, error), len(rs))
	for i, r := range rs {
		out[i] = r.fetch
	}
	return out
}

func describe(rs []*replica) string {
	s := ""
	for _, r := range rs {
		res := r.name
		if r.err != nil {
			res = "fail"
		}
		s += fmt.Sprintf("[%s %v %s]", r.name, r.delay, res)
	}
	return s
}

func TestFirst(t *testing.T) {
	errA, errB := errors.New("a failed"), errors.New("b failed")
	tests := []struct {
		name         string
		replicas     func() []*replica
		want         string
		wantAfter    time.Duration
		wantCanceled []string
	}{
		{
			name: "fastest wins",
			replicas: func() []*replica {
				return []*replica{{name: "A", delay: 3 * time.Second}, {name: "B", delay: time.Second}, {name: "C", delay: 2 * time.Second}}
			},
			want: "B", wantAfter: time.Second, wantCanceled: []string{"A", "C"},
		},
		{
			name: "failure is skipped",
			replicas: func() []*replica {
				return []*replica{{name: "A", delay: 100 * time.Millisecond, err: errA}, {name: "B", delay: 2 * time.Second}, {name: "C", delay: 3 * time.Second}}
			},
			want: "B", wantAfter: 2 * time.Second, wantCanceled: []string{"C"},
		},
		{
			name: "single replica",
			replicas: func() []*replica {
				return []*replica{{name: "solo", delay: 5 * time.Second}}
			},
			want: "solo", wantAfter: 5 * time.Second,
		},
		{
			name: "only the last succeeds",
			replicas: func() []*replica {
				return []*replica{{name: "A", delay: time.Second, err: errA}, {name: "B", delay: 2 * time.Second, err: errB}, {name: "C", delay: 4 * time.Second}}
			},
			want: "C", wantAfter: 4 * time.Second,
		},
	}
	for _, tt := range tests {
		synctest.Test(t, func(t *testing.T) {
			rs := tt.replicas()
			start := time.Now()
			got, err := First(t.Context(), funcs(rs)...)
			elapsed := time.Since(start)
			if got != tt.want || err != nil {
				t.Fatalf("First(%s) = (%q, %v), want (%q, <nil>)", describe(rs), got, err, tt.want)
			}
			if elapsed != tt.wantAfter {
				t.Errorf("First(%s) returned after %v, want %v (run the fetchers concurrently)", describe(rs), elapsed, tt.wantAfter)
			}
			synctest.Wait()
			for _, r := range rs {
				if !r.started.Load() {
					t.Errorf("First(%s): fetcher %s never started", describe(rs), r.name)
				}
			}
			for _, name := range tt.wantCanceled {
				for _, r := range rs {
					if r.name == name && !r.canceled.Load() {
						t.Errorf("First(%s): loser %s did not see its ctx canceled, want canceled", describe(rs), name)
					}
				}
			}
		})
	}
}

func TestAllFail(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		errA, errB, errC := errors.New("a failed"), errors.New("b failed"), errors.New("c failed")
		rs := []*replica{{name: "A", delay: 3 * time.Second, err: errA}, {name: "B", delay: time.Second, err: errB}, {name: "C", delay: 2 * time.Second, err: errC}}
		start := time.Now()
		got, err := First(t.Context(), funcs(rs)...)
		elapsed := time.Since(start)
		if got != "" || err == nil {
			t.Fatalf("First(%s) = (%q, %v), want (\"\", joined errors)", describe(rs), got, err)
		}
		for _, e := range []error{errA, errB, errC} {
			if !errors.Is(err, e) {
				t.Errorf("First(%s) error = %v, want it to wrap %v", describe(rs), err, e)
			}
		}
		if want := "a failed\nb failed\nc failed"; err.Error() != want {
			t.Errorf("First(%s) error = %q, want %q (argument order)", describe(rs), err.Error(), want)
		}
		if elapsed != 3*time.Second {
			t.Errorf("First(%s) returned after %v, want 3s", describe(rs), elapsed)
		}
	})
}

func TestNoFetchers(t *testing.T) {
	got, err := First(t.Context())
	if got != "" || !errors.Is(err, ErrNoFetchers) {
		t.Errorf("First() = (%q, %v), want (\"\", %v)", got, err, ErrNoFetchers)
	}
}

func TestParentTimeout(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		ctx, cancel := context.WithTimeout(t.Context(), time.Second)
		defer cancel()
		rs := []*replica{{name: "A", delay: 5 * time.Second}, {name: "B", delay: 6 * time.Second}}
		start := time.Now()
		got, err := First(ctx, funcs(rs)...)
		elapsed := time.Since(start)
		if got != "" || !errors.Is(err, context.DeadlineExceeded) {
			t.Errorf("First(timeout 1s, %s) = (%q, %v), want (\"\", %v)", describe(rs), got, err, context.DeadlineExceeded)
		}
		if elapsed != time.Second {
			t.Errorf("First(timeout 1s, %s) returned after %v, want 1s", describe(rs), elapsed)
		}
	})
}

func TestChildCanceledAfterSuccess(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		var seen context.Context
		f := func(ctx context.Context) (string, error) {
			seen = ctx
			return "ok", nil
		}
		got, err := First(t.Context(), f)
		if got != "ok" || err != nil {
			t.Fatalf("First(instant ok) = (%q, %v), want (\"ok\", <nil>)", got, err)
		}
		synctest.Wait()
		if seen == nil {
			t.Fatalf("First(instant ok): fetcher never called")
		}
		if seen.Err() == nil {
			t.Errorf("First(instant ok): fetcher's ctx Err() after return = <nil>, want %v (cancel the child)", context.Canceled)
		}
		if t.Context().Err() != nil {
			t.Errorf("First canceled the caller's ctx; cancel only the child")
		}
	})
}

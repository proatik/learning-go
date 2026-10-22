package fetch

import (
	"context"
	"errors"
	"testing"
	"testing/synctest"
	"time"
)

var ErrNotFound = errors.New("not found")

type key struct{}

// service answers "ok" after delay, or returns ctx.Err() if ctx ends first.
// It records the context it was called with.
type service struct {
	delay time.Duration
	calls int
	ctx   context.Context
}

func (s *service) call(ctx context.Context) (string, error) {
	s.calls++
	s.ctx = ctx
	select {
	case <-time.After(s.delay):
		return "ok", nil
	case <-ctx.Done():
		return "", ctx.Err()
	}
}

func TestFetch(t *testing.T) {
	tests := []struct {
		name      string
		parentTO  time.Duration // > 0: parent has this timeout
		cancelAt  time.Duration // > 0: parent canceled at this time
		d         time.Duration
		delay     time.Duration
		want      string
		wantErr   error
		wantMsg   string
		wantAfter time.Duration
		wantDL    time.Duration // deadline f should see, relative to start; 0 = none
	}{
		{name: "d=1s, answers after 100ms", d: time.Second, delay: 100 * time.Millisecond, want: "ok", wantAfter: 100 * time.Millisecond, wantDL: time.Second},
		{name: "d=1s, takes 5s", d: time.Second, delay: 5 * time.Second, wantErr: context.DeadlineExceeded, wantMsg: "fetch: context deadline exceeded", wantAfter: time.Second, wantDL: time.Second},
		{name: "d=1s, parent times out at 500ms", parentTO: 500 * time.Millisecond, d: time.Second, delay: 5 * time.Second, wantErr: context.DeadlineExceeded, wantAfter: 500 * time.Millisecond, wantDL: 500 * time.Millisecond},
		{name: "d=1s, parent timeout 3s", parentTO: 3 * time.Second, d: time.Second, delay: 2 * time.Second, wantErr: context.DeadlineExceeded, wantAfter: time.Second, wantDL: time.Second},
		{name: "d=1s, parent canceled at 200ms", cancelAt: 200 * time.Millisecond, d: time.Second, delay: 5 * time.Second, wantErr: context.Canceled, wantMsg: "fetch: context canceled", wantAfter: 200 * time.Millisecond, wantDL: time.Second},
		{name: "d=0, answers after 1h", d: 0, delay: time.Hour, want: "ok", wantAfter: time.Hour},
		{name: "d=-1s, answers after 2s", d: -time.Second, delay: 2 * time.Second, want: "ok", wantAfter: 2 * time.Second},
		{name: "d=0, parent times out at 3s", parentTO: 3 * time.Second, d: 0, delay: time.Hour, wantErr: context.DeadlineExceeded, wantAfter: 3 * time.Second, wantDL: 3 * time.Second},
	}
	for _, tt := range tests {
		synctest.Test(t, func(t *testing.T) {
			start := time.Now()
			parent, cancel := context.WithCancel(context.WithValue(t.Context(), key{}, "v"))
			defer cancel()
			if tt.parentTO > 0 {
				parent, cancel = context.WithTimeout(parent, tt.parentTO)
				defer cancel()
			}
			if tt.cancelAt > 0 {
				time.AfterFunc(tt.cancelAt, cancel)
			}
			svc := &service{delay: tt.delay}
			got, err := Fetch(parent, tt.d, svc.call)
			elapsed := time.Since(start)

			if svc.calls != 1 {
				t.Fatalf("Fetch(%s) called f %d times, want 1", tt.name, svc.calls)
			}
			if got != tt.want || !errors.Is(err, tt.wantErr) {
				t.Errorf("Fetch(%s) = (%q, %v), want (%q, %v)", tt.name, got, err, tt.want, tt.wantErr)
			}
			if tt.wantMsg != "" && (err == nil || err.Error() != tt.wantMsg) {
				t.Errorf("Fetch(%s) error = %v, want %q", tt.name, err, tt.wantMsg)
			}
			if elapsed != tt.wantAfter {
				t.Errorf("Fetch(%s) returned after %v, want %v", tt.name, elapsed, tt.wantAfter)
			}
			dl, ok := svc.ctx.Deadline()
			switch {
			case tt.wantDL == 0 && ok:
				t.Errorf("Fetch(%s): f's ctx deadline = start+%v, want no deadline", tt.name, dl.Sub(start))
			case tt.wantDL != 0 && (!ok || dl.Sub(start) != tt.wantDL):
				t.Errorf("Fetch(%s): f's ctx deadline = start+%v (set %v), want start+%v", tt.name, dl.Sub(start), ok, tt.wantDL)
			}
			if v := svc.ctx.Value(key{}); v != "v" {
				t.Errorf("Fetch(%s): f's ctx Value(key) = %v, want \"v\" (derive from the ctx you were given)", tt.name, v)
			}
			if svc.ctx.Err() == nil {
				t.Errorf("Fetch(%s): f's ctx Err() after Fetch returned = <nil>, want %v (call cancel)", tt.name, context.Canceled)
			}
		})
	}
}

func TestFetchError(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		var fctx context.Context
		f := func(ctx context.Context) (string, error) {
			fctx = ctx
			return "partial", ErrNotFound
		}
		got, err := Fetch(t.Context(), time.Second, f)
		if fctx == nil {
			t.Fatalf("Fetch(f returns ErrNotFound) never called f")
		}
		if got != "" || !errors.Is(err, ErrNotFound) {
			t.Errorf("Fetch(f returns ErrNotFound) = (%q, %v), want (\"\", %v)", got, err, ErrNotFound)
		}
		if err == nil || err.Error() != "fetch: not found" {
			t.Errorf("Fetch(f returns ErrNotFound) error = %v, want %q", err, "fetch: not found")
		}
		if fctx.Err() != context.Canceled {
			t.Errorf("Fetch(f returns ErrNotFound): f's ctx Err() afterwards = %v, want %v (cancel on every path)", fctx.Err(), context.Canceled)
		}
	})
}

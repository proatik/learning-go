package bridge

import (
	"context"
	"errors"
	"sync"
	"sync/atomic"
	"testing"
	"testing/synctest"
	"time"
)

var ErrBroken = errors.New("connection broken")

// fakeReader delivers messages from msgs; Close makes a blocked Read return
// ErrClosed.
type fakeReader struct {
	msgs   chan string
	fail   error // if set, Read returns it at once
	closed chan struct{}
	once   sync.Once
	closes atomic.Int32
	reads  atomic.Int32
}

func newReader() *fakeReader {
	return &fakeReader{msgs: make(chan string, 1), closed: make(chan struct{})}
}

func (r *fakeReader) Read() (string, error) {
	r.reads.Add(1)
	if r.fail != nil {
		return "", r.fail
	}
	select {
	case m := <-r.msgs:
		return m, nil
	case <-r.closed:
		return "", ErrClosed
	}
}

func (r *fakeReader) Close() error {
	r.closes.Add(1)
	r.once.Do(func() { close(r.closed) })
	return nil
}

func TestReadCtx(t *testing.T) {
	tests := []struct {
		name      string
		msgAt     time.Duration // > 0: message "hi" arrives then; < 0: already waiting
		cancelAt  time.Duration // > 0: cancel then; < 0: canceled before the call
		timeout   time.Duration
		fail      error
		want      string
		wantErr   error
		wantAfter time.Duration
		wantClose int32
	}{
		{name: "message at 1s", msgAt: time.Second, want: "hi", wantAfter: time.Second},
		{name: "message already waiting", msgAt: -1, want: "hi"},
		{name: "no message, cancel at 2s", cancelAt: 2 * time.Second, wantErr: context.Canceled, wantAfter: 2 * time.Second, wantClose: 1},
		{name: "message at 3s, cancel at 2s", msgAt: 3 * time.Second, cancelAt: 2 * time.Second, wantErr: context.Canceled, wantAfter: 2 * time.Second, wantClose: 1},
		{name: "no message, timeout 500ms", timeout: 500 * time.Millisecond, wantErr: context.DeadlineExceeded, wantAfter: 500 * time.Millisecond, wantClose: 1},
		{name: "already canceled, message waiting", msgAt: -1, cancelAt: -1, wantErr: context.Canceled, wantClose: 1},
		{name: "Read fails with ErrBroken", fail: ErrBroken, wantErr: ErrBroken},
	}
	for _, tt := range tests {
		synctest.Test(t, func(t *testing.T) {
			ctx, cancel := context.WithCancel(t.Context())
			defer cancel()
			if tt.timeout > 0 {
				ctx, cancel = context.WithTimeout(ctx, tt.timeout)
				defer cancel()
			}
			r := newReader()
			r.fail = tt.fail
			switch {
			case tt.msgAt < 0:
				r.msgs <- "hi"
			case tt.msgAt > 0:
				time.AfterFunc(tt.msgAt, func() {
					select {
					case r.msgs <- "hi":
					default:
					}
				})
			}
			switch {
			case tt.cancelAt < 0:
				cancel()
			case tt.cancelAt > 0:
				time.AfterFunc(tt.cancelAt, cancel)
			}
			start := time.Now()
			got, err := ReadCtx(ctx, r)
			elapsed := time.Since(start)
			if n := r.reads.Load(); n != 1 {
				t.Fatalf("ReadCtx(%s) called Read %d times, want 1", tt.name, n)
			}
			if got != tt.want || !errors.Is(err, tt.wantErr) {
				t.Errorf("ReadCtx(%s) = (%q, %v), want (%q, %v)", tt.name, got, err, tt.want, tt.wantErr)
			}
			if elapsed != tt.wantAfter {
				t.Errorf("ReadCtx(%s) returned after %v, want %v", tt.name, elapsed, tt.wantAfter)
			}
			// Cancel afterwards: a stopped AfterFunc must not close r now.
			cancel()
			synctest.Wait()
			if n := r.closes.Load(); n != tt.wantClose {
				t.Errorf("ReadCtx(%s): Close called %d times (checked after a later cancel), want %d", tt.name, n, tt.wantClose)
			}
		})
	}
}

type key struct{}

func TestMerge(t *testing.T) {
	ErrX := errors.New("cause X")
	ErrY := errors.New("cause Y")
	tests := []struct {
		name      string
		act       func(cancelA, cancelB context.CancelCauseFunc, cancelM context.CancelFunc)
		wantErr   error
		wantCause error
	}{
		{"nothing canceled", func(_, _ context.CancelCauseFunc, _ context.CancelFunc) {}, nil, nil},
		{"a canceled with X", func(a, _ context.CancelCauseFunc, _ context.CancelFunc) { a(ErrX) }, context.Canceled, ErrX},
		{"b canceled with Y", func(_, b context.CancelCauseFunc, _ context.CancelFunc) { b(ErrY) }, context.Canceled, ErrY},
		{"b canceled with nil", func(_, b context.CancelCauseFunc, _ context.CancelFunc) { b(nil) }, context.Canceled, context.Canceled},
		{"merged cancel()", func(_, _ context.CancelCauseFunc, m context.CancelFunc) { m() }, context.Canceled, context.Canceled},
		{"merged cancel(), then b canceled with Y", func(_, b context.CancelCauseFunc, m context.CancelFunc) { m(); b(ErrY) }, context.Canceled, context.Canceled},
		{"b canceled with Y, then a with X", func(a, b context.CancelCauseFunc, _ context.CancelFunc) { b(ErrY); time.Sleep(time.Second); a(ErrX) }, context.Canceled, ErrY},
	}
	for _, tt := range tests {
		synctest.Test(t, func(t *testing.T) {
			a, cancelA := context.WithCancelCause(context.WithValue(t.Context(), key{}, "from a"))
			defer cancelA(nil)
			b, cancelB := context.WithCancelCause(t.Context())
			defer cancelB(nil)
			m, cancelM := Merge(a, b)
			if m == nil {
				t.Fatalf("Merge(a, b) returned a nil context")
			}
			defer cancelM()
			tt.act(cancelA, cancelB, cancelM)
			synctest.Wait()
			if err, cause := m.Err(), context.Cause(m); err != tt.wantErr || cause != tt.wantCause {
				t.Errorf("Merge(a, b), %s: (Err, Cause) = (%v, %v), want (%v, %v)", tt.name, err, cause, tt.wantErr, tt.wantCause)
			}
			if v := m.Value(key{}); v != "from a" {
				t.Errorf("Merge(a, b).Value(key) = %v, want \"from a\" (derive from a)", v)
			}
		})
	}
}

func TestMergeUntouchedParents(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		a, cancelA := context.WithCancel(t.Context())
		defer cancelA()
		b, cancelB := context.WithCancel(t.Context())
		defer cancelB()
		_, cancel := Merge(a, b)
		cancel()
		synctest.Wait()
		if a.Err() != nil || b.Err() != nil {
			t.Errorf("Merge(a, b) cancel(): (a.Err, b.Err) = (%v, %v), want (<nil>, <nil>)", a.Err(), b.Err())
		}
	})
}

func TestMergeTimeouts(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		a, cancelA := context.WithTimeout(t.Context(), time.Hour)
		defer cancelA()
		b, cancelB := context.WithTimeout(t.Context(), 2*time.Second)
		defer cancelB()
		start := time.Now()
		m, cancel := Merge(a, b)
		defer cancel()
		if dl, ok := m.Deadline(); !ok || dl.Sub(start) != time.Hour {
			t.Errorf("Merge(a 1h, b 2s).Deadline() = (start+%v, %v), want (start+1h0m0s, true) (a's deadline)", dl.Sub(start), ok)
		}
		time.Sleep(2*time.Second - time.Nanosecond)
		synctest.Wait()
		if err := m.Err(); err != nil {
			t.Fatalf("Merge(a 1h, b 2s) at 2s-1ns: Err = %v, want <nil>", err)
		}
		time.Sleep(time.Nanosecond)
		synctest.Wait()
		if err, cause := m.Err(), context.Cause(m); err != context.Canceled || cause != context.DeadlineExceeded {
			t.Errorf("Merge(a 1h, b 2s) at 2s: (Err, Cause) = (%v, %v), want (%v, %v)", err, cause, context.Canceled, context.DeadlineExceeded)
		}
	})
}

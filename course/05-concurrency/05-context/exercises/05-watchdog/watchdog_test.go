package watchdog

import (
	"context"
	"errors"
	"testing"
	"testing/synctest"
	"time"
)

// beat sends one heartbeat. Watch's goroutine must be waiting for it.
func beat(t *testing.T, beats chan struct{}) {
	t.Helper()
	synctest.Wait()
	select {
	case beats <- struct{}{}:
	default:
		t.Fatalf("Watch is not receiving from beats (start a goroutine that reads them)")
	}
}

// expect checks ctx.Err() and context.Cause(ctx) once every goroutine has
// caught up.
func expect(t *testing.T, step string, ctx context.Context, wantErr, wantCause error) {
	t.Helper()
	synctest.Wait()
	if err, cause := ctx.Err(), context.Cause(ctx); err != wantErr || cause != wantCause {
		t.Fatalf("%s: (Err, Cause) = (%v, %v), want (%v, %v)", step, err, cause, wantErr, wantCause)
	}
}

func TestSteadyBeats(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		beats := make(chan struct{})
		ctx, cancel := Watch(t.Context(), beats, time.Second)
		defer cancel()
		for range 20 {
			time.Sleep(500 * time.Millisecond)
			beat(t, beats)
		}
		expect(t, "a beat every 500ms for 10s (timeout 1s)", ctx, nil, nil)
		time.Sleep(999 * time.Millisecond)
		expect(t, "999ms after the last beat", ctx, nil, nil)
	})
}

func TestStall(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		beats := make(chan struct{})
		ctx, cancel := Watch(t.Context(), beats, time.Second)
		defer cancel()
		time.Sleep(500 * time.Millisecond)
		beat(t, beats)
		time.Sleep(700 * time.Millisecond)
		beat(t, beats)
		time.Sleep(999 * time.Millisecond)
		expect(t, "beats at 0.5s and 1.2s, now 2.199s", ctx, nil, nil)
		time.Sleep(time.Millisecond)
		expect(t, "beats at 0.5s and 1.2s, now 2.2s", ctx, context.Canceled, ErrStalled)
		if t.Context().Err() != nil {
			t.Fatalf("the parent was canceled too; Watch must cancel only its child")
		}
	})
}

func TestNoBeatAtAll(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		ctx, cancel := Watch(t.Context(), make(chan struct{}), 3*time.Second)
		defer cancel()
		time.Sleep(3*time.Second - time.Nanosecond)
		expect(t, "no beats, timeout 3s, now 3s-1ns", ctx, nil, nil)
		time.Sleep(time.Nanosecond)
		expect(t, "no beats, timeout 3s, now 3s", ctx, context.Canceled, ErrStalled)
	})
}

func TestClosed(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		beats := make(chan struct{})
		ctx, cancel := Watch(t.Context(), beats, time.Second)
		defer cancel()
		time.Sleep(300 * time.Millisecond)
		close(beats)
		expect(t, "beats closed at 300ms", ctx, context.Canceled, ErrStopped)
		time.Sleep(5 * time.Second)
		expect(t, "beats closed at 300ms, 5s later", ctx, context.Canceled, ErrStopped)
	})
}

func TestCancelFunc(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		beats := make(chan struct{})
		ctx, cancel := Watch(t.Context(), beats, time.Second)
		time.Sleep(500 * time.Millisecond)
		beat(t, beats)
		cancel()
		expect(t, "cancel() at 500ms", ctx, context.Canceled, context.Canceled)
		time.Sleep(5 * time.Second)
		expect(t, "cancel() at 500ms, 5s later (a stall must not change the cause)", ctx, context.Canceled, context.Canceled)
		cancel() // a second call is harmless
	})
}

func TestParentCanceledWithCause(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		ErrShutdown := errors.New("shutdown")
		parent, stop := context.WithCancelCause(t.Context())
		defer stop(nil)
		beats := make(chan struct{})
		ctx, cancel := Watch(parent, beats, time.Second)
		defer cancel()
		time.Sleep(400 * time.Millisecond)
		beat(t, beats)
		stop(ErrShutdown)
		expect(t, "parent canceled with ErrShutdown", ctx, context.Canceled, ErrShutdown)
	})
}

func TestParentTimeout(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		parent, stop := context.WithTimeout(t.Context(), 3*time.Second)
		defer stop()
		beats := make(chan struct{})
		ctx, cancel := Watch(parent, beats, time.Second)
		defer cancel()
		for range 5 {
			time.Sleep(500 * time.Millisecond)
			beat(t, beats)
		}
		expect(t, "parent timeout 3s, beats until 2.5s, now 2.5s", ctx, nil, nil)
		time.Sleep(500 * time.Millisecond)
		expect(t, "parent timeout 3s, now 3s", ctx, context.DeadlineExceeded, context.DeadlineExceeded)
	})
}

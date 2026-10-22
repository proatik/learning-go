package tee

import (
	"context"
	"slices"
	"testing"
	"testing/synctest"
	"time"
)

// closed returns a closed channel holding vals.
func closed[T any](vals ...T) <-chan T {
	ch := make(chan T, len(vals))
	for _, v := range vals {
		ch <- v
	}
	close(ch)
	return ch
}

func notNil[T any](t *testing.T, a, b <-chan T) {
	t.Helper()
	if a == nil || b == nil {
		t.Fatalf("Tee returned a nil channel (%v, %v)", a, b)
	}
}

// drainBoth reads a and b at the same time until both are closed, failing
// after 2s.
func drainBoth[T any](t *testing.T, a, b <-chan T) (gotA, gotB []T) {
	t.Helper()
	notNil(t, a, b)
	timeout := time.After(2 * time.Second)
	for a != nil || b != nil {
		select {
		case v, ok := <-a:
			if !ok {
				a = nil
				continue
			}
			gotA = append(gotA, v)
		case v, ok := <-b:
			if !ok {
				b = nil
				continue
			}
			gotB = append(gotB, v)
		case <-timeout:
			t.Fatalf("Tee: outputs not both closed within 2s (got %v and %v so far)", gotA, gotB)
		}
	}
	return gotA, gotB
}

func TestTee(t *testing.T) {
	tests := [][]int{
		{1, 2, 3, 4, 5},
		{0, -1, 0},
		{9},
		{},
	}
	for _, in := range tests {
		o1, o2 := Tee(context.Background(), closed(in...))
		a, b := drainBoth(t, o1, o2)
		if !slices.Equal(a, in) || !slices.Equal(b, in) {
			t.Errorf("Tee(%v) = (%v, %v), want (%v, %v)", in, a, b, in, in)
		}
	}
	words := []string{"héllo", "", "世界"}
	w1, w2 := Tee(context.Background(), closed(words...))
	a, b := drainBoth(t, w1, w2)
	if !slices.Equal(a, words) || !slices.Equal(b, words) {
		t.Errorf("Tee(%q) = (%q, %q), want (%q, %q)", words, a, b, words, words)
	}
}

func TestTeeMany(t *testing.T) {
	in := make([]int, 2000)
	for i := range in {
		in[i] = i
	}
	o1, o2 := Tee(context.Background(), closed(in...))
	a, b := drainBoth(t, o1, o2)
	if !slices.Equal(a, in) || !slices.Equal(b, in) {
		t.Errorf("Tee(0..1999): outputs got %d and %d values, want 2000 each in order", len(a), len(b))
	}
}

// TestTeeEitherOrder reads the SECOND output first for every value. A Tee
// that always sends to the first output first blocks forever here.
func TestTeeEitherOrder(t *testing.T) {
	a, b := Tee(context.Background(), closed(10, 20, 30))
	notNil(t, a, b)
	recv := func(name string, ch <-chan int, want int) {
		t.Helper()
		select {
		case v, ok := <-ch:
			if !ok || v != want {
				t.Fatalf("Tee([10 20 30]): %s received (%d, %v), want (%d, true)", name, v, ok, want)
			}
		case <-time.After(2 * time.Second):
			t.Fatalf("Tee([10 20 30]): %s got nothing within 2s while waiting for %d (send to whichever output is ready: select with both sends)", name, want)
		}
	}
	for _, want := range []int{10, 20, 30} {
		recv("second output", b, want)
		recv("first output", a, want)
	}
	gotA, gotB := drainBoth(t, a, b)
	if len(gotA) != 0 || len(gotB) != 0 {
		t.Errorf("Tee([10 20 30]) after 3 values = (%v, %v), want both closed", gotA, gotB)
	}
}

// mustBeClosed checks, after synctest.Wait, that ch is closed and empty.
func mustBeClosed(t *testing.T, name string, ch <-chan int) {
	t.Helper()
	select {
	case v, ok := <-ch:
		if ok {
			t.Fatalf("after cancel: %s received %d, want it closed", name, v)
		}
	default:
		t.Fatalf("after cancel: %s is still open, want it closed (select on ctx.Done() everywhere you block)", name)
	}
}

// release closes in and drains both outputs, so a Tee that ignores ctx
// fails with a message instead of leaving a blocked goroutine in the bubble.
func release(in chan int, a, b <-chan int) {
	close(in)
	for _, ch := range []<-chan int{a, b} {
		go func() {
			for range ch {
			}
		}()
	}
}

// TestTeeCancelWhileIdle: in never closes; cancelling ctx must close both
// outputs and end the goroutine (synctest fails if it leaks).
func TestTeeCancelWhileIdle(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		ctx, cancel := context.WithCancel(t.Context())
		in := make(chan int)
		a, b := Tee(ctx, in)
		notNil(t, a, b)
		defer release(in, a, b)
		synctest.Wait()
		cancel()
		synctest.Wait()
		mustBeClosed(t, "first output", a)
		mustBeClosed(t, "second output", b)
	})
}

// TestTeeCancelWhileSending: the first output took the value, nobody reads
// the second. Cancelling must still close both.
func TestTeeCancelWhileSending(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		ctx, cancel := context.WithCancel(t.Context())
		in := make(chan int)
		a, b := Tee(ctx, in)
		notNil(t, a, b)
		defer release(in, a, b)
		select {
		case in <- 7:
		case <-time.After(time.Hour):
			t.Fatalf("Tee did not receive from in")
		}
		select {
		case v := <-a:
			if v != 7 {
				t.Fatalf("Tee: first output received %d, want 7", v)
			}
		case <-time.After(time.Hour):
			t.Fatalf("Tee: first output got nothing while the second is not being read (send to whichever output is ready)")
		}
		synctest.Wait()
		cancel()
		synctest.Wait()
		mustBeClosed(t, "first output", a)
		mustBeClosed(t, "second output", b)
	})
}

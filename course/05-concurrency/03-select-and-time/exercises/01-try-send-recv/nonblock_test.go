package nonblock

import (
	"slices"
	"testing"
	"testing/synctest"
	"time"
)

// within runs f in a goroutine and fails the test if it hasn't returned
// after 2s, which means f blocked.
func within(t *testing.T, call string, f func()) {
	t.Helper()
	done := make(chan struct{})
	go func() {
		defer close(done)
		f()
	}()
	select {
	case <-done:
	case <-time.After(2 * time.Second):
		t.Fatalf("%s blocked: did not return within 2s", call)
	}
}

func TestTrySend(t *testing.T) {
	ch := make(chan int, 2)
	var got bool
	within(t, "TrySend(empty buffer of 2, 1)", func() { got = TrySend(ch, 1) })
	if !got {
		t.Fatalf("TrySend(empty buffer of 2, 1) = false, want true")
	}
	within(t, "TrySend(buffer 1/2, 2)", func() { got = TrySend(ch, 2) })
	if !got {
		t.Fatalf("TrySend(buffer 1/2, 2) = false, want true")
	}
	within(t, "TrySend(full buffer, 3)", func() { got = TrySend(ch, 3) })
	if got {
		t.Fatalf("TrySend(full buffer, 3) = true, want false")
	}
	if n := len(ch); n != 2 {
		t.Fatalf("after TrySend 1, 2, 3 on a buffer of 2: len(ch) = %d, want 2", n)
	}
	if a, b := <-ch, <-ch; a != 1 || b != 2 {
		t.Fatalf("buffer holds %d %d, want 1 2", a, b)
	}

	within(t, "TrySend(unbuffered, no receiver, 1)", func() { got = TrySend(make(chan int), 1) })
	if got {
		t.Errorf("TrySend(unbuffered, no receiver, 1) = true, want false")
	}
	within(t, "TrySend(nil, 1)", func() { got = TrySend[int](nil, 1) })
	if got {
		t.Errorf("TrySend(nil, 1) = true, want false")
	}
}

func TestTrySendToWaitingReceiver(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		ch := make(chan int)
		got := make(chan int, 1)
		go func() { got <- <-ch }()
		synctest.Wait() // the receiver is now blocked in <-ch

		ok := TrySend(ch, 7)
		synctest.Wait()
		select {
		case v := <-got:
			if !ok || v != 7 {
				t.Errorf("TrySend(unbuffered with a waiting receiver, 7) = %v and the receiver got %d, want true and 7", ok, v)
			}
		default:
			t.Errorf("TrySend(unbuffered with a waiting receiver, 7) = %v and the receiver got nothing, want true and 7", ok)
			ch <- 0 // release the receiver so the test can finish
		}
	})
}

func TestTryRecv(t *testing.T) {
	ch := make(chan string, 2)
	ch <- "a"
	ch <- "b"
	for _, want := range []string{"a", "b"} {
		var v string
		var ok bool
		within(t, "TryRecv(ch)", func() { v, ok = TryRecv(ch) })
		if v != want || !ok {
			t.Fatalf("TryRecv(ch) = (%q, %v), want (%q, true)", v, ok, want)
		}
	}

	cases := []struct {
		name string
		ch   chan string
	}{
		{"empty buffered", make(chan string, 1)},
		{"unbuffered, no sender", make(chan string)},
		{"nil", nil},
		{"closed", func() chan string { c := make(chan string); close(c); return c }()},
	}
	for _, c := range cases {
		var v string
		var ok bool
		within(t, "TryRecv("+c.name+")", func() { v, ok = TryRecv(c.ch) })
		if v != "" || ok {
			t.Errorf("TryRecv(%s) = (%q, %v), want (\"\", false)", c.name, v, ok)
		}
	}
}

func TestTryRecvFromWaitingSender(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		ch := make(chan int)
		go func() { ch <- 7 }()
		synctest.Wait() // the sender is now blocked in ch <- 7

		v, ok := TryRecv(ch)
		if v != 7 || !ok {
			t.Errorf("TryRecv(unbuffered with a waiting sender of 7) = (%d, %v), want (7, true)", v, ok)
		}
		synctest.Wait()
		select {
		case <-ch: // release the sender if TryRecv didn't take the value
		default:
		}
	})
}

func TestDrainReady(t *testing.T) {
	ch := make(chan int, 5)
	for _, v := range []int{1, 2, 3} {
		ch <- v
	}
	var got []int
	within(t, "DrainReady(ch)", func() { got = DrainReady(ch) })
	if want := []int{1, 2, 3}; !slices.Equal(got, want) {
		t.Fatalf("DrainReady(ch holding 1 2 3) = %v, want %v", got, want)
	}
	if n := len(ch); n != 0 {
		t.Fatalf("after DrainReady: len(ch) = %d, want 0", n)
	}

	within(t, "DrainReady(empty)", func() { got = DrainReady(ch) })
	if got != nil {
		t.Errorf("DrainReady(empty) = %v, want nil", got)
	}

	ch <- 9
	close(ch)
	within(t, "DrainReady(closed, holding 9)", func() { got = DrainReady(ch) })
	if want := []int{9}; !slices.Equal(got, want) {
		t.Errorf("DrainReady(closed, holding 9) = %v, want %v", got, want)
	}

	within(t, "DrainReady(nil)", func() { got = DrainReady[int](nil) })
	if got != nil {
		t.Errorf("DrainReady(nil) = %v, want nil", got)
	}
}

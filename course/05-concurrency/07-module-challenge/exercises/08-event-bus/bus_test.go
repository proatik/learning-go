package bus

import (
	"errors"
	"fmt"
	"slices"
	"sync"
	"sync/atomic"
	"testing"
	"testing/synctest"
	"time"
)

// recorder is a handler that records every message, optionally taking
// delay per message.
type recorder struct {
	mu    sync.Mutex
	msgs  []string
	delay time.Duration
}

func (r *recorder) handle(msg string) {
	if r.delay > 0 {
		time.Sleep(r.delay)
	}
	r.mu.Lock()
	r.msgs = append(r.msgs, msg)
	r.mu.Unlock()
}

func (r *recorder) got() []string {
	r.mu.Lock()
	defer r.mu.Unlock()
	return slices.Clone(r.msgs)
}

func TestReadmeScenario(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		b := New(2)
		slow := &recorder{delay: time.Second}
		s := b.Subscribe("news", slow.handle)
		start := time.Now()
		if err := b.Publish("news", "m1"); err != nil {
			t.Fatalf(`Publish("news", "m1") error = %v, want <nil>`, err)
		}
		synctest.Wait() // the handler is now busy with m1
		for _, m := range []string{"m2", "m3", "m4", "m5"} {
			b.Publish("news", m)
		}
		if el := time.Since(start); el != 0 {
			t.Fatalf("5 Publish calls to a busy subscriber took %v, want 0s (Publish must not block)", el)
		}
		if d := s.Dropped(); d != 2 {
			t.Errorf("Dropped() with queue size 2 after m2..m5 = %d, want 2", d)
		}
		b.Close()
		if el := time.Since(start); el != 3*time.Second {
			t.Errorf("Close() returned after %v, want 3s (it must wait for m1, m4, m5)", el)
		}
		if got, want := slow.got(), []string{"m1", "m4", "m5"}; !slices.Equal(got, want) {
			t.Errorf("handler received %q, want %q (oldest dropped)", got, want)
		}
		if err := b.Publish("news", "m6"); !errors.Is(err, ErrClosed) {
			t.Errorf(`Publish("news", "m6") after Close: error = %v, want %v`, err, ErrClosed)
		}
		b.Close() // a second Close must not panic or block
	})
}

func TestTopicsAndFanOut(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		b := New(10)
		a1, a2, other := &recorder{}, &recorder{}, &recorder{}
		b.Subscribe("a", a1.handle)
		b.Subscribe("a", a2.handle)
		b.Subscribe("b", other.handle)
		for _, p := range [][2]string{{"a", "x1"}, {"b", "y1"}, {"a", "x2"}, {"c", "z1"}, {"a", "x3"}} {
			if err := b.Publish(p[0], p[1]); err != nil {
				t.Fatalf("Publish(%q, %q) error = %v, want <nil>", p[0], p[1], err)
			}
		}
		b.Close()
		want := []string{"x1", "x2", "x3"}
		if got := a1.got(); !slices.Equal(got, want) {
			t.Errorf(`first subscriber of "a" received %q, want %q`, got, want)
		}
		if got := a2.got(); !slices.Equal(got, want) {
			t.Errorf(`second subscriber of "a" received %q, want %q`, got, want)
		}
		if got := other.got(); !slices.Equal(got, []string{"y1"}) {
			t.Errorf(`subscriber of "b" received %q, want ["y1"]`, got)
		}
	})
}

func TestSlowSubscriberDoesNotDelayOthers(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		b := New(10)
		slow, fast := &recorder{delay: time.Second}, &recorder{}
		b.Subscribe("t", slow.handle)
		b.Subscribe("t", fast.handle)
		start := time.Now()
		for i := range 3 {
			b.Publish("t", fmt.Sprint(i))
		}
		synctest.Wait()
		if got := fast.got(); len(got) != 3 {
			t.Errorf("fast subscriber received %q at %v, want 3 messages at once", got, time.Since(start))
		}
		b.Close()
		if got := slow.got(); len(got) != 3 {
			t.Errorf("slow subscriber received %q after Close, want 3 messages", got)
		}
		if el := time.Since(start); el != 3*time.Second {
			t.Errorf("Close() returned after %v, want 3s", el)
		}
	})
}

func TestUnsubscribe(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		b := New(5)
		gone, stays := &recorder{delay: time.Second}, &recorder{}
		s := b.Subscribe("t", gone.handle)
		b.Subscribe("t", stays.handle)
		for _, m := range []string{"a", "b", "c"} {
			b.Publish("t", m)
		}
		start := time.Now()
		s.Unsubscribe()
		if el := time.Since(start); el != 3*time.Second {
			t.Errorf("Unsubscribe() returned after %v, want 3s (finish the 3 queued messages first)", el)
		}
		if got := gone.got(); !slices.Equal(got, []string{"a", "b", "c"}) {
			t.Errorf("unsubscribed handler received %q, want [a b c]", got)
		}
		if err := b.Publish("t", "d"); err != nil {
			t.Errorf(`Publish("t", "d") after one Unsubscribe: error = %v, want <nil>`, err)
		}
		s.Unsubscribe() // a second Unsubscribe must not panic or block
		b.Close()
		if got := gone.got(); len(got) != 3 {
			t.Errorf("unsubscribed handler received %q, want no messages after Unsubscribe", got)
		}
		if got := stays.got(); !slices.Equal(got, []string{"a", "b", "c", "d"}) {
			t.Errorf("remaining subscriber received %q, want [a b c d]", got)
		}
		s.Unsubscribe() // after Close: must not panic or block
	})
}

func TestSubscribeAfterClose(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		b := New(1)
		b.Close()
		r := &recorder{}
		s := b.Subscribe("t", r.handle)
		if s == nil {
			t.Fatalf(`Subscribe("t", h) after Close = nil, want a subscription`)
		}
		if err := b.Publish("t", "x"); !errors.Is(err, ErrClosed) {
			t.Errorf(`Publish("t", "x") after Close: error = %v, want %v`, err, ErrClosed)
		}
		s.Unsubscribe()
		if got := r.got(); len(got) != 0 {
			t.Errorf("handler subscribed after Close received %q, want nothing", got)
		}
	})
}

func TestQueueSizeBelowOne(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		b := New(0)
		r := &recorder{delay: time.Second}
		s := b.Subscribe("t", r.handle)
		b.Publish("t", "m1")
		synctest.Wait()
		b.Publish("t", "m2")
		b.Publish("t", "m3")
		b.Close()
		if got := r.got(); !slices.Equal(got, []string{"m1", "m3"}) {
			t.Errorf("New(0) (queue 1): handler received %q, want [m1 m3]", got)
		}
		if d := s.Dropped(); d != 1 {
			t.Errorf("New(0) (queue 1): Dropped() = %d, want 1", d)
		}
	})
}

func TestEveryMessageHandledOrDropped(t *testing.T) {
	const publishers, each = 4, 1000
	b := New(3)
	gate := make(chan struct{})
	var handled atomic.Int64
	s := b.Subscribe("load", func(string) {
		<-gate
		handled.Add(1)
	})
	published := make(chan struct{})
	go func() {
		var wg sync.WaitGroup
		for p := range publishers {
			wg.Go(func() {
				for i := range each {
					b.Publish("load", fmt.Sprintf("%d-%d", p, i))
				}
			})
		}
		wg.Wait()
		close(published)
	}()
	select {
	case <-published:
	case <-time.After(2 * time.Second):
		close(gate)
		t.Fatalf("%d publishers x %d messages to a stuck subscriber did not finish within 2s (Publish blocked?)", publishers, each)
	}
	close(gate)
	closed := make(chan struct{})
	go func() {
		b.Close()
		close(closed)
	}()
	select {
	case <-closed:
	case <-time.After(2 * time.Second):
		t.Fatalf("Close() did not return within 2s")
	}
	h, d := handled.Load(), int64(s.Dropped())
	if h+d != publishers*each {
		t.Errorf("handled (%d) + Dropped() (%d) = %d, want %d (every message is handled or dropped)", h, d, h+d, publishers*each)
	}
	if h < 1 || h > 4 {
		t.Errorf("handled = %d, want 1 to 4 (the one in the handler plus at most 3 queued)", h)
	}
}

package pubsub

import (
	"slices"
	"sync"
	"testing"
	"testing/synctest"
)

// pending takes every value already buffered in ch without blocking and
// reports whether ch is closed.
func pending[T any](ch <-chan T) (vals []T, closed bool) {
	for {
		select {
		case v, ok := <-ch:
			if !ok {
				return vals, true
			}
			vals = append(vals, v)
		default:
			return vals, false
		}
	}
}

func subscribe[T any](t *testing.T, b *Broker[T]) (<-chan T, func()) {
	t.Helper()
	ch, unsub := b.Subscribe()
	if ch == nil || unsub == nil {
		t.Fatalf("Subscribe() returned a nil channel or func")
	}
	return ch, unsub
}

// noPanic runs f and reports a panic as a test failure.
func noPanic(t *testing.T, name string, f func()) {
	t.Helper()
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("%s panicked: %v, want no panic", name, r)
		}
	}()
	f()
}

func TestPublishToAll(t *testing.T) {
	b := New[int](10)
	if n := b.Publish(99); n != 0 {
		t.Errorf("Publish(99) with no subscribers = %d, want 0", n)
	}
	s1, _ := subscribe(t, b)
	s2, _ := subscribe(t, b)
	if n := b.Subscribers(); n != 2 {
		t.Errorf("Subscribers() = %d, want 2", n)
	}
	for _, v := range []int{1, 2, 3} {
		if n := b.Publish(v); n != 2 {
			t.Errorf("Publish(%d) with 2 subscribers = %d, want 2", v, n)
		}
	}
	for i, ch := range []<-chan int{s1, s2} {
		got, closed := pending(ch)
		if !slices.Equal(got, []int{1, 2, 3}) || closed {
			t.Errorf("subscriber %d received %v (closed %v), want [1 2 3] (open)", i+1, got, closed)
		}
	}
	s3, _ := subscribe(t, b)
	b.Publish(4)
	if got, _ := pending(s3); !slices.Equal(got, []int{4}) {
		t.Errorf("late subscriber received %v, want [4] (only values published after Subscribe)", got)
	}
}

func TestUnsubscribe(t *testing.T) {
	b := New[string](4)
	s1, unsub1 := subscribe(t, b)
	s2, _ := subscribe(t, b)
	b.Publish("a")
	unsub1()
	if n := b.Subscribers(); n != 1 {
		t.Errorf("Subscribers() after one unsubscribe = %d, want 1", n)
	}
	if n := b.Publish("b"); n != 1 {
		t.Errorf("Publish(%q) after one unsubscribe = %d, want 1", "b", n)
	}
	got, closed := pending(s1)
	if !slices.Equal(got, []string{"a"}) || !closed {
		t.Errorf("unsubscribed channel yields %q (closed %v), want [a] then closed", got, closed)
	}
	if got, _ := pending(s2); !slices.Equal(got, []string{"a", "b"}) {
		t.Errorf("remaining subscriber received %q, want [a b]", got)
	}
	noPanic(t, "calling unsubscribe twice", unsub1)
	if n := b.Subscribers(); n != 1 {
		t.Errorf("Subscribers() after a second unsubscribe = %d, want 1", n)
	}
}

// TestSlowSubscriberMisses: a full buffer drops the value for that
// subscriber only; Publish never blocks.
func TestSlowSubscriberMisses(t *testing.T) {
	b := New[int](2)
	slow, _ := subscribe(t, b)
	fast, _ := subscribe(t, b)
	var fastGot []int
	wantN := []int{2, 2, 1, 1}
	for i, v := range []int{1, 2, 3, 4} {
		if n := b.Publish(v); n != wantN[i] {
			t.Errorf("Publish(%d) = %d, want %d (the slow subscriber's buffer of 2 is full)", v, n, wantN[i])
		}
		got, _ := pending(fast)
		fastGot = append(fastGot, got...)
	}
	if got, _ := pending(slow); !slices.Equal(got, []int{1, 2}) {
		t.Errorf("slow subscriber (never reads, buffer 2) received %v, want [1 2]", got)
	}
	if !slices.Equal(fastGot, []int{1, 2, 3, 4}) {
		t.Errorf("fast subscriber received %v, want [1 2 3 4]", fastGot)
	}
}

// TestUnbuffered: with buffer 0 a subscriber receives only while it is
// waiting in a receive.
func TestUnbuffered(t *testing.T) {
	for _, buffer := range []int{0, -3} {
		synctest.Test(t, func(t *testing.T) {
			b := New[int](buffer)
			ch, unsub := subscribe(t, b)
			if n := b.Publish(1); n != 0 {
				t.Errorf("New(%d): Publish(1) with nobody receiving = %d, want 0", buffer, n)
			}
			got := make(chan int, 1)
			go func() {
				v, ok := <-ch
				if ok {
					got <- v
				}
				close(got)
			}()
			synctest.Wait() // the goroutine is now blocked receiving
			if n := b.Publish(2); n != 1 {
				t.Errorf("New(%d): Publish(2) with a waiting receiver = %d, want 1", buffer, n)
			}
			unsub() // ends the goroutine if Publish didn't reach it
			if v, ok := <-got; !ok || v != 2 {
				t.Errorf("New(%d): waiting receiver got (%d, %v), want (2, true)", buffer, v, ok)
			}
		})
	}
}

func TestClose(t *testing.T) {
	b := New[int](4)
	s1, unsub1 := subscribe(t, b)
	s2, _ := subscribe(t, b)
	b.Publish(7)
	b.Close()
	for i, ch := range []<-chan int{s1, s2} {
		got, closed := pending(ch)
		if !slices.Equal(got, []int{7}) || !closed {
			t.Errorf("after Close: subscriber %d yields %v (closed %v), want [7] then closed", i+1, got, closed)
		}
	}
	if n := b.Subscribers(); n != 0 {
		t.Errorf("Subscribers() after Close = %d, want 0", n)
	}
	noPanic(t, "Publish after Close", func() {
		if n := b.Publish(8); n != 0 {
			t.Errorf("Publish(8) after Close = %d, want 0", n)
		}
	})
	noPanic(t, "unsubscribe after Close", unsub1)
	noPanic(t, "a second Close", b.Close)
	noPanic(t, "Subscribe after Close", func() {
		ch, unsub := b.Subscribe()
		if ch == nil || unsub == nil {
			t.Fatalf("Subscribe() after Close returned a nil channel or func")
		}
		if _, closed := pending(ch); !closed {
			t.Errorf("Subscribe() after Close returned an open channel, want a closed one")
		}
		unsub()
	})
	if n := b.Subscribers(); n != 0 {
		t.Errorf("Subscribers() after Subscribe on a closed broker = %d, want 0", n)
	}
}

// TestConcurrent: 50 publishers × 100 values reach every steady subscriber
// exactly once, while other goroutines subscribe and unsubscribe.
func TestConcurrent(t *testing.T) {
	const subs, publishers, per = 10, 50, 100
	b := New[int](publishers * per)
	var chans []<-chan int
	for range subs {
		ch, _ := subscribe(t, b)
		chans = append(chans, ch)
	}
	var wg sync.WaitGroup
	var mu sync.Mutex
	delivered := 0
	for p := range publishers {
		wg.Go(func() {
			n := 0
			for i := range per {
				n += b.Publish(p*per + i)
			}
			mu.Lock()
			delivered += n
			mu.Unlock()
		})
	}
	for range 10 {
		wg.Go(func() {
			for range 100 {
				ch, unsub := b.Subscribe()
				if ch == nil || unsub == nil {
					return
				}
				unsub()
				unsub()
			}
		})
	}
	wg.Wait()
	if n := b.Subscribers(); n != subs {
		t.Errorf("Subscribers() after churn = %d, want %d", n, subs)
	}
	if delivered < subs*publishers*per {
		t.Errorf("Publish deliveries to steady subscribers = %d, want at least %d", delivered, subs*publishers*per)
	}
	for i, ch := range chans {
		got, _ := pending(ch)
		slices.Sort(got)
		if len(got) != publishers*per {
			t.Errorf("subscriber %d received %d values, want %d", i+1, len(got), publishers*per)
			continue
		}
		for j, v := range got {
			if v != j {
				t.Errorf("subscriber %d: value %d missing or duplicated", i+1, j)
				break
			}
		}
	}
	b.Close()
}

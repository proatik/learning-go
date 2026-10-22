package pool

import (
	"context"
	"errors"
	"slices"
	"sync"
	"sync/atomic"
	"testing"
	"testing/synctest"
	"time"
)

type result struct {
	v   string
	err error
}

func TestNewIdle(t *testing.T) {
	if got := New([]string{"a", "b"}).Idle(); got != 2 {
		t.Errorf(`New([a b]).Idle() = %d, want 2`, got)
	}
	if got := New[int](nil).Idle(); got != 0 {
		t.Errorf(`New(nil).Idle() = %d, want 0`, got)
	}
}

func TestGetAndPut(t *testing.T) {
	p := New([]string{"a", "b"})
	var got []string
	for range 2 {
		v, err := p.Get(context.Background())
		if err != nil {
			t.Fatalf("Get() error = %v, want nil", err)
		}
		got = append(got, v)
	}
	slices.Sort(got)
	if !slices.Equal(got, []string{"a", "b"}) {
		t.Fatalf("two Get() calls on New([a b]) = %v, want [a b] in any order", got)
	}
	if n := p.Idle(); n != 0 {
		t.Errorf("Idle() after borrowing both = %d, want 0", n)
	}
	p.Put("a")
	p.Put("b")
	if n := p.Idle(); n != 2 {
		t.Errorf("Idle() after putting both back = %d, want 2", n)
	}
}

func TestPutIntoFullPoolPanics(t *testing.T) {
	p := New([]int{1})
	defer func() {
		r := recover()
		if r != "pool: Put into a full pool" {
			t.Errorf("Put(2) into a full pool: panic = %v, want %q", r, "pool: Put into a full pool")
		}
	}()
	p.Put(2)
}

func TestGetWaitsForPut(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		p := New([]string{"a"})
		if _, err := p.Get(context.Background()); err != nil {
			t.Fatalf("Get() error = %v, want nil", err)
		}
		results := make(chan result, 1)
		go func() {
			v, err := p.Get(context.Background())
			results <- result{v, err}
		}()
		synctest.Wait()
		if len(results) != 0 {
			t.Fatalf("Get() on an empty pool returned %v, want it to block until Put", <-results)
		}
		p.Put("a")
		synctest.Wait()
		select {
		case r := <-results:
			if r.v != "a" || r.err != nil {
				t.Errorf(`waiting Get() after Put("a") = (%q, %v), want ("a", <nil>)`, r.v, r.err)
			}
		default:
			t.Fatalf(`waiting Get() did not return after Put("a")`)
		}
	})
}

func TestGetTimeout(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		p := New[string](nil)
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()
		start := time.Now()
		v, err := p.Get(ctx)
		if !errors.Is(err, context.DeadlineExceeded) {
			t.Errorf("Get(ctx with 1s timeout) on an empty pool: error = %v, want %v", err, context.DeadlineExceeded)
		}
		if v != "" {
			t.Errorf("Get(ctx with 1s timeout) on an empty pool = %q, want the zero value", v)
		}
		if el := time.Since(start); el != time.Second {
			t.Errorf("Get(ctx with 1s timeout) returned after %v, want 1s", el)
		}
	})
}

func TestClose(t *testing.T) {
	p := New([]int{7})
	p.Close()
	if v, err := p.Get(context.Background()); !errors.Is(err, ErrClosed) || v != 0 {
		t.Errorf("Get() after Close = (%d, %v), want (0, %v)", v, err, ErrClosed)
	}
	p.Close() // a second Close must not panic

	p = New([]int{7})
	v, _ := p.Get(context.Background())
	p.Close()
	p.Put(v) // ignored, must not panic
	if n := p.Idle(); n != 0 {
		t.Errorf("Idle() after Put on a closed pool = %d, want 0", n)
	}
}

func TestCloseWakesWaiters(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		p := New[string](nil)
		results := make(chan result, 3)
		for range 3 {
			go func() {
				v, err := p.Get(context.Background())
				results <- result{v, err}
			}()
		}
		synctest.Wait()
		if len(results) != 0 {
			t.Fatalf("Get() on an empty pool returned %v, want it to block", <-results)
		}
		p.Close()
		synctest.Wait()
		if len(results) != 3 {
			t.Fatalf("after Close, waiting Get() calls returned = %d, want 3", len(results))
		}
		for range 3 {
			if r := <-results; !errors.Is(r.err, ErrClosed) {
				t.Errorf("waiting Get() after Close: error = %v, want %v", r.err, ErrClosed)
			}
		}
	})
}

func TestLimitsConcurrency(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		p := New([]int{1, 2, 3})
		var running, maxRunning atomic.Int64
		var mu sync.Mutex
		inUse := map[int]bool{}
		lentTwice := 0
		start := time.Now()
		var wg sync.WaitGroup
		for range 12 {
			wg.Go(func() {
				r, err := p.Get(context.Background())
				if err != nil {
					return
				}
				mu.Lock()
				if inUse[r] {
					lentTwice++
				}
				inUse[r] = true
				mu.Unlock()
				n := running.Add(1)
				for {
					m := maxRunning.Load()
					if n <= m || maxRunning.CompareAndSwap(m, n) {
						break
					}
				}
				time.Sleep(time.Second)
				running.Add(-1)
				mu.Lock()
				inUse[r] = false
				mu.Unlock()
				p.Put(r)
			})
		}
		wg.Wait()
		if lentTwice != 0 {
			t.Errorf("resources lent to two goroutines at once = %d, want 0", lentTwice)
		}
		if m := maxRunning.Load(); m != 3 {
			t.Errorf("12 workers, 3 resources: max holding a resource at once = %d, want 3", m)
		}
		if el := time.Since(start); el != 4*time.Second {
			t.Errorf("12 tasks of 1s, 3 resources: took %v, want 4s", el)
		}
	})
}

func TestHammer(t *testing.T) {
	const workers, rounds = 50, 200
	p := New([]int{1, 2, 3, 4})
	var held [5]atomic.Int32
	var gets, clashes atomic.Int64
	done := make(chan struct{})
	go func() {
		var wg sync.WaitGroup
		for range workers {
			wg.Go(func() {
				for range rounds {
					r, err := p.Get(context.Background())
					if err != nil || r < 1 || r > 4 {
						clashes.Add(1)
						continue
					}
					if held[r].Add(1) != 1 {
						clashes.Add(1)
					}
					gets.Add(1)
					held[r].Add(-1)
					p.Put(r)
				}
			})
		}
		wg.Wait()
		close(done)
	}()
	select {
	case <-done:
	case <-time.After(2 * time.Second):
		t.Fatalf("%d goroutines x %d Get/Put did not finish within 2s", workers, rounds)
	}
	if c := clashes.Load(); c != 0 {
		t.Errorf("bad or shared resources seen = %d, want 0", c)
	}
	if g := gets.Load(); g != workers*rounds {
		t.Errorf("successful Get() calls = %d, want %d", g, workers*rounds)
	}
	if n := p.Idle(); n != 4 {
		t.Errorf("Idle() after all workers finished = %d, want 4", n)
	}
}

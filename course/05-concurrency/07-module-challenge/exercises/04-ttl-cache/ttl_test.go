package ttl

import (
	"sync"
	"testing"
	"testing/synctest"
	"time"
)

func wantGet(t *testing.T, step string, c *Cache[string, int], k string, want int, wantOK bool) {
	t.Helper()
	if v, ok := c.Get(k); v != want || ok != wantOK {
		t.Fatalf("%s: Get(%q) = (%d, %v), want (%d, %v)", step, k, v, ok, want, wantOK)
	}
}

func wantLen(t *testing.T, step string, c *Cache[string, int], want int) {
	t.Helper()
	if got := c.Len(); got != want {
		t.Fatalf("%s: Len() = %d, want %d", step, got, want)
	}
}

func TestSetGetExpire(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		c := New[string, int](time.Hour)
		defer c.Close()
		c.Set("a", 1, 10*time.Second)
		wantGet(t, `Set("a", 1, 10s)`, c, "a", 1, true)
		wantGet(t, "missing key", c, "zz", 0, false)
		wantLen(t, `Set("a", 1, 10s)`, c, 1)
		time.Sleep(10*time.Second - time.Millisecond)
		wantGet(t, "9.999s later", c, "a", 1, true)
		time.Sleep(time.Millisecond)
		wantGet(t, "10s later (expired)", c, "a", 0, false)
		wantLen(t, "10s later, before any sweep", c, 1)
	})
}

func TestZeroAndNegativeTTL(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		c := New[string, int](time.Hour)
		defer c.Close()
		c.Set("zero", 1, 0)
		c.Set("neg", 2, -time.Second)
		wantGet(t, `Set("zero", 1, 0)`, c, "zero", 0, false)
		wantGet(t, `Set("neg", 2, -1s)`, c, "neg", 0, false)
	})
}

func TestOverwriteResetsExpiry(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		c := New[string, int](time.Hour)
		defer c.Close()
		c.Set("a", 1, time.Second)
		c.Set("a", 2, time.Minute)
		time.Sleep(2 * time.Second)
		wantGet(t, `Set("a", 1, 1s), Set("a", 2, 1m), 2s later`, c, "a", 2, true)
		wantLen(t, "after overwrite", c, 1)
	})
}

func TestDelete(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		c := New[string, int](time.Hour)
		defer c.Close()
		c.Set("a", 1, time.Minute)
		c.Set("b", 2, time.Minute)
		c.Delete("a")
		c.Delete("missing")
		wantGet(t, `Delete("a")`, c, "a", 0, false)
		wantGet(t, `Delete("a")`, c, "b", 2, true)
		wantLen(t, `Delete("a")`, c, 1)
	})
}

func TestJanitorSweeps(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		c := New[string, int](time.Minute)
		defer c.Close()
		c.Set("a", 1, 10*time.Second)
		c.Set("b", 2, 90*time.Second)
		c.Set("c", 3, time.Hour)
		time.Sleep(59 * time.Second)
		synctest.Wait()
		wantLen(t, "59s, before the first sweep", c, 3)
		time.Sleep(time.Second)
		synctest.Wait()
		wantLen(t, "first sweep at 1m (a expired)", c, 2)
		wantGet(t, "after the first sweep", c, "b", 2, true)
		time.Sleep(time.Minute)
		synctest.Wait()
		wantLen(t, "second sweep at 2m (b expired)", c, 1)
		wantGet(t, "after the second sweep", c, "c", 3, true)
	})
}

func TestDefaultSweepInterval(t *testing.T) {
	for _, every := range []time.Duration{0, -time.Minute} {
		synctest.Test(t, func(t *testing.T) {
			c := New[string, int](every)
			defer c.Close()
			c.Set("a", 1, time.Millisecond)
			time.Sleep(time.Second)
			synctest.Wait()
			if got := c.Len(); got != 0 {
				t.Fatalf("New(%v): Len() 1s after an entry expired = %d, want 0 (sweep every 1s)", every, got)
			}
		})
	}
}

func TestCloseStopsJanitor(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		c := New[string, int](time.Second)
		c.Set("a", 1, 10*time.Millisecond)
		time.Sleep(time.Second)
		synctest.Wait()
		wantLen(t, "before Close", c, 0)
		c.Close()
		c.Close() // a second Close must not panic or block
		c.Set("b", 2, 10*time.Millisecond)
		time.Sleep(5 * time.Second)
		synctest.Wait()
		wantLen(t, "5s after Close (no more sweeps)", c, 1)
		wantGet(t, "after Close", c, "b", 0, false)
	})
}

func TestConcurrentUse(t *testing.T) {
	const workers, each = 50, 1000
	c := New[string, int](10 * time.Millisecond)
	defer c.Close()
	done := make(chan struct{})
	go func() {
		var wg sync.WaitGroup
		for w := range workers {
			wg.Go(func() {
				for i := range each {
					k := string(rune('A'+w)) + string(rune('0'+i%10)) + string(rune(i))
					c.Set(k, i, time.Hour)
					c.Get(k)
					c.Len()
				}
			})
		}
		wg.Wait()
		close(done)
	}()
	select {
	case <-done:
	case <-time.After(2 * time.Second):
		t.Fatalf("%d goroutines x %d Set/Get did not finish within 2s", workers, each)
	}
	if got := c.Len(); got != workers*each {
		t.Errorf("Len() after %d concurrent Sets of distinct keys = %d, want %d", workers*each, got, workers*each)
	}
}

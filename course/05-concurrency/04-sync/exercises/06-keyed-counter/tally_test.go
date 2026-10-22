package tally

import (
	"fmt"
	"maps"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

// run calls f in a goroutine and fails the test if it doesn't return within
// 2s, so a self-deadlock fails the test instead of hanging it.
func run(t *testing.T, what string, f func()) {
	t.Helper()
	done := make(chan struct{})
	go func() {
		defer close(done)
		f()
	}()
	select {
	case <-done:
	case <-time.After(2 * time.Second):
		t.Fatalf("%s did not finish within 2s (deadlock?)", what)
	}
}

func TestSequential(t *testing.T) {
	var c Counter
	run(t, "Inc/Add/Get", func() {
		steps := []struct {
			call string
			got  int
			want int
		}{
			{`Get("go")`, c.Get("go"), 0},
			{`Inc("go")`, c.Inc("go"), 1},
			{`Inc("go")`, c.Inc("go"), 2},
			{`Add("go", 10)`, c.Add("go", 10), 12},
			{`Inc("rust")`, c.Inc("rust"), 1},
			{`Add("rust", -3)`, c.Add("rust", -3), -2},
			{`Inc("")`, c.Inc(""), 1},
			{`Inc("héllo")`, c.Inc("héllo"), 1},
			{`Get("go")`, c.Get("go"), 12},
			{`Get("missing")`, c.Get("missing"), 0},
		}
		for _, s := range steps {
			if s.got != s.want {
				t.Errorf("%s = %d, want %d", s.call, s.got, s.want)
			}
		}
	})
}

func TestSnapshot(t *testing.T) {
	var c Counter
	run(t, "Snapshot", func() {
		if got := c.Snapshot(); got == nil || len(got) != 0 {
			t.Errorf("Counter{}.Snapshot() = %#v, want map[string]int{}", got)
		}
		c.Add("a", 2)
		c.Inc("b")
		want := map[string]int{"a": 2, "b": 1}
		snap := c.Snapshot()
		if !maps.Equal(snap, want) {
			t.Fatalf("Snapshot() = %v, want %v", snap, want)
		}
		snap["a"] = 100
		snap["z"] = 1
		delete(snap, "b")
		if got := c.Snapshot(); !maps.Equal(got, want) {
			t.Errorf("after changing the returned map: Snapshot() = %v, want %v (return a copy)", got, want)
		}
		snap2 := c.Snapshot()
		c.Inc("a")
		c.Inc("c")
		if !maps.Equal(snap2, want) {
			t.Errorf("after Inc(\"a\"), Inc(\"c\"): earlier snapshot = %v, want it unchanged %v", snap2, want)
		}
	})
}

func TestReset(t *testing.T) {
	var c Counter
	run(t, "Reset", func() {
		if got := c.Reset(); got == nil || len(got) != 0 {
			t.Errorf("Counter{}.Reset() = %#v, want map[string]int{}", got)
		}
		c.Add("a", 5)
		c.Inc("b")
		if got, want := c.Reset(), map[string]int{"a": 5, "b": 1}; !maps.Equal(got, want) {
			t.Errorf("Reset() = %v, want %v", got, want)
		}
		if got := c.Snapshot(); len(got) != 0 {
			t.Errorf("after Reset(): Snapshot() = %v, want map[]", got)
		}
		if got := c.Inc("a"); got != 1 {
			t.Errorf("after Reset(): Inc(\"a\") = %d, want 1", got)
		}
	})
}

func TestConcurrentInc(t *testing.T) {
	const goroutines, perG, keys = 50, 1000, 5
	const perKey = goroutines * perG / keys
	var c Counter
	// seen[k][v-1] records that some Inc returned v for key k: every value
	// 1..perKey must come back exactly once.
	var seen [keys][perKey]atomic.Bool
	var bad atomic.Int64
	run(t, "50 goroutines × 1000 Inc()", func() {
		var wg sync.WaitGroup
		for g := range goroutines {
			wg.Go(func() {
				for i := range perG {
					k := (g + i) % keys
					v := c.Inc(fmt.Sprint("k", k))
					if v < 1 || v > perKey || seen[k][v-1].Swap(true) {
						bad.Add(1)
					}
					if i%100 == 0 {
						_ = len(c.Snapshot())
					}
				}
			})
		}
		wg.Wait()
	})
	if got := bad.Load(); got != 0 {
		t.Errorf("concurrent Inc(): %d return values were duplicated or out of range, want 0 (each Inc returns its own new count)", got)
	}
	want := map[string]int{}
	for k := range keys {
		want[fmt.Sprint("k", k)] = perKey
	}
	if got := c.Snapshot(); !maps.Equal(got, want) {
		t.Errorf("%d goroutines × %d Inc() over %d keys: Snapshot() = %v, want %v", goroutines, perG, keys, got, want)
	}
}

func TestResetLosesNothing(t *testing.T) {
	var c Counter
	total := map[string]int{}
	run(t, "Inc() with a concurrent Reset()", func() {
		stop := make(chan struct{})
		var resetter sync.WaitGroup
		resetter.Go(func() {
			for {
				for k, v := range c.Reset() {
					total[k] += v
				}
				select {
				case <-stop:
					return
				default:
				}
			}
		})
		var wg sync.WaitGroup
		for range 20 {
			wg.Go(func() {
				for i := range 1000 {
					c.Inc(fmt.Sprint("k", i%4))
				}
			})
		}
		wg.Wait()
		close(stop)
		resetter.Wait()
	})
	for k, v := range c.Reset() {
		total[k] += v
	}
	want := map[string]int{"k0": 5000, "k1": 5000, "k2": 5000, "k3": 5000}
	if !maps.Equal(total, want) {
		t.Errorf("sum of every Reset() during 20 goroutines × 1000 Inc() = %v, want %v", total, want)
	}
}

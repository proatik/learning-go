package counter

import (
	"sync"
	"testing"
	"time"
)

// wait fails the test if done isn't closed within 2s.
func wait(t *testing.T, done <-chan struct{}, what string) {
	t.Helper()
	select {
	case <-done:
	case <-time.After(2 * time.Second):
		t.Fatalf("%s did not finish within 2s (deadlock?)", what)
	}
}

// run calls f in a goroutine and fails the test if it doesn't return within
// 2s, so a self-deadlock fails the test instead of hanging it.
func run(t *testing.T, what string, f func()) {
	t.Helper()
	done := make(chan struct{})
	go func() {
		defer close(done)
		f()
	}()
	wait(t, done, what)
}

func TestZeroValue(t *testing.T) {
	var c Counter
	got := -1
	run(t, "Value()", func() { got = c.Value() })
	if got != 0 {
		t.Errorf("Counter{}.Value() = %d, want 0", got)
	}
}

func TestSequential(t *testing.T) {
	var c Counter
	steps := []struct {
		name string
		call func()
		want int
	}{
		{"Inc()", c.Inc, 1},
		{"Inc()", c.Inc, 2},
		{"Add(10)", func() { c.Add(10) }, 12},
		{"Add(-15)", func() { c.Add(-15) }, -3},
		{"Add(0)", func() { c.Add(0) }, -3},
	}
	for _, s := range steps {
		var got int
		run(t, s.name, func() {
			s.call()
			got = c.Value()
		})
		if got != s.want {
			t.Errorf("after %s: Value() = %d, want %d", s.name, got, s.want)
		}
	}
}

func TestConcurrentInc(t *testing.T) {
	const goroutines, perG = 50, 1000
	var c Counter
	run(t, "50 goroutines × 1000 Inc()", func() {
		var wg sync.WaitGroup
		for range goroutines {
			wg.Go(func() {
				for range perG {
					c.Inc()
				}
			})
		}
		wg.Wait()
	})
	if got := c.Value(); got != goroutines*perG {
		t.Errorf("%d goroutines × %d Inc(): Value() = %d, want %d", goroutines, perG, got, goroutines*perG)
	}
}

func TestConcurrentMixed(t *testing.T) {
	var c Counter
	run(t, "40 goroutines × 1000 Add()", func() {
		var wg sync.WaitGroup
		for g := range 40 {
			wg.Go(func() {
				for range 1000 {
					if g%2 == 0 {
						c.Add(3)
					} else {
						c.Add(-1)
					}
					_ = c.Value()
				}
			})
		}
		wg.Wait()
	})
	want := 20*1000*3 - 20*1000
	if got := c.Value(); got != want {
		t.Errorf("20 goroutines × 1000 Add(3) and 20 × 1000 Add(-1): Value() = %d, want %d", got, want)
	}
}

func TestMethodsTakeTheLock(t *testing.T) {
	var c Counter
	tests := []struct {
		name string
		call func()
	}{
		{"Inc()", c.Inc},
		{"Add(5)", func() { c.Add(5) }},
		{"Value()", func() { _ = c.Value() }},
	}
	for _, tt := range tests {
		c.mu.Lock()
		done := make(chan struct{})
		go func() {
			defer close(done)
			tt.call()
		}()
		time.Sleep(20 * time.Millisecond)
		select {
		case <-done:
			c.mu.Unlock()
			t.Errorf("%s while c.mu is locked = returned, want it to wait for the lock", tt.name)
			continue
		default:
		}
		c.mu.Unlock()
		wait(t, done, tt.name+" after c.mu.Unlock()")
	}
	if c.n != 6 {
		t.Errorf("after Inc() and Add(5): c.n = %d, want 6", c.n)
	}
}

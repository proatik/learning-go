package cache

import (
	"fmt"
	"sync"
	"testing"
	"time"
)

// start runs f in a goroutine and returns a channel closed when f returns.
func start(f func()) <-chan struct{} {
	done := make(chan struct{})
	go func() {
		defer close(done)
		f()
	}()
	return done
}

// finishes fails the test if done isn't closed within 2s.
func finishes(t *testing.T, done <-chan struct{}, what string) {
	t.Helper()
	select {
	case <-done:
	case <-time.After(2 * time.Second):
		t.Fatalf("%s did not finish within 2s (deadlock?)", what)
	}
}

// waits reports whether done is still open after 20ms.
func waits(done <-chan struct{}) bool {
	time.Sleep(20 * time.Millisecond)
	select {
	case <-done:
		return false
	default:
		return true
	}
}

func TestBasics(t *testing.T) {
	c := New()
	if c == nil || c.items == nil {
		t.Fatalf("New() = %+v, want a cache with a non-nil items map", c)
	}
	finishes(t, start(func() {
		if got := c.Len(); got != 0 {
			t.Errorf("New().Len() = %d, want 0", got)
		}
		if v, ok := c.Get("a"); v != "" || ok {
			t.Errorf("Get(%q) on empty cache = (%q, %v), want (\"\", false)", "a", v, ok)
		}
		c.Set("a", "1")
		c.Set("b", "2")
		c.Set("", "empty key")
		c.Set("ключ", "")
		if v, ok := c.Get("a"); v != "1" || !ok {
			t.Errorf("Get(%q) = (%q, %v), want (%q, true)", "a", v, ok, "1")
		}
		if v, ok := c.Get(""); v != "empty key" || !ok {
			t.Errorf("Get(%q) = (%q, %v), want (%q, true)", "", v, ok, "empty key")
		}
		if v, ok := c.Get("ключ"); v != "" || !ok {
			t.Errorf("Get(%q) = (%q, %v), want (\"\", true)", "ключ", v, ok)
		}
		if got := c.Len(); got != 4 {
			t.Errorf("Len() after 4 Set = %d, want 4", got)
		}
		c.Set("a", "one")
		if v, _ := c.Get("a"); v != "one" {
			t.Errorf("Get(%q) after overwrite = %q, want %q", "a", v, "one")
		}
		if got := c.Len(); got != 4 {
			t.Errorf("Len() after overwriting a key = %d, want 4", got)
		}
		c.Delete("a")
		c.Delete("missing")
		if v, ok := c.Get("a"); v != "" || ok {
			t.Errorf("Get(%q) after Delete = (%q, %v), want (\"\", false)", "a", v, ok)
		}
		if got := c.Len(); got != 3 {
			t.Errorf("Len() after Delete = %d, want 3", got)
		}
	}), "Get/Set/Delete/Len")
}

func TestReadersShareTheLock(t *testing.T) {
	c := &Cache{items: map[string]string{"a": "1"}}
	c.mu.RLock()
	defer c.mu.RUnlock()
	var v string
	finishes(t, start(func() { v, _ = c.Get("a") }), "Get while another reader holds c.mu.RLock()")
	if v != "1" {
		t.Errorf("Get(%q) = %q, want %q", "a", v, "1")
	}
	n := 0
	finishes(t, start(func() { n = c.Len() }), "Len while another reader holds c.mu.RLock()")
	if n != 1 {
		t.Errorf("Len() = %d, want 1", n)
	}
}

func TestWritersWaitForReaders(t *testing.T) {
	for _, op := range []string{"Set", "Delete"} {
		c := &Cache{items: map[string]string{"a": "1"}}
		c.mu.RLock()
		done := start(func() {
			if op == "Set" {
				c.Set("b", "2")
			} else {
				c.Delete("a")
			}
		})
		if !waits(done) {
			c.mu.RUnlock()
			t.Errorf("%s while a reader holds c.mu.RLock() = returned, want it to wait (use c.mu.Lock())", op)
			continue
		}
		c.mu.RUnlock()
		finishes(t, done, op+" after the reader unlocked")
	}
}

func TestReadersWaitForWriters(t *testing.T) {
	for _, op := range []string{"Get", "Len"} {
		c := &Cache{items: map[string]string{"a": "1"}}
		c.mu.Lock()
		done := start(func() {
			if op == "Get" {
				c.Get("a")
			} else {
				c.Len()
			}
		})
		if !waits(done) {
			c.mu.Unlock()
			t.Errorf("%s while a writer holds c.mu.Lock() = returned, want it to wait (use c.mu.RLock())", op)
			continue
		}
		c.mu.Unlock()
		finishes(t, done, op+" after the writer unlocked")
	}
}

func TestConcurrent(t *testing.T) {
	c := New()
	if c == nil || c.items == nil {
		t.Fatalf("New() = %+v, want a cache with a non-nil items map", c)
	}
	var wg sync.WaitGroup
	for g := range 20 {
		wg.Go(func() { // writer: its own 100 keys
			for i := range 100 {
				c.Set(fmt.Sprintf("g%d-%d", g, i), fmt.Sprint(i))
			}
		})
		wg.Go(func() { // reader
			for i := range 500 {
				if v, ok := c.Get(fmt.Sprintf("g%d-%d", g, i%100)); ok && v != fmt.Sprint(i%100) {
					t.Errorf("Get(%q) = %q, want %q", fmt.Sprintf("g%d-%d", g, i%100), v, fmt.Sprint(i%100))
					return
				}
				_ = c.Len()
			}
		})
	}
	finishes(t, start(wg.Wait), "20 writers and 20 readers")
	if got := c.Len(); got != 2000 {
		t.Errorf("Len() after 20 goroutines × 100 Set = %d, want 2000", got)
	}
	for g := range 20 {
		wg.Go(func() {
			for i := range 50 {
				c.Delete(fmt.Sprintf("g%d-%d", g, i))
			}
		})
	}
	finishes(t, start(wg.Wait), "20 goroutines × 50 Delete")
	if got := c.Len(); got != 1000 {
		t.Errorf("Len() after deleting 1000 of 2000 keys = %d, want 1000", got)
	}
}

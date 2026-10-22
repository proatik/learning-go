package first

import (
	"fmt"
	"testing"
	"testing/synctest"
	"time"
)

const (
	isNil  = -1 // the channel is nil
	never  = -2 // the channel exists but nothing is ever sent
	before = -3 // the value is already buffered before FirstOf is called
)

// source returns a buffered channel that receives v at the given fake time,
// or a nil / silent / pre-filled channel for the special values above.
func source(v string, at time.Duration) chan string {
	switch at {
	case isNil:
		return nil
	case never:
		return make(chan string, 1)
	case before:
		ch := make(chan string, 1)
		ch <- v
		return ch
	}
	ch := make(chan string, 1) // buffered: the sender never blocks
	go func() {
		time.Sleep(at)
		ch <- v
	}()
	return ch
}

func describe(at time.Duration) string {
	switch at {
	case isNil:
		return "nil"
	case never:
		return "silent"
	case before:
		return "already holds a value"
	}
	return fmt.Sprintf("sends at %v", at)
}

func TestFirstOf(t *testing.T) {
	tests := []struct {
		aAt, bAt, timeout time.Duration
		want              string
		wantFrom          int
		wantAt            time.Duration
	}{
		{time.Second, 2 * time.Second, 5 * time.Second, "x", 1, time.Second},
		{2 * time.Second, time.Second, 5 * time.Second, "y", 2, time.Second},
		{3 * time.Second, 4 * time.Second, 2 * time.Second, "", 0, 2 * time.Second},
		{before, 2 * time.Second, 5 * time.Second, "x", 1, 0},
		{3 * time.Second, before, 5 * time.Second, "y", 2, 0},
		{isNil, time.Second, 5 * time.Second, "y", 2, time.Second},
		{1500 * time.Millisecond, isNil, 5 * time.Second, "x", 1, 1500 * time.Millisecond},
		{isNil, isNil, time.Second, "", 0, time.Second},
		{never, never, 300 * time.Millisecond, "", 0, 300 * time.Millisecond},
		{never, 4 * time.Second, 5 * time.Second, "y", 2, 4 * time.Second},
	}
	for _, tt := range tests {
		synctest.Test(t, func(t *testing.T) {
			a, b := source("x", tt.aAt), source("y", tt.bAt)
			call := fmt.Sprintf("FirstOf(a %s, b %s, %v)", describe(tt.aAt), describe(tt.bAt), tt.timeout)

			start := time.Now()
			got, from := FirstOf(a, b, tt.timeout)
			at := time.Since(start)

			if got != tt.want || from != tt.wantFrom {
				t.Errorf("%s = (%q, %d), want (%q, %d)", call, got, from, tt.want, tt.wantFrom)
			}
			if at != tt.wantAt {
				t.Errorf("%s returned after %v, want %v", call, at, tt.wantAt)
			}
			time.Sleep(time.Hour) // let every sender finish
		})
	}
}

func TestFirstOfReceivesOnlyOne(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		a, b := source("x", time.Second), source("y", 2*time.Second)
		FirstOf(a, b, time.Minute)
		time.Sleep(time.Hour)
		if len(a) != 0 || len(b) != 1 {
			t.Errorf("FirstOf(a sends at 1s, b sends at 2s, 1m): values left in a, b = %d, %d, want 0, 1", len(a), len(b))
		}
	})
}

func TestFirstOfInts(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		a := make(chan int, 1)
		b := make(chan int, 1)
		b <- 0 // a zero value is still a real value
		got, from := FirstOf(a, b, time.Second)
		if got != 0 || from != 2 {
			t.Errorf("FirstOf(a silent, b holds 0, 1s) = (%d, %d), want (0, 2)", got, from)
		}
	})
}

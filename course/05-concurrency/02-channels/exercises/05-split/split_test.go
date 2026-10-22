package split

import (
	"fmt"
	"slices"
	"testing"
	"time"
)

// closed returns a buffered channel holding vals, already closed.
func closed(vals []int) <-chan int {
	ch := make(chan int, len(vals))
	for _, v := range vals {
		ch <- v
	}
	close(ch)
	return ch
}

// split calls Split(in) and reads both outputs concurrently until they are
// closed, failing if anything takes over 2s.
func split(t *testing.T, call string, in <-chan int) (evens, odds []int) {
	t.Helper()
	timeout := time.After(2 * time.Second)
	type pair struct{ e, o <-chan int }
	started := make(chan pair, 1)
	go func() {
		e, o := Split(in)
		started <- pair{e, o}
	}()
	var p pair
	select {
	case p = <-started:
	case <-timeout:
		t.Fatalf("%s did not return within 2s, want it to return at once and work in a goroutine", call)
	}
	if p.e == nil || p.o == nil {
		t.Fatalf("%s returned a nil channel (evens nil: %v, odds nil: %v)", call, p.e == nil, p.o == nil)
	}
	read := func(ch <-chan int, out chan<- []int) {
		var got []int
		for v := range ch {
			got = append(got, v)
		}
		out <- got
	}
	evCh, odCh := make(chan []int, 1), make(chan []int, 1)
	go read(p.e, evCh)
	go read(p.o, odCh)
	evDone, odDone := false, false
	for !evDone || !odDone {
		select {
		case evens = <-evCh:
			evDone = true
		case odds = <-odCh:
			odDone = true
		case <-timeout:
			t.Fatalf("%s: evens closed: %v, odds closed: %v within 2s, want both closed", call, evDone, odDone)
		}
	}
	return evens, odds
}

func TestSplit(t *testing.T) {
	tests := []struct {
		in          []int
		evens, odds []int
	}{
		{[]int{1, 2, 3, 4, 5}, []int{2, 4}, []int{1, 3, 5}},
		{[]int{-3, -2, 0, 7}, []int{-2, 0}, []int{-3, 7}},
		{[]int{2, 4}, []int{2, 4}, nil},
		{[]int{-1, -5, 9}, nil, []int{-1, -5, 9}},
		{[]int{0}, []int{0}, nil},
		{nil, nil, nil},
	}
	for _, tt := range tests {
		call := fmt.Sprintf("Split(%v)", tt.in)
		ev, od := split(t, call, closed(tt.in))
		if !slices.Equal(ev, tt.evens) || !slices.Equal(od, tt.odds) {
			t.Errorf("%s = evens %v, odds %v, want evens %v, odds %v", call, ev, od, tt.evens, tt.odds)
		}
	}
}

func TestSplitStream(t *testing.T) {
	in := make(chan int) // unbuffered: values arrive one at a time
	go func() {
		defer close(in)
		for i := range 1000 {
			in <- i - 500
		}
	}()
	ev, od := split(t, "Split(-500..499)", in)
	if len(ev) != 500 || len(od) != 500 {
		t.Fatalf("Split(-500..499) = %d evens, %d odds, want 500 and 500", len(ev), len(od))
	}
	for i := range 500 {
		if ev[i] != 2*i-500 || od[i] != 2*i-499 {
			t.Fatalf("Split(-500..499): evens[%d] = %d, odds[%d] = %d, want %d and %d (input order)", i, ev[i], i, od[i], 2*i-500, 2*i-499)
		}
	}
}

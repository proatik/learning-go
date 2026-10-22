package batch

import (
	"fmt"
	"slices"
	"strconv"
	"testing"
	"time"
)

// closed returns a buffered channel holding vals, already closed.
func closed(vals []string) <-chan string {
	ch := make(chan string, len(vals))
	for _, v := range vals {
		ch <- v
	}
	close(ch)
	return ch
}

// batches calls Batch(in, size) and keeps every batch it receives until the
// output is closed, failing if anything takes over 2s.
func batches(t *testing.T, call string, in <-chan string, size int) [][]string {
	t.Helper()
	timeout := time.After(2 * time.Second)
	started := make(chan (<-chan []string), 1)
	go func() { started <- Batch(in, size) }()
	var ch <-chan []string
	select {
	case ch = <-started:
	case <-timeout:
		t.Fatalf("%s did not return within 2s, want it to return at once and work in a goroutine", call)
	}
	if ch == nil {
		t.Fatalf("%s returned a nil channel", call)
	}
	var got [][]string
	for {
		select {
		case b, ok := <-ch:
			if !ok {
				return got
			}
			got = append(got, b)
			if len(got) > 10000 {
				t.Fatalf("%s sent over 10000 batches, want it to stop", call)
			}
		case <-timeout:
			t.Fatalf("%s: output not closed within 2s (received %q so far)", call, got)
		}
	}
}

func TestBatch(t *testing.T) {
	tests := []struct {
		in   []string
		size int
		want [][]string
	}{
		{[]string{"a", "b", "c", "d", "e"}, 2, [][]string{{"a", "b"}, {"c", "d"}, {"e"}}},
		{[]string{"a", "b", "c", "d"}, 2, [][]string{{"a", "b"}, {"c", "d"}}},
		{[]string{"a", "b"}, 5, [][]string{{"a", "b"}}},
		{[]string{"a", "b", "c"}, 1, [][]string{{"a"}, {"b"}, {"c"}}},
		{[]string{"a", "b", "c"}, 3, [][]string{{"a", "b", "c"}}},
		{[]string{"", "é", "世界"}, 2, [][]string{{"", "é"}, {"世界"}}},
		{nil, 3, nil},
	}
	for _, tt := range tests {
		call := fmt.Sprintf("Batch(%q, %d)", tt.in, tt.size)
		got := batches(t, call, closed(tt.in), tt.size)
		if !slices.EqualFunc(got, tt.want, slices.Equal) {
			t.Errorf("%s = %q, want %q", call, got, tt.want)
		}
	}
}

func TestBatchesAreIndependent(t *testing.T) {
	in := make(chan string) // unbuffered: one value at a time
	go func() {
		defer close(in)
		for i := range 1003 {
			in <- strconv.Itoa(i)
		}
	}()
	got := batches(t, "Batch(0..1002, 10)", in, 10)
	if len(got) != 101 {
		t.Fatalf("Batch(0..1002, 10) = %d batches, want 101", len(got))
	}
	for i, b := range got {
		want := 10
		if i == 100 {
			want = 3
		}
		if len(b) != want {
			t.Fatalf("Batch(0..1002, 10): batch %d has %d values, want %d", i, len(b), want)
		}
		for j, v := range b {
			if w := strconv.Itoa(i*10 + j); v != w {
				t.Fatalf("Batch(0..1002, 10): batch %d = %q, want it to start at %q (was the slice reused after sending?)", i, b, strconv.Itoa(i*10))
			}
		}
	}
}

func TestBatchPanics(t *testing.T) {
	for _, size := range []int{0, -1} {
		func() {
			defer func() {
				r := recover()
				if r != "batch: size must be at least 1" {
					t.Errorf("Batch(in, %d) panic = %v, want %q", size, r, "batch: size must be at least 1")
				}
			}()
			Batch(closed(nil), size)
		}()
	}
}

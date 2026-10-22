package batch

import (
	"fmt"
	"slices"
	"strings"
	"testing"
	"testing/synctest"
	"time"
)

const ms = time.Millisecond

type item struct {
	v  int
	at time.Duration
}

type sent struct {
	batch []int
	at    time.Duration
}

// feed sends each item at its fake time on a buffered channel (so the feeder
// never blocks), then closes the channel at closeAt.
func feed(items []item, closeAt time.Duration) chan int {
	ch := make(chan int, len(items))
	go func() {
		start := time.Now()
		for _, it := range items {
			time.Sleep(it.at - time.Since(start))
			ch <- it.v
		}
		time.Sleep(closeAt - time.Since(start))
		close(ch)
	}()
	return ch
}

func listItems(items []item) string {
	var parts []string
	for _, it := range items {
		parts = append(parts, fmt.Sprintf("%d@%v", it.v, it.at))
	}
	return "[" + strings.Join(parts, " ") + "]"
}

func listSent(s []sent) string {
	var parts []string
	for _, b := range s {
		parts = append(parts, fmt.Sprintf("%v@%v", b.batch, b.at))
	}
	return "[" + strings.Join(parts, " ") + "]"
}

func equalSent(a, b []sent) bool {
	return slices.EqualFunc(a, b, func(x, y sent) bool {
		return x.at == y.at && slices.Equal(x.batch, y.batch)
	})
}

// drain records every batch received from out with its fake time, and the
// time out was closed. A silent output fails after one (fake) hour.
func drain(t *testing.T, call string, out <-chan []int, start time.Time) (got []sent, closedAt time.Duration) {
	t.Helper()
	for len(got) <= 100 {
		select {
		case b, ok := <-out:
			if !ok {
				return got, time.Since(start)
			}
			got = append(got, sent{b, time.Since(start)})
		case <-time.After(time.Hour):
			t.Fatalf("%s: output sent %s and then nothing for an hour, want it closed", call, listSent(got))
		}
	}
	t.Fatalf("%s: output sent more than 100 batches: %s", call, listSent(got[:10]))
	return nil, 0
}

func TestBatch(t *testing.T) {
	tests := []struct {
		size       int
		wait       time.Duration
		in         []item
		closeAt    time.Duration
		want       []sent
		wantClosed time.Duration
	}{
		{3, time.Second, []item{{1, 0}, {2, 0}, {3, 0}, {4, 0}, {5, 500 * ms}}, 3 * time.Second,
			[]sent{{[]int{1, 2, 3}, 0}, {[]int{4, 5}, time.Second}}, 3 * time.Second},
		{3, time.Second, []item{{1, 0}, {2, 200 * ms}}, 400 * ms,
			[]sent{{[]int{1, 2}, 400 * ms}}, 400 * ms},
		{3, time.Second, []item{{1, 0}, {2, 100 * ms}, {3, 200 * ms}, {4, 1100 * ms}}, 5 * time.Second,
			[]sent{{[]int{1, 2, 3}, 200 * ms}, {[]int{4}, 2100 * ms}}, 5 * time.Second},
		{3, time.Second, nil, 100 * ms, nil, 100 * ms},
		{3, time.Second, nil, 0, nil, 0},
		{3, time.Second, []item{{1, 0}, {2, 1500 * ms}, {3, 1700 * ms}, {4, 1800 * ms}, {5, 1900 * ms}}, 10 * time.Second,
			[]sent{{[]int{1}, time.Second}, {[]int{2, 3, 4}, 1800 * ms}, {[]int{5}, 2900 * ms}}, 10 * time.Second},
		{3, time.Second, []item{{1, 0}, {2, 999 * ms}, {3, 1001 * ms}}, 1500 * ms,
			[]sent{{[]int{1, 2}, time.Second}, {[]int{3}, 1500 * ms}}, 1500 * ms},
		{1, time.Second, []item{{1, 0}, {2, 10 * ms}, {3, 10 * ms}}, time.Minute,
			[]sent{{[]int{1}, 0}, {[]int{2}, 10 * ms}, {[]int{3}, 10 * ms}}, time.Minute},
		{0, time.Second, []item{{7, 0}, {8, 5 * ms}}, time.Minute,
			[]sent{{[]int{7}, 0}, {[]int{8}, 5 * ms}}, time.Minute},
		{2, 50 * ms, []item{{1, 0}, {2, 0}, {3, 0}, {4, 0}, {5, 0}}, time.Second,
			[]sent{{[]int{1, 2}, 0}, {[]int{3, 4}, 0}, {[]int{5}, 50 * ms}}, time.Second},
	}
	for _, tt := range tests {
		synctest.Test(t, func(t *testing.T) {
			call := fmt.Sprintf("Batch(%s closed@%v, %d, %v)", listItems(tt.in), tt.closeAt, tt.size, tt.wait)
			in := feed(tt.in, tt.closeAt)
			start := time.Now()
			out := Batch(in, tt.size, tt.wait)
			if out == nil {
				time.Sleep(time.Hour) // let the feeder finish
				t.Fatalf("%s returned a nil channel", call)
			}
			got, closedAt := drain(t, call, out, start)
			// Compare only after every batch has been received, so a reused
			// backing array shows up as overwritten earlier batches.
			if !equalSent(got, tt.want) {
				t.Errorf("%s sent %s, want %s", call, listSent(got), listSent(tt.want))
			}
			if closedAt != tt.wantClosed {
				t.Errorf("%s closed the output at %v, want %v", call, closedAt, tt.wantClosed)
			}
		})
	}
}

func TestBatchStrings(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		in := make(chan string, 4)
		for _, s := range []string{"a", "b", "c", "d"} {
			in <- s
		}
		close(in)
		out := Batch(in, 3, time.Second)
		if out == nil {
			t.Fatalf("Batch(a b c d closed, 3, 1s) returned a nil channel")
		}
		var got [][]string
		for b := range out { // safe: the input is already closed, so a correct Batch closes out at once
			got = append(got, b)
		}
		want := [][]string{{"a", "b", "c"}, {"d"}}
		if !slices.EqualFunc(got, want, slices.Equal) {
			t.Errorf("Batch(a b c d closed, 3, 1s) sent %q, want %q", got, want)
		}
	})
}

package debounce

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

func (it item) String() string { return fmt.Sprintf("%d@%v", it.v, it.at) }

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

func list(items []item) string {
	var parts []string
	for _, it := range items {
		parts = append(parts, it.String())
	}
	return "[" + strings.Join(parts, " ") + "]"
}

// drain records every value received from out with its fake time, and the
// time out was closed. A silent output fails after one (fake) hour.
func drain(t *testing.T, call string, out <-chan int, start time.Time) (got []item, closedAt time.Duration) {
	t.Helper()
	for len(got) <= 100 {
		select {
		case v, ok := <-out:
			if !ok {
				return got, time.Since(start)
			}
			got = append(got, item{v, time.Since(start)})
		case <-time.After(time.Hour):
			t.Fatalf("%s: output sent %v and then nothing for an hour, want it closed", call, list(got))
		}
	}
	t.Fatalf("%s: output sent more than 100 values: %v", call, list(got[:10]))
	return nil, 0
}

func TestDebounce(t *testing.T) {
	tests := []struct {
		d          time.Duration
		in         []item
		closeAt    time.Duration
		want       []item
		wantClosed time.Duration
	}{
		{200 * ms, []item{{1, 0}, {2, 50 * ms}, {3, 100 * ms}, {4, 400 * ms}, {5, 450 * ms}}, time.Second,
			[]item{{3, 300 * ms}, {5, 650 * ms}}, time.Second},
		{200 * ms, []item{{1, 0}, {2, 300 * ms}, {3, 600 * ms}}, 2 * time.Second,
			[]item{{1, 200 * ms}, {2, 500 * ms}, {3, 800 * ms}}, 2 * time.Second},
		{200 * ms, []item{{1, 0}}, 100 * ms,
			[]item{{1, 100 * ms}}, 100 * ms},
		{200 * ms, nil, 500 * ms, nil, 500 * ms},
		{200 * ms, nil, 0, nil, 0},
		{200 * ms, []item{{1, 0}, {2, 150 * ms}, {3, 300 * ms}, {4, 450 * ms}, {5, 600 * ms}, {6, 750 * ms}, {7, 900 * ms}}, 2 * time.Second,
			[]item{{7, 1100 * ms}}, 2 * time.Second},
		{200 * ms, []item{{1, 0}, {2, 500 * ms}, {3, 550 * ms}}, 600 * ms,
			[]item{{1, 200 * ms}, {3, 600 * ms}}, 600 * ms},
		{time.Second, []item{{0, 0}, {0, 2 * time.Second}}, time.Minute,
			[]item{{0, time.Second}, {0, 3 * time.Second}}, time.Minute},
	}
	for _, tt := range tests {
		synctest.Test(t, func(t *testing.T) {
			call := fmt.Sprintf("Debounce(%s closed@%v, %v)", list(tt.in), tt.closeAt, tt.d)
			in := feed(tt.in, tt.closeAt)
			start := time.Now()
			out := Debounce(in, tt.d)
			if out == nil {
				time.Sleep(time.Hour) // let the feeder finish
				t.Fatalf("%s returned a nil channel", call)
			}
			got, closedAt := drain(t, call, out, start)
			if !slices.Equal(got, tt.want) {
				t.Errorf("%s sent %v, want %v", call, list(got), list(tt.want))
			}
			if closedAt != tt.wantClosed {
				t.Errorf("%s closed the output at %v, want %v", call, closedAt, tt.wantClosed)
			}
		})
	}
}

func TestDebounceStrings(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		in := make(chan string, 3)
		in <- "g"
		in <- "go"
		in <- "gopher"
		close(in)
		out := Debounce(in, time.Second)
		if out == nil {
			t.Fatalf("Debounce(g go gopher closed, 1s) returned a nil channel")
		}
		var got []string
		for v := range out { // safe: the input is already closed, so a correct Debounce closes out at once
			got = append(got, v)
		}
		if want := []string{"gopher"}; !slices.Equal(got, want) {
			t.Errorf("Debounce(g go gopher closed, 1s) sent %q, want %q", got, want)
		}
	})
}

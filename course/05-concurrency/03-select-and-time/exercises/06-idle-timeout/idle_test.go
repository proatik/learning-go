package idle

import (
	"fmt"
	"slices"
	"strings"
	"testing"
	"testing/synctest"
	"time"
)

const ms = time.Millisecond

type event struct {
	name string
	at   time.Duration
}

// feed sends each event at its fake time on a buffered channel (so the
// feeder never blocks), then closes the channel at closeAt.
func feed(events []event, closeAt time.Duration) <-chan string {
	ch := make(chan string, len(events))
	go func() {
		start := time.Now()
		for _, e := range events {
			time.Sleep(e.at - time.Since(start))
			ch <- e.name
		}
		time.Sleep(closeAt - time.Since(start))
		close(ch)
	}()
	return ch
}

func describe(events []event, closeAt time.Duration) string {
	var parts []string
	for _, e := range events {
		parts = append(parts, fmt.Sprintf("%s@%v", e.name, e.at))
	}
	return fmt.Sprintf("events [%s] closed@%v", strings.Join(parts, " "), closeAt)
}

func TestUntilIdle(t *testing.T) {
	tests := []struct {
		d       time.Duration
		events  []event
		closeAt time.Duration
		want    []string
		wantAt  time.Duration
	}{
		{250 * ms, []event{{"a", 100 * ms}, {"b", 300 * ms}, {"c", 600 * ms}}, 2 * time.Second, []string{"a", "b"}, 550 * ms},
		{250 * ms, []event{{"a", 200 * ms}, {"b", 400 * ms}, {"c", 600 * ms}}, 5 * time.Second, []string{"a", "b", "c"}, 850 * ms},
		{250 * ms, []event{{"a", 100 * ms}, {"b", 300 * ms}, {"c", 500 * ms}, {"d", 700 * ms}}, 800 * ms, []string{"a", "b", "c", "d"}, 800 * ms},
		{250 * ms, nil, time.Second, nil, 250 * ms},
		{250 * ms, nil, 0, nil, 0},
		{250 * ms, []event{{"a", 300 * ms}}, time.Second, nil, 250 * ms},
		{time.Second, []event{{"x", 0}, {"y", 0}, {"z", 999 * ms}}, time.Minute, []string{"x", "y", "z"}, 1999 * ms},
		{100 * ms, []event{{"a", 50 * ms}, {"b", 100 * ms}, {"c", 150 * ms}, {"d", 200 * ms}, {"e", 250 * ms}, {"f", 351 * ms}}, time.Minute, []string{"a", "b", "c", "d", "e"}, 350 * ms},
	}
	for _, tt := range tests {
		synctest.Test(t, func(t *testing.T) {
			events := feed(tt.events, tt.closeAt)
			call := fmt.Sprintf("UntilIdle(%s, %v)", describe(tt.events, tt.closeAt), tt.d)

			start := time.Now()
			got := UntilIdle(events, tt.d)
			at := time.Since(start)

			if !slices.Equal(got, tt.want) {
				t.Errorf("%s = %q, want %q", call, got, tt.want)
			}
			if at != tt.wantAt {
				t.Errorf("%s returned after %v, want %v", call, at, tt.wantAt)
			}
			time.Sleep(time.Hour) // let the feeder finish
		})
	}
}

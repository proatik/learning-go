package timeout

import (
	"errors"
	"testing"
	"testing/synctest"
	"time"
)

func TestDo(t *testing.T) {
	tests := []struct {
		d       time.Duration
		work    time.Duration // how long f takes
		want    int
		wantErr error
		wantAt  time.Duration // when Do should return
	}{
		{d: 2 * time.Second, work: time.Second, want: 42, wantAt: time.Second},
		{d: 2 * time.Second, work: 3 * time.Second, wantErr: ErrTimeout, wantAt: 2 * time.Second},
		{d: time.Second, work: 0, want: 42, wantAt: 0},
		{d: 500 * time.Millisecond, work: 499 * time.Millisecond, want: 42, wantAt: 499 * time.Millisecond},
		{d: 500 * time.Millisecond, work: 501 * time.Millisecond, wantErr: ErrTimeout, wantAt: 500 * time.Millisecond},
	}
	for _, tt := range tests {
		synctest.Test(t, func(t *testing.T) {
			calls := 0
			f := func() int {
				calls++
				time.Sleep(tt.work)
				return 42
			}
			start := time.Now()
			got, err := Do(tt.d, f)
			at := time.Since(start)

			if got != tt.want || !errors.Is(err, tt.wantErr) {
				t.Errorf("Do(%v, f taking %v) = (%d, %v), want (%d, %v)", tt.d, tt.work, got, err, tt.want, tt.wantErr)
			}
			if at != tt.wantAt {
				t.Errorf("Do(%v, f taking %v) returned after %v, want %v", tt.d, tt.work, at, tt.wantAt)
			}
			time.Sleep(time.Hour) // let f finish
			synctest.Wait()
			if calls != 1 {
				t.Errorf("Do(%v, f taking %v) called f %d times, want 1", tt.d, tt.work, calls)
			}
		})
	}
}

func TestDoNonPositive(t *testing.T) {
	for _, d := range []time.Duration{0, -time.Second} {
		synctest.Test(t, func(t *testing.T) {
			calls := 0
			start := time.Now()
			got, err := Do(d, func() string { calls++; return "hi" })
			at := time.Since(start)
			time.Sleep(time.Hour)
			synctest.Wait()
			if got != "" || !errors.Is(err, ErrTimeout) {
				t.Errorf("Do(%v, f) = (%q, %v), want (\"\", %v)", d, got, err, ErrTimeout)
			}
			if at != 0 || calls != 0 {
				t.Errorf("Do(%v, f) took %v and called f %d times, want 0s and 0 calls", d, at, calls)
			}
		})
	}
}

func TestDoString(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		got, err := Do(time.Minute, func() string {
			time.Sleep(10 * time.Second)
			return "done"
		})
		if got != "done" || err != nil {
			t.Errorf("Do(1m, f returning \"done\" after 10s) = (%q, %v), want (\"done\", <nil>)", got, err)
		}
		time.Sleep(time.Hour)
	})
}

package gen

import (
	"context"
	"testing"
	"testing/synctest"
)

// recv takes the next value once every goroutine has caught up. It fails if
// no value is ready or the channel is closed.
func recv(t *testing.T, name string, ch <-chan int) int {
	t.Helper()
	synctest.Wait()
	select {
	case v, ok := <-ch:
		if !ok {
			t.Fatalf("%s: channel closed early", name)
		}
		return v
	default:
		t.Fatalf("%s: no value ready (is a goroutine sending?)", name)
	}
	return 0
}

// mustClose checks that ch is closed after ctx was canceled, allowing one
// value that was already on offer.
func mustClose(t *testing.T, name string, ch <-chan int) {
	t.Helper()
	for i := range 2 {
		synctest.Wait()
		select {
		case v, ok := <-ch:
			if !ok {
				return
			}
			if i == 1 {
				t.Fatalf("%s: still sending after cancel (got %d)", name, v)
			}
		default:
			t.Fatalf("%s: channel not closed after cancel (the goroutine is stuck: watch ctx.Done())", name)
		}
	}
}

func notNil(t *testing.T, name string, ch <-chan int) {
	t.Helper()
	if ch == nil {
		t.Fatalf("%s returned a nil channel", name)
	}
	if c := cap(ch); c != 0 {
		t.Fatalf("%s returned a channel with capacity %d, want 0 (unbuffered)", name, c)
	}
}

func TestCount(t *testing.T) {
	tests := []struct {
		name        string
		start, step int
		want        []int
	}{
		{"Count(ctx, 1, 1)", 1, 1, []int{1, 2, 3, 4, 5}},
		{"Count(ctx, 10, -5)", 10, -5, []int{10, 5, 0, -5, -10}},
		{"Count(ctx, 7, 0)", 7, 0, []int{7, 7, 7}},
		{"Count(ctx, -3, 100)", -3, 100, []int{-3, 97, 197}},
	}
	for _, tt := range tests {
		synctest.Test(t, func(t *testing.T) {
			ctx, cancel := context.WithCancel(t.Context())
			defer cancel()
			ch := Count(ctx, tt.start, tt.step)
			notNil(t, tt.name, ch)
			for i, want := range tt.want {
				if got := recv(t, tt.name, ch); got != want {
					t.Fatalf("%s value %d = %d, want %d", tt.name, i, got, want)
				}
			}
			cancel()
			mustClose(t, tt.name, ch)
		})
	}
}

func TestCountCancelUnread(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		ctx, cancel := context.WithCancel(t.Context())
		ch := Count(ctx, 0, 1)
		notNil(t, "Count(ctx, 0, 1)", ch)
		synctest.Wait()
		cancel()
		mustClose(t, "Count(ctx, 0, 1) never read, then canceled", ch)
	})
}

func TestCountAlreadyCanceled(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		ctx, cancel := context.WithCancel(t.Context())
		cancel()
		ch := Count(ctx, 0, 1)
		notNil(t, "Count(canceled ctx, 0, 1)", ch)
		mustClose(t, "Count(canceled ctx, 0, 1)", ch)
	})
}

func TestMapOverCount(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		ctx, cancel := context.WithCancel(t.Context())
		defer cancel()
		nums := Count(ctx, 1, 1)
		notNil(t, "Count(ctx, 1, 1)", nums)
		sq := Map(ctx, nums, func(n int) int { return n * n })
		notNil(t, "Map(ctx, Count(ctx, 1, 1), square)", sq)
		for i, want := range []int{1, 4, 9, 16, 25} {
			if got := recv(t, "Map(ctx, Count(ctx, 1, 1), square)", sq); got != want {
				t.Fatalf("Map(ctx, Count(ctx, 1, 1), square) value %d = %d, want %d", i, got, want)
			}
		}
		cancel()
		mustClose(t, "Map(ctx, Count(…), square) after cancel", sq)
		mustClose(t, "Count(ctx, 1, 1) after cancel", nums)
	})
}

func TestMapInputClosed(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		in := make(chan int, 2)
		in <- 3
		in <- 4
		close(in)
		out := Map(t.Context(), in, func(n int) int { return 2 * n })
		notNil(t, "Map(ctx, [3 4], double)", out)
		for i, want := range []int{6, 8} {
			if got := recv(t, "Map(ctx, [3 4], double)", out); got != want {
				t.Fatalf("Map(ctx, [3 4], double) value %d = %d, want %d", i, got, want)
			}
		}
		synctest.Wait()
		select {
		case v, ok := <-out:
			if ok {
				t.Fatalf("Map(ctx, [3 4], double) sent extra value %d, want closed", v)
			}
		default:
			t.Fatalf("Map(ctx, [3 4], double): channel not closed after in was closed")
		}
	})
}

func TestMapCancelWhileReceiving(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		ctx, cancel := context.WithCancel(t.Context())
		in := make(chan int) // never sends, never closed
		out := Map(ctx, in, func(n int) int { return n })
		notNil(t, "Map(ctx, silent in, id)", out)
		synctest.Wait()
		cancel()
		mustClose(t, "Map(ctx, silent in, id) after cancel", out)
	})
}

func TestMapCancelWhileSending(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		ctx, cancel := context.WithCancel(t.Context())
		in := make(chan int, 4)
		for _, v := range []int{1, 2, 3, 4} {
			in <- v
		}
		close(in)
		out := Map(ctx, in, func(n int) int { return -n })
		notNil(t, "Map(ctx, [1 2 3 4], negate)", out)
		if got := recv(t, "Map(ctx, [1 2 3 4], negate)", out); got != -1 {
			t.Fatalf("Map(ctx, [1 2 3 4], negate) value 0 = %d, want -1", got)
		}
		cancel()
		mustClose(t, "Map(ctx, [1 2 3 4], negate) canceled after 1 value", out)
	})
}

package merge

import (
	"runtime"
	"slices"
	"testing"
	"time"
)

// gen sends vals on a new unbuffered channel, then closes it.
func gen(vals ...int) <-chan int {
	ch := make(chan int)
	go func() {
		defer close(ch)
		for _, v := range vals {
			ch <- v
		}
	}()
	return ch
}

// recv receives one value from out, failing the test after 2s.
func recv(t *testing.T, call string, out <-chan int) (int, bool) {
	t.Helper()
	select {
	case v, ok := <-out:
		return v, ok
	case <-time.After(2 * time.Second):
		t.Fatalf("%s: output did not send or close within 2s", call)
		return 0, false
	}
}

// collect receives from out until it's closed (at most 1000 values).
func collect(t *testing.T, call string, out <-chan int) []int {
	t.Helper()
	if out == nil {
		t.Fatalf("%s returned a nil channel", call)
	}
	var got []int
	for {
		v, ok := recv(t, call, out)
		if !ok {
			return got
		}
		got = append(got, v)
		if len(got) > 1000 {
			t.Fatalf("%s sent more than 1000 values, first ones %v (forwarding zero values from a closed input?)", call, got[:10])
		}
	}
}

func TestMerge(t *testing.T) {
	got := collect(t, "Merge(1 2 3, 10 20 30 40)", Merge(gen(1, 2, 3), gen(10, 20, 30, 40)))

	var fromA, fromB []int
	for _, v := range got {
		if v < 10 {
			fromA = append(fromA, v)
		} else {
			fromB = append(fromB, v)
		}
	}
	if !slices.Equal(fromA, []int{1, 2, 3}) || !slices.Equal(fromB, []int{10, 20, 30, 40}) {
		t.Errorf("Merge(1 2 3, 10 20 30 40) = %v, want all 7 values with each input's order kept", got)
	}
}

func TestMergeClosedInput(t *testing.T) {
	closed := make(chan int)
	close(closed)
	if got, want := collect(t, "Merge(closed, 5 6)", Merge(closed, gen(5, 6))), []int{5, 6}; !slices.Equal(got, want) {
		t.Errorf("Merge(closed, 5 6) = %v, want %v", got, want)
	}
	if got, want := collect(t, "Merge(7 8 9, closed)", Merge(gen(7, 8, 9), closed)), []int{7, 8, 9}; !slices.Equal(got, want) {
		t.Errorf("Merge(7 8 9, closed) = %v, want %v", got, want)
	}
	if got := collect(t, "Merge(closed, closed)", Merge(closed, closed)); got != nil {
		t.Errorf("Merge(closed, closed) = %v, want no values", got)
	}
}

func TestMergeNilInputs(t *testing.T) {
	if got, want := collect(t, "Merge(nil, 7)", Merge(nil, gen(7))), []int{7}; !slices.Equal(got, want) {
		t.Errorf("Merge(nil, 7) = %v, want %v", got, want)
	}
	if got, want := collect(t, "Merge(4 5, nil)", Merge(gen(4, 5), nil)), []int{4, 5}; !slices.Equal(got, want) {
		t.Errorf("Merge(4 5, nil) = %v, want %v", got, want)
	}
	if got := collect(t, "Merge(nil, nil)", Merge[int](nil, nil)); got != nil {
		t.Errorf("Merge(nil, nil) = %v, want no values", got)
	}
}

func TestMergeDoesNotWaitForOneInput(t *testing.T) {
	a := make(chan int) // stays open and quiet until b is done
	out := Merge(a, gen(1, 2, 3))
	if out == nil {
		t.Fatalf("Merge(a, 1 2 3) returned a nil channel")
	}
	const call = "Merge(a open and quiet, 1 2 3)"
	for _, want := range []int{1, 2, 3} {
		if v, ok := recv(t, call, out); v != want || !ok {
			t.Fatalf("%s: received (%d, %v), want (%d, true)", call, v, ok, want)
		}
	}
	go func() {
		a <- 100
		close(a)
	}()
	if v, ok := recv(t, call, out); v != 100 || !ok {
		t.Fatalf("Merge: after b is done, a sends 100: received (%d, %v), want (100, true)", v, ok)
	}
	if v, ok := recv(t, call, out); ok {
		t.Fatalf("Merge: after both inputs closed: received (%d, %v), want the output closed", v, ok)
	}
}

func TestMergeStrings(t *testing.T) {
	a := make(chan string, 2)
	a <- "go"
	a <- "pher"
	close(a)
	out := Merge(a, nil)
	if out == nil {
		t.Fatalf("Merge(go pher, nil) returned a nil channel")
	}
	var got []string
	for range 3 {
		select {
		case v, ok := <-out:
			if !ok {
				if want := []string{"go", "pher"}; !slices.Equal(got, want) {
					t.Errorf("Merge(go pher, nil) = %q, want %q", got, want)
				}
				return
			}
			got = append(got, v)
		case <-time.After(2 * time.Second):
			t.Fatalf("Merge(go pher, nil): output did not send or close within 2s")
		}
	}
	t.Errorf("Merge(go pher, nil) = %q and still open, want [go pher] then closed", got)
}

func TestMergeNoLeak(t *testing.T) {
	before := runtime.NumGoroutine()
	for range 20 {
		collect(t, "Merge(1 2, 3 4)", Merge(gen(1, 2), gen(3, 4)))
	}
	deadline := time.Now().Add(time.Second)
	for runtime.NumGoroutine() > before && time.Now().Before(deadline) {
		time.Sleep(5 * time.Millisecond)
	}
	if n := runtime.NumGoroutine(); n > before {
		t.Errorf("goroutines after 20 finished Merges = %d, want %d (does the forwarding goroutine return?)", n, before)
	}
}

package queue

import (
	"fmt"
	"testing"
	"time"
)

// within runs f and fails if it doesn't return within 2s (a blocked send or
// receive). f must only use t.Errorf.
func within(t *testing.T, what string, f func()) {
	t.Helper()
	done := make(chan struct{})
	go func() {
		defer close(done)
		f()
	}()
	select {
	case <-done:
	case <-time.After(2 * time.Second):
		t.Fatalf("%s blocked for 2s, want Push and Pop never to block", what)
	}
}

func TestNew(t *testing.T) {
	for _, c := range []int{0, 1, 5, 100} {
		q := New(c)
		if q == nil {
			t.Fatalf("New(%d) = nil, want a queue", c)
		}
		if q.ch == nil {
			t.Fatalf("New(%d).ch = nil, want a channel made with make(chan string, %d)", c, c)
		}
		if got := q.Cap(); got != c {
			t.Errorf("New(%d).Cap() = %d, want %d", c, got, c)
		}
		if got := q.Len(); got != 0 {
			t.Errorf("New(%d).Len() = %d, want 0", c, got)
		}
	}
}

func TestReadmeScenario(t *testing.T) {
	q := New(2)
	if q == nil || q.ch == nil {
		t.Fatalf("New(2) has no channel, want make(chan string, 2)")
	}
	within(t, "the README scenario", func() {
		steps := []struct {
			name string
			got  any
			want any
		}{
			{`Push("a")`, q.Push("a"), true},
			{`Push("b")`, q.Push("b"), true},
			{`Push("c") on a full queue`, q.Push("c"), false},
			{"Len()", q.Len(), 2},
			{"Pop()", pop(q), `("a", true)`},
			{"Len()", q.Len(), 1},
			{"Pop()", pop(q), `("b", true)`},
			{"Pop() on an empty queue", pop(q), `("", false)`},
			{"Len()", q.Len(), 0},
			{"Cap()", q.Cap(), 2},
		}
		for _, s := range steps {
			if s.got != s.want {
				t.Errorf("%s = %v, want %v", s.name, s.got, s.want)
			}
		}
	})
}

func TestFIFOAndReuse(t *testing.T) {
	q := New(3)
	if q == nil || q.ch == nil {
		t.Fatalf("New(3) has no channel, want make(chan string, 3)")
	}
	within(t, "Push/Pop", func() {
		var got []string
		for round := range 4 {
			for i := range 3 {
				s := fmt.Sprintf("r%d-%d", round, i)
				if !q.Push(s) {
					t.Errorf("round %d: Push(%q) = false with Len() %d, want true", round, s, q.Len())
				}
			}
			if q.Push("extra") {
				t.Errorf("round %d: Push(%q) on a full queue = true, want false", round, "extra")
			}
			for range 3 {
				v, _ := q.Pop()
				got = append(got, v)
			}
		}
		want := "[r0-0 r0-1 r0-2 r1-0 r1-1 r1-2 r2-0 r2-1 r2-2 r3-0 r3-1 r3-2]"
		if fmt.Sprint(got) != want {
			t.Errorf("Pop order = %v, want %v", got, want)
		}
	})
}

func TestEmptyStringIsAValue(t *testing.T) {
	q := New(1)
	if q == nil || q.ch == nil {
		t.Fatalf("New(1) has no channel, want make(chan string, 1)")
	}
	within(t, "Push/Pop", func() {
		if !q.Push("") {
			t.Errorf(`Push("") = false, want true`)
		}
		if v, ok := q.Pop(); v != "" || !ok {
			t.Errorf(`Pop() after Push("") = (%q, %v), want ("", true)`, v, ok)
		}
		if v, ok := q.Pop(); v != "" || ok {
			t.Errorf(`Pop() on an empty queue = (%q, %v), want ("", false)`, v, ok)
		}
	})
}

func TestCapacityZero(t *testing.T) {
	q := New(0)
	if q == nil || q.ch == nil {
		t.Fatalf("New(0) has no channel, want make(chan string, 0)")
	}
	within(t, "Push/Pop on New(0)", func() {
		if q.Push("x") {
			t.Errorf(`New(0).Push("x") = true, want false`)
		}
		if v, ok := q.Pop(); ok {
			t.Errorf(`New(0).Pop() = (%q, true), want ("", false)`, v)
		}
	})
}

// pop formats the result of q.Pop as ("v", ok).
func pop(q *Queue) string {
	v, ok := q.Pop()
	return fmt.Sprintf("(%q, %v)", v, ok)
}

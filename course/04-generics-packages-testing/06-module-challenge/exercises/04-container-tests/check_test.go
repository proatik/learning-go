package containers

import (
	"iter"
	"os"
	"os/exec"
	"regexp"
	"slices"
	"strings"
	"testing"
)

// This file holds the code under test (a generic Stack and Queue) and the
// checker. During the bug hunt one planted bug is switched on at a time, and
// your tests in containers_test.go must fail for every one of them.

// bugs lists the planted bugs; the code below checks active to enable one.
var bugs = []string{
	"stack Pop takes the bottom item",
	"stack Pop on empty reports ok",
	"stack Peek removes the item",
	"stack All yields bottom to top",
	"stack All ignores break",
	"queue Pop forgets to shrink Len",
	"queue All ignores the head position",
	"queue grow scrambles a wrapped buffer",
	"queue Pop on empty returns stale value",
}

const bugEnv = "LEARNGO_BUG"

var active = os.Getenv(bugEnv)

// Stack is a last-in, first-out stack.
type Stack[T any] struct {
	items []T
}

// NewStack returns an empty stack.
func NewStack[T any]() *Stack[T] { return &Stack[T]{} }

// Push adds v on top.
func (s *Stack[T]) Push(v T) { s.items = append(s.items, v) }

// Pop removes and returns the top item, or the zero value and false.
func (s *Stack[T]) Pop() (T, bool) {
	var zero T
	if len(s.items) == 0 {
		return zero, active == "stack Pop on empty reports ok"
	}
	if active == "stack Pop takes the bottom item" {
		v := s.items[0]
		s.items = s.items[1:]
		return v, true
	}
	v := s.items[len(s.items)-1]
	s.items = s.items[:len(s.items)-1]
	return v, true
}

// Peek returns the top item without removing it.
func (s *Stack[T]) Peek() (T, bool) {
	if active == "stack Peek removes the item" {
		return s.Pop()
	}
	if len(s.items) == 0 {
		var zero T
		return zero, false
	}
	return s.items[len(s.items)-1], true
}

// Len returns the number of items.
func (s *Stack[T]) Len() int { return len(s.items) }

// All yields the items from top to bottom.
func (s *Stack[T]) All() iter.Seq[T] {
	return func(yield func(T) bool) {
		if active == "stack All yields bottom to top" {
			for _, v := range s.items {
				if !yield(v) {
					return
				}
			}
			return
		}
		for i := len(s.items) - 1; i >= 0; i-- {
			if !yield(s.items[i]) && active != "stack All ignores break" {
				return
			}
		}
	}
}

// Queue is a first-in, first-out queue stored in a ring buffer that starts
// with room for 4 items and doubles when full.
type Queue[T any] struct {
	buf  []T
	head int // index of the front item
	n    int // number of items
}

// NewQueue returns an empty queue.
func NewQueue[T any]() *Queue[T] { return &Queue[T]{buf: make([]T, 4)} }

// Push adds v at the back.
func (q *Queue[T]) Push(v T) {
	if q.n == len(q.buf) {
		q.grow()
	}
	q.buf[(q.head+q.n)%len(q.buf)] = v
	q.n++
}

// grow doubles the buffer, moving the front item to index 0.
func (q *Queue[T]) grow() {
	nb := make([]T, 2*len(q.buf))
	if active == "queue grow scrambles a wrapped buffer" {
		copy(nb, q.buf)
	} else {
		for i := range q.n {
			nb[i] = q.buf[(q.head+i)%len(q.buf)]
		}
	}
	q.buf, q.head = nb, 0
}

// Pop removes and returns the front item, or the zero value and false.
func (q *Queue[T]) Pop() (T, bool) {
	var zero T
	if q.n == 0 {
		if active == "queue Pop on empty returns stale value" {
			return q.buf[(q.head+len(q.buf)-1)%len(q.buf)], false
		}
		return zero, false
	}
	v := q.buf[q.head]
	if active != "queue Pop on empty returns stale value" {
		q.buf[q.head] = zero
	}
	q.head = (q.head + 1) % len(q.buf)
	if active != "queue Pop forgets to shrink Len" {
		q.n--
	}
	return v, true
}

// Len returns the number of items.
func (q *Queue[T]) Len() int { return q.n }

// All yields the items from front to back.
func (q *Queue[T]) All() iter.Seq[T] {
	return func(yield func(T) bool) {
		for i := range q.n {
			idx := (q.head + i) % len(q.buf)
			if active == "queue All ignores the head position" {
				idx = i
			}
			if !yield(q.buf[idx]) {
				return
			}
		}
	}
}

// TestYourTestsCatchEveryBug re-runs your tests once per bug: each run must
// fail. It also requires at least one subtest (t.Run).
func TestYourTestsCatchEveryBug(t *testing.T) {
	if active != "" {
		t.Skip("inside the bug hunt")
	}
	run := func(bug string, args ...string) (string, error) {
		cmd := exec.Command(os.Args[0], append([]string{"-test.count=1", "-test.skip=^TestYourTestsCatchEveryBug$"}, args...)...)
		cmd.Env = append(os.Environ(), bugEnv+"="+bug)
		out, err := cmd.CombinedOutput()
		return string(out), err
	}
	list, _ := run("", "-test.list=^Test")
	if !slices.ContainsFunc(strings.Fields(list), func(n string) bool { return n != "TestYourTestsCatchEveryBug" && strings.HasPrefix(n, "Test") }) {
		t.Fatal("no Test functions found: write at least one test")
	}
	out, err := run("", "-test.run=^Test", "-test.v")
	if err != nil {
		t.Fatalf("your tests fail on the correct Stack and Queue:\n%s", out)
	}
	if !regexp.MustCompile(`(?m)^=== RUN\s+Test\w+/`).MatchString(out) {
		t.Errorf("subtests = none, want at least one t.Run")
	}
	for _, name := range bugs {
		if _, err := run(name, "-test.run=^Test"); err == nil {
			t.Errorf("Bug %q = not caught, want caught", name)
		}
	}
}

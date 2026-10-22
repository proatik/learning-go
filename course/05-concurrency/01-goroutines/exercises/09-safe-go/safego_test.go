package safego_test

import (
	"errors"
	"fmt"
	"runtime"
	"sync/atomic"
	"testing"
	"testing/synctest"
	"time"

	safego "learngo/course/05-concurrency/01-goroutines/exercises/09-safe-go"
)

var ErrDisk = errors.New("disk full")

type codeError struct{ code int }

func (e *codeError) Error() string { return fmt.Sprintf("code %d", e.code) }

// within runs f in a goroutine and fails the test if it doesn't return in 2s.
func within(t *testing.T, name string, f func()) {
	t.Helper()
	done := make(chan struct{})
	go func() {
		defer close(done)
		f()
	}()
	select {
	case <-done:
	case <-time.After(2 * time.Second):
		t.Fatalf("%s did not return within 2s", name)
	}
}

func messages(errs []error) []string {
	out := make([]string, len(errs))
	for i, err := range errs {
		if err != nil {
			out[i] = err.Error()
		}
	}
	return out
}

func TestSafeRun(t *testing.T) {
	ok := func() {}
	tests := []struct {
		name  string
		tasks []func()
		want  []string // "" means nil
	}{
		{"ok, panic(boom), ok", []func(){ok, func() { panic("boom") }, ok}, []string{"", "task 1 panicked: boom", ""}},
		{"panic(42)", []func(){func() { panic(42) }}, []string{"task 0 panicked: 42"}},
		{"panic(ErrDisk)", []func(){func() { panic(ErrDisk) }}, []string{"task 0 panicked: disk full"}},
		{"ok, nil task", []func(){ok, nil}, []string{"", "task 1 panicked: runtime error: invalid memory address or nil pointer dereference"}},
		{"ok, ok", []func(){ok, ok}, []string{"", ""}},
		{"3 panics", []func(){func() { panic("a") }, func() { panic("b") }, func() { panic([]int{1, 2}) }},
			[]string{"task 0 panicked: a", "task 1 panicked: b", "task 2 panicked: [1 2]"}},
		{"index out of range", []func(){ok, ok, ok, func() { s := []int{}; _ = s[len(s)+2] }},
			[]string{"", "", "", "task 3 panicked: runtime error: index out of range [2] with length 0"}},
	}
	for _, tt := range tests {
		var errs []error
		within(t, "SafeRun("+tt.name+")", func() { errs = safego.SafeRun(tt.tasks) })
		if len(errs) != len(tt.tasks) {
			t.Errorf("SafeRun(%s) = %d results, want %d", tt.name, len(errs), len(tt.tasks))
			continue
		}
		for i, got := range messages(errs) {
			if (errs[i] == nil) != (tt.want[i] == "") || got != tt.want[i] {
				t.Errorf("SafeRun(%s)[%d] = %v, want %q", tt.name, i, errs[i], tt.want[i])
			}
		}
	}
}

func TestWrapsErrorPanics(t *testing.T) {
	var errs []error
	tasks := []func(){
		func() { panic(ErrDisk) },
		func() { panic(&codeError{7}) },
		func() { var m map[string]int; m["x"] = 1 },
		func() { panic("just a string") },
	}
	within(t, "SafeRun", func() { errs = safego.SafeRun(tasks) })
	if len(errs) != len(tasks) {
		t.Fatalf("SafeRun(4 panicking tasks) = %d results, want 4", len(errs))
	}
	if !errors.Is(errs[0], ErrDisk) {
		t.Errorf("SafeRun(panic(ErrDisk)): errors.Is(err, ErrDisk) = false, want true (wrap with %%w)")
	}
	var ce *codeError
	if !errors.As(errs[1], &ce) || ce.code != 7 {
		t.Errorf("SafeRun(panic(&codeError{7})): errors.As(err, *codeError) = false, want true with code 7")
	}
	var re runtime.Error
	if !errors.As(errs[2], &re) {
		t.Errorf("SafeRun(write to nil map): errors.As(err, runtime.Error) = false, want true")
	}
	if errs[3] == nil || errors.Unwrap(errs[3]) != nil {
		t.Errorf("SafeRun(panic(string)): errors.Unwrap(err) = %v, want nil (only errors are wrapped)", errors.Unwrap(errs[3]))
	}
	if got, want := fmt.Sprint(errs[1]), "task 1 panicked: code 7"; got != want {
		t.Errorf("SafeRun(panic(&codeError{7}))[1] = %q, want %q", got, want)
	}
}

func TestConcurrentAndIsolated(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		var finished atomic.Int32
		sleepThen := func(d time.Duration, fail bool) func() {
			return func() {
				time.Sleep(d)
				if fail {
					panic(fmt.Sprintf("failed after %v", d))
				}
				finished.Add(1)
			}
		}
		tasks := []func(){
			sleepThen(time.Second, true),
			sleepThen(3*time.Second, false),
			sleepThen(2*time.Second, true),
			sleepThen(3*time.Second, false),
		}
		start := time.Now()
		errs := safego.SafeRun(tasks)
		elapsed := time.Since(start)
		want := []string{"task 0 panicked: failed after 1s", "", "task 2 panicked: failed after 2s", ""}
		if len(errs) != len(want) {
			t.Fatalf("SafeRun(4 sleeping tasks) = %d results, want %d", len(errs), len(want))
		}
		for i, got := range messages(errs) {
			if got != want[i] {
				t.Errorf("SafeRun(4 sleeping tasks)[%d] = %v, want %q", i, errs[i], want[i])
			}
		}
		if f := finished.Load(); f != 2 {
			t.Errorf("SafeRun(4 sleeping tasks, 2 panic): tasks finished normally = %d, want 2", f)
		}
		if elapsed != 3*time.Second {
			t.Errorf("SafeRun(tasks of 1s, 3s, 2s, 3s) took %v, want 3s", elapsed)
		}
	})
}

func TestEmpty(t *testing.T) {
	for _, tasks := range [][]func(){nil, {}} {
		var errs []error
		within(t, "SafeRun(empty)", func() { errs = safego.SafeRun(tasks) })
		if len(errs) != 0 {
			t.Errorf("SafeRun(%d tasks) = %v, want []", len(tasks), errs)
		}
	}
}

func TestMany(t *testing.T) {
	const n = 1000
	tasks := make([]func(), n)
	for i := range tasks {
		if i%10 == 7 {
			tasks[i] = func() { panic(i) }
		} else {
			tasks[i] = func() {}
		}
	}
	var errs []error
	within(t, "SafeRun(1000 tasks)", func() { errs = safego.SafeRun(tasks) })
	if len(errs) != n {
		t.Fatalf("SafeRun(1000 tasks) = %d results, want %d", len(errs), n)
	}
	for i, err := range errs {
		want := ""
		if i%10 == 7 {
			want = fmt.Sprintf("task %d panicked: %d", i, i)
		}
		if got := messages([]error{err})[0]; got != want {
			t.Fatalf("SafeRun(1000 tasks)[%d] = %v, want %q", i, err, want)
		}
	}
}

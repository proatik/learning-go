package catch_test

import (
	"errors"
	"fmt"
	"io"
	"runtime"
	"slices"
	"testing"

	catch "learngo/course/03-methods-interfaces-errors/05-panic-recover/exercises/04-try"
)

// noPanic runs f and reports a panic that escaped instead of crashing.
func noPanic(t *testing.T, name string, f func()) {
	t.Helper()
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("%s panicked: %v", name, r)
		}
	}()
	f()
}

var errBoom = errors.New("boom error")

func TestTryNoPanic(t *testing.T) {
	noPanic(t, "Try", func() {
		called := false
		err := catch.Try(func() { called = true })
		if err != nil {
			t.Errorf("Try(no panic) = %v, want <nil>", err)
		}
		if !called {
			t.Errorf("Try(f) did not call f")
		}
	})
}

func TestTry(t *testing.T) {
	tests := []struct {
		name string
		f    func()
		want string
	}{
		{"string", func() { panic("boom") }, "recovered: boom"},
		{"int", func() { panic(42) }, "recovered: 42"},
		{"error", func() { panic(io.EOF) }, "recovered: EOF"},
		{"wrapped", func() { panic(fmt.Errorf("load: %w", errBoom)) }, "recovered: load: boom error"},
		{"panic(nil)", func() { panic(nil) }, "recovered: runtime error: panic called with nil argument"},
		{"nil map", func() {
			var m map[string]bool
			m["x"] = true
		}, "recovered: assignment to entry in nil map"},
		{"nil func", nil, "recovered: runtime error: invalid memory address or nil pointer dereference"},
	}
	for _, tt := range tests {
		noPanic(t, "Try", func() {
			err := catch.Try(tt.f)
			if err == nil {
				t.Errorf("Try(%s) = <nil>, want %q", tt.name, tt.want)
				return
			}
			if err.Error() != tt.want {
				t.Errorf("Try(%s) = %q, want %q", tt.name, err.Error(), tt.want)
			}
		})
	}
}

func TestTryWrapsErrors(t *testing.T) {
	noPanic(t, "Try", func() {
		if err := catch.Try(func() { panic(io.EOF) }); !errors.Is(err, io.EOF) {
			t.Errorf("Try(panic(io.EOF)) = %v, want errors.Is(err, io.EOF)", err)
		}
		if err := catch.Try(func() { panic(fmt.Errorf("load: %w", errBoom)) }); !errors.Is(err, errBoom) {
			t.Errorf("Try(panic(wrapped errBoom)) = %v, want errors.Is(err, errBoom)", err)
		}
		var pn *runtime.PanicNilError
		if err := catch.Try(func() { panic(nil) }); !errors.As(err, &pn) {
			t.Errorf("Try(panic(nil)) = %v, want errors.As(err, *runtime.PanicNilError)", err)
		}
		var re runtime.Error
		if err := catch.Try(func() {
			s := []int{}
			_ = s[len(s)]
		}); !errors.As(err, &re) {
			t.Errorf("Try(index out of range) = %v, want errors.As(err, runtime.Error)", err)
		}
	})
}

func TestTryRunsDeferredCallsOfF(t *testing.T) {
	noPanic(t, "Try", func() {
		var log []string
		err := catch.Try(func() {
			defer func() { log = append(log, "cleanup") }()
			log = append(log, "work")
			panic("fail")
		})
		if err == nil || err.Error() != "recovered: fail" {
			t.Errorf("Try(f) = %v, want %q", err, "recovered: fail")
		}
		if want := []string{"work", "cleanup"}; !slices.Equal(log, want) {
			t.Errorf("log = %q, want %q", log, want)
		}
	})
}

func TestTryInt(t *testing.T) {
	noPanic(t, "TryInt", func() {
		n, err := catch.TryInt(func() int { return 7 })
		if n != 7 || err != nil {
			t.Errorf("TryInt(return 7) = %d, %v; want 7, <nil>", n, err)
		}
		n, err = catch.TryInt(func() int { return -1 })
		if n != -1 || err != nil {
			t.Errorf("TryInt(return -1) = %d, %v; want -1, <nil>", n, err)
		}
	})
	noPanic(t, "TryInt", func() {
		n, err := catch.TryInt(func() int { panic("x") })
		if n != 0 || err == nil || err.Error() != "recovered: x" {
			t.Errorf("TryInt(panic(\"x\")) = %d, %v; want 0, recovered: x", n, err)
		}
	})
	noPanic(t, "TryInt", func() {
		zero := 0
		n, err := catch.TryInt(func() int { return 10 / zero })
		want := "recovered: runtime error: integer divide by zero"
		if n != 0 || err == nil || err.Error() != want {
			t.Errorf("TryInt(10 / 0) = %d, %v; want 0, %s", n, err, want)
		}
		var re runtime.Error
		if !errors.As(err, &re) {
			t.Errorf("TryInt(10 / 0) error = %v, want it to wrap a runtime.Error", err)
		}
	})
}

func TestTryAll(t *testing.T) {
	noPanic(t, "TryAll", func() {
		var ran []int
		jobs := []func(){
			func() { ran = append(ran, 0) },
			func() { ran = append(ran, 1); panic("boom") },
			func() { ran = append(ran, 2) },
			func() { ran = append(ran, 3); panic(io.ErrUnexpectedEOF) },
			func() { ran = append(ran, 4) },
		}
		errs := catch.TryAll(jobs)
		if len(errs) != len(jobs) {
			t.Fatalf("TryAll(5 jobs) returned %d errors, want %d", len(errs), len(jobs))
		}
		want := []string{"", "recovered: boom", "", "recovered: unexpected EOF", ""}
		for i, err := range errs {
			got := ""
			if err != nil {
				got = err.Error()
			}
			if got != want[i] {
				t.Errorf("TryAll(jobs)[%d] = %q, want %q", i, got, want[i])
			}
		}
		if !errors.Is(errs[3], io.ErrUnexpectedEOF) {
			t.Errorf("TryAll(jobs)[3] = %v, want errors.Is(err, io.ErrUnexpectedEOF)", errs[3])
		}
		if wantRan := []int{0, 1, 2, 3, 4}; !slices.Equal(ran, wantRan) {
			t.Errorf("TryAll ran jobs %v, want %v", ran, wantRan)
		}
	})
	noPanic(t, "TryAll", func() {
		if errs := catch.TryAll(nil); len(errs) != 0 {
			t.Errorf("TryAll(nil) = %v, want empty", errs)
		}
	})
}

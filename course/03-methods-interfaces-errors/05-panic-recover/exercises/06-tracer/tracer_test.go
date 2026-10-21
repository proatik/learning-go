package tracer

import (
	"fmt"
	"os"
	"strings"
	"testing"
)

// safe runs f and reports an escaped panic instead of crashing.
func safe(f func()) (r any) {
	defer func() { r = recover() }()
	f()
	return nil
}

func fact(t *Tracer, n int) int {
	defer t.Enter(fmt.Sprintf("fact(%d)", n))()
	if n <= 1 {
		return 1
	}
	return n * fact(t, n-1)
}

func ExampleTracer_Enter() {
	t := &Tracer{W: os.Stdout}
	if r := safe(func() { fmt.Println("result:", fact(t, 3)) }); r != nil {
		fmt.Println("PANIC:", r)
	}
	// Output:
	// > fact(3)
	// | > fact(2)
	// | | > fact(1)
	// | | < fact(1)
	// | < fact(2)
	// < fact(3)
	// result: 6
}

func ExampleTracer_Enter_siblings() {
	t := &Tracer{W: os.Stdout}
	a := func() { defer t.Enter("a")() }
	b := func() { defer t.Enter("b")() }
	run := func() {
		defer t.Enter("main")()
		a()
		b()
	}
	if r := safe(run); r != nil {
		fmt.Println("PANIC:", r)
	}
	// Output:
	// > main
	// | > a
	// | < a
	// | > b
	// | < b
	// < main
}

func TestEnterPanicStillExits(t *testing.T) {
	var sb strings.Builder
	tr := &Tracer{W: &sb}
	r := safe(func() {
		defer tr.Enter("outer")()
		func() {
			defer tr.Enter("inner")()
			panic("boom")
		}()
	})
	if r != "boom" {
		t.Errorf("panic value = %v, want boom (Enter must not recover)", r)
	}
	if r := safe(func() { defer tr.Enter("after")() }); r != nil {
		t.Fatalf("Enter panicked: %v", r)
	}
	want := "> outer\n| > inner\n| < inner\n< outer\n> after\n< after\n"
	if got := sb.String(); got != want {
		t.Errorf("trace after a panic =\n%s\nwant\n%s", got, want)
	}
}

func TestEnterManual(t *testing.T) {
	var sb strings.Builder
	tr := &Tracer{W: &sb}
	if r := safe(func() {
		exitA := tr.Enter("a")
		exitB := tr.Enter("b")
		exitB()
		exitC := tr.Enter("c")
		exitC()
		exitA()
	}); r != nil {
		t.Fatalf("Enter panicked: %v", r)
	}
	want := "> a\n| > b\n| < b\n| > c\n| < c\n< a\n"
	if got := sb.String(); got != want {
		t.Errorf("trace =\n%s\nwant\n%s", got, want)
	}
}

func TestEnterUnicode(t *testing.T) {
	var sb strings.Builder
	tr := &Tracer{W: &sb}
	if r := safe(func() { defer tr.Enter("größe → λ")() }); r != nil {
		t.Fatalf("Enter panicked: %v", r)
	}
	if got, want := sb.String(), "> größe → λ\n< größe → λ\n"; got != want {
		t.Errorf("Enter(%q) trace = %q, want %q", "größe → λ", got, want)
	}
}

func TestEnterNilWriter(t *testing.T) {
	tr := &Tracer{}
	var exit func()
	if r := safe(func() { exit = tr.Enter("quiet") }); r != nil {
		t.Fatalf("Enter with nil W panicked: %v", r)
	}
	if exit == nil {
		t.Fatalf("Enter with nil W returned a nil func, want a no-op func")
	}
	if r := safe(exit); r != nil {
		t.Errorf("exit func with nil W panicked: %v", r)
	}
	if tr.depth != 0 {
		t.Errorf("depth after Enter with nil W = %d, want 0", tr.depth)
	}
}

package commands_test

import (
	"fmt"
	"slices"
	"strings"
	"testing"

	commands "learngo/course/03-methods-interfaces-errors/02-interfaces/exercises/09-command-registry"
)

// greeter is a handler type the package knows nothing about.
type greeter struct{ greeting string }

func (g greeter) Handle(args []string) string {
	if len(args) == 0 {
		return g.greeting + ", stranger"
	}
	return g.greeting + ", " + args[0]
}

func TestHandlerFunc(t *testing.T) {
	var got []string
	f := commands.HandlerFunc(func(args []string) string {
		got = args
		return fmt.Sprint(len(args))
	})
	var h commands.Handler = f
	if out := h.Handle([]string{"a", "b"}); out != "2" {
		t.Errorf("HandlerFunc(f).Handle([a b]) = %q, want %q", out, "2")
	}
	if !slices.Equal(got, []string{"a", "b"}) {
		t.Errorf("HandlerFunc(f).Handle([a b]) passed %q to f, want [a b]", got)
	}
}

func TestEcho(t *testing.T) {
	tests := []struct {
		args []string
		want string
	}{
		{[]string{"hi", "there"}, "hi there"},
		{[]string{"one"}, "one"},
		{nil, ""},
		{[]string{"ünï", "cödé"}, "ünï cödé"},
	}
	for _, tt := range tests {
		if got := (commands.Echo{}).Handle(tt.args); got != tt.want {
			t.Errorf("Echo{}.Handle(%q) = %q, want %q", tt.args, got, tt.want)
		}
	}
}

func TestCounter(t *testing.T) {
	c := &commands.Counter{}
	for i := 1; i <= 3; i++ {
		want := fmt.Sprintf("count=%d", i)
		if got := c.Handle(nil); got != want {
			t.Errorf("call %d: Counter.Handle(nil) = %q, want %q", i, got, want)
		}
	}
	other := &commands.Counter{}
	if got := other.Handle([]string{"x"}); got != "count=1" {
		t.Errorf("new Counter: Handle([x]) = %q, want %q", got, "count=1")
	}
}

func TestUpper(t *testing.T) {
	h := commands.Upper(commands.Echo{})
	if h == nil {
		t.Fatalf("Upper(Echo{}) = nil, want a Handler")
	}
	if got := h.Handle([]string{"go", "rocks"}); got != "GO ROCKS" {
		t.Errorf("Upper(Echo{}).Handle([go rocks]) = %q, want %q", got, "GO ROCKS")
	}
	g := commands.Upper(greeter{"hello"})
	if got := g.Handle(nil); got != "HELLO, STRANGER" {
		t.Errorf("Upper(greeter).Handle(nil) = %q, want %q", got, "HELLO, STRANGER")
	}
}

func TestRegister(t *testing.T) {
	var r commands.Registry // zero value must work
	tests := []struct {
		name string
		h    commands.Handler
		want bool
	}{
		{"echo", commands.Echo{}, true},
		{"echo", commands.Echo{}, false},
		{"count", &commands.Counter{}, true},
		{"Echo", commands.Echo{}, true}, // case-sensitive
		{"", commands.Echo{}, false},
		{"two words", commands.Echo{}, false},
		{"tab\there", commands.Echo{}, false},
		{"nl\n", commands.Echo{}, false},
		{"nilface", nil, false},
		{"nilfunc", commands.HandlerFunc(nil), false},
		{"nilcounter", (*commands.Counter)(nil), false},
		{"greet", greeter{"hi"}, true},
		{"héllo", commands.Echo{}, true},
	}
	for _, tt := range tests {
		if got := r.Register(tt.name, tt.h); got != tt.want {
			t.Errorf("Register(%q, %T) = %v, want %v", tt.name, tt.h, got, tt.want)
		}
	}
	want := []string{"Echo", "count", "echo", "greet", "héllo"}
	if got := r.Names(); !slices.Equal(got, want) {
		t.Errorf("Names() = %q, want %q", got, want)
	}
}

func TestNamesEmpty(t *testing.T) {
	var r commands.Registry
	if got := r.Names(); len(got) != 0 {
		t.Errorf("Names() on empty registry = %q, want []", got)
	}
}

func TestDispatch(t *testing.T) {
	var r commands.Registry
	var seen []string
	r.Register("echo", commands.Echo{})
	r.Register("count", &commands.Counter{})
	r.Register("shout", commands.Upper(commands.Echo{}))
	r.Register("greet", greeter{"hi"})
	r.Register("args", commands.HandlerFunc(func(args []string) string {
		seen = args
		return fmt.Sprintf("%d args", len(args))
	}))

	tests := []struct {
		line, want string
	}{
		{"echo a b", "a b"},
		{"  echo   a  b ", "a b"},
		{"echo", ""},
		{"count", "count=1"},
		{"count x y", "count=2"},
		{"shout hey you", "HEY YOU"},
		{"greet", "hi, stranger"},
		{"greet ann bob", "hi, ann"},
		{"args 1 2 3", "3 args"},
		{"args", "0 args"},
		{"nope x", "unknown command: nope"},
		{"ECHO a", "unknown command: ECHO"},
		{"", ""},
		{" \t\n ", ""},
		{"help", "commands: args count echo greet shout"},
		{"help me", "commands: args count echo greet shout"},
	}
	for _, tt := range tests {
		if got := r.Dispatch(tt.line); got != tt.want {
			t.Errorf("Dispatch(%q) = %q, want %q", tt.line, got, tt.want)
		}
	}

	r.Dispatch("args")
	if seen == nil || len(seen) != 0 {
		t.Errorf("Dispatch(%q) passed %#v to the handler, want an empty non-nil slice", "args", seen)
	}
}

func TestDispatchHelp(t *testing.T) {
	var empty commands.Registry
	if got := empty.Dispatch("help"); got != "no commands" {
		t.Errorf("Dispatch(%q) on empty registry = %q, want %q", "help", got, "no commands")
	}
	var r commands.Registry
	r.Register("help", commands.HandlerFunc(func(args []string) string {
		return "custom help: " + strings.Join(args, ",")
	}))
	if got := r.Dispatch("help a b"); got != "custom help: a,b" {
		t.Errorf("Dispatch(%q) with registered help = %q, want %q", "help a b", got, "custom help: a,b")
	}
}

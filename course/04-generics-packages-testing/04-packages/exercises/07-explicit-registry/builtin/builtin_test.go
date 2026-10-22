package builtin_test

import (
	"errors"
	"go/ast"
	"go/parser"
	"go/token"
	"testing"

	"learngo/course/04-generics-packages-testing/04-packages/exercises/07-explicit-registry/builtin"
	"learngo/course/04-generics-packages-testing/04-packages/exercises/07-explicit-registry/render"
)

func TestFormats(t *testing.T) {
	tests := []struct {
		format string
		f      func(string, []string) string
		title  string
		items  []string
		want   string
	}{
		{"Plain", builtin.Plain, "Todo", []string{"a", "b"}, "Todo: a, b"},
		{"Plain", builtin.Plain, "Todo", []string{"one"}, "Todo: one"},
		{"Plain", builtin.Plain, "Todo", nil, "Todo: (none)"},
		{"Plain", builtin.Plain, "", []string{}, ": (none)"},
		{"Markdown", builtin.Markdown, "Todo", []string{"a", "b"}, "# Todo\n\n- a\n- b"},
		{"Markdown", builtin.Markdown, "Tâches", []string{"café"}, "# Tâches\n\n- café"},
		{"Markdown", builtin.Markdown, "Todo", nil, "# Todo\n\n_No items._"},
		{"Numbered", builtin.Numbered, "Todo", []string{"a", "b"}, "Todo\n1. a\n2. b"},
		{"Numbered", builtin.Numbered, "Steps", []string{"x", "y", "z"}, "Steps\n1. x\n2. y\n3. z"},
		{"Numbered", builtin.Numbered, "Todo", nil, "Todo"},
	}
	for _, tt := range tests {
		if got := tt.f(tt.title, tt.items); got != tt.want {
			t.Errorf("%s(%q, %q) = %q, want %q", tt.format, tt.title, tt.items, got, tt.want)
		}
	}
}

func TestRegisterAll(t *testing.T) {
	var r render.Registry
	if err := builtin.RegisterAll(&r); err != nil {
		t.Fatalf("RegisterAll error = %v, want <nil>", err)
	}
	got := r.Names()
	want := []string{"markdown", "numbered", "plain"}
	if len(got) != len(want) || got[0] != want[0] || got[1] != want[1] || got[2] != want[2] {
		t.Fatalf("Names() = %q, want %q", got, want)
	}
	out, err := r.Render("numbered", "T", []string{"a"})
	if out != "T\n1. a" || err != nil {
		t.Errorf("Render(%q, %q, [a]) = %q, %v, want %q, <nil>", "numbered", "T", out, err, "T\n1. a")
	}
}

func TestRegisterAllTwice(t *testing.T) {
	var r render.Registry
	if err := builtin.RegisterAll(&r); err != nil {
		t.Fatalf("first RegisterAll error = %v, want <nil>", err)
	}
	if err := builtin.RegisterAll(&r); !errors.Is(err, render.ErrDuplicate) {
		t.Errorf("second RegisterAll error = %v, want %v", err, render.ErrDuplicate)
	}
}

func TestRegisterAllStopsAtFirstError(t *testing.T) {
	var r render.Registry
	if err := r.Register("markdown", builtin.Plain); err != nil {
		t.Fatalf("Register(%q) error = %v, want <nil>", "markdown", err)
	}
	err := builtin.RegisterAll(&r)
	if want := `render: already registered: "markdown"`; err == nil || err.Error() != want {
		t.Errorf("RegisterAll with %q taken error = %v, want %q", "markdown", err, want)
	}
	if got := r.Names(); len(got) != 2 {
		t.Errorf("Names() = %q, want [markdown plain] (numbered not reached)", got)
	}
}

func TestNoInit(t *testing.T) {
	f, err := parser.ParseFile(token.NewFileSet(), "builtin.go", nil, 0)
	if err != nil {
		t.Fatalf("parse builtin.go: %v", err)
	}
	for _, d := range f.Decls {
		if fn, ok := d.(*ast.FuncDecl); ok && fn.Name.Name == "init" {
			t.Errorf("func init in builtin.go = present, want removed (callers use RegisterAll)")
		}
	}
}

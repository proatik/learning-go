package render_test

import (
	"errors"
	"go/ast"
	"go/parser"
	"go/token"
	"strings"
	"testing"

	"learngo/course/04-generics-packages-testing/04-packages/exercises/07-explicit-registry/render"
)

func upper(title string, items []string) string {
	return strings.ToUpper(title) + ":" + strings.Join(items, ",")
}

func TestRegisterAndRender(t *testing.T) {
	var r render.Registry // zero value must work
	if got := r.Names(); len(got) != 0 {
		t.Errorf("empty Names() = %q, want []", got)
	}
	if err := r.Register("upper", upper); err != nil {
		t.Fatalf("Register(%q) error = %v, want <nil>", "upper", err)
	}
	got, err := r.Render("upper", "todo", []string{"a", "b"})
	if got != "TODO:a,b" || err != nil {
		t.Errorf("Render(%q, %q, [a b]) = %q, %v, want %q, <nil>", "upper", "todo", got, err, "TODO:a,b")
	}
}

func TestRegisterErrors(t *testing.T) {
	var r render.Registry
	if err := r.Register("x", upper); err != nil {
		t.Fatalf("Register(%q) error = %v, want <nil>", "x", err)
	}
	tests := []struct {
		name string
		f    render.Func
		want error
		msg  string
	}{
		{"", upper, render.ErrEmptyName, "render: empty name"},
		{"", nil, render.ErrEmptyName, "render: empty name"},
		{"y", nil, render.ErrNilFunc, "render: nil func"},
		{"x", upper, render.ErrDuplicate, `render: already registered: "x"`},
	}
	for _, tt := range tests {
		err := r.Register(tt.name, tt.f)
		if !errors.Is(err, tt.want) {
			t.Errorf("Register(%q) error = %v, want %v", tt.name, err, tt.want)
			continue
		}
		if err.Error() != tt.msg {
			t.Errorf("Register(%q) error = %q, want %q", tt.name, err.Error(), tt.msg)
		}
	}
	if got := r.Names(); len(got) != 1 || got[0] != "x" {
		t.Errorf("Names() after failed Registers = %q, want %q", got, []string{"x"})
	}
}

func TestRenderUnknown(t *testing.T) {
	var r render.Registry
	got, err := r.Render("nope", "T", nil)
	if got != "" || !errors.Is(err, render.ErrUnknown) {
		t.Fatalf("Render(%q) = %q, %v, want %q, %v", "nope", got, err, "", render.ErrUnknown)
	}
	if want := `render: unknown format: "nope"`; err.Error() != want {
		t.Errorf("Render(%q) error = %q, want %q", "nope", err.Error(), want)
	}
}

func TestNamesSorted(t *testing.T) {
	var r render.Registry
	for _, n := range []string{"zeta", "alpha", "Mid", "beta"} {
		if err := r.Register(n, upper); err != nil {
			t.Fatalf("Register(%q) error = %v, want <nil>", n, err)
		}
	}
	want := []string{"Mid", "alpha", "beta", "zeta"}
	got := r.Names()
	if strings.Join(got, ",") != strings.Join(want, ",") {
		t.Errorf("Names() = %q, want %q", got, want)
	}
}

func TestRegistriesAreIndependent(t *testing.T) {
	var a, b render.Registry
	if err := a.Register("only-a", upper); err != nil {
		t.Fatalf("Register(%q) error = %v, want <nil>", "only-a", err)
	}
	if _, err := b.Render("only-a", "T", nil); !errors.Is(err, render.ErrUnknown) {
		t.Errorf("second registry Render(%q) error = %v, want %v", "only-a", err, render.ErrUnknown)
	}
}

func TestOldGlobalAPIRemoved(t *testing.T) {
	f, err := parser.ParseFile(token.NewFileSet(), "render.go", nil, 0)
	if err != nil {
		t.Fatalf("parse render.go: %v", err)
	}
	for _, d := range f.Decls {
		switch d := d.(type) {
		case *ast.FuncDecl:
			if d.Recv == nil && (d.Name.Name == "Register" || d.Name.Name == "Render") {
				t.Errorf("package-level func %s = present, want removed (use Registry methods)", d.Name.Name)
			}
		case *ast.GenDecl:
			if d.Tok != token.VAR {
				continue
			}
			for _, s := range d.Specs {
				for _, v := range s.(*ast.ValueSpec).Values {
					if _, isLit := v.(*ast.CompositeLit); isLit {
						t.Errorf("package-level map var = present, want removed (state belongs in a Registry)")
					}
				}
			}
		}
	}
}

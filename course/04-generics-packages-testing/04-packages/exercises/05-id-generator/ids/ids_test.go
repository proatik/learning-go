package ids_test

import (
	"go/ast"
	"go/parser"
	"go/token"
	"strconv"
	"strings"
	"testing"

	"learngo/course/04-generics-packages-testing/04-packages/exercises/05-id-generator/ids"
)

func TestGenerator(t *testing.T) {
	g := ids.New("u-")
	if g == nil {
		t.Fatalf("New(%q) = nil, want a generator", "u-")
	}
	for i, want := range []string{"u-1", "u-2", "u-3"} {
		if got := g.Next(); got != want {
			t.Errorf("Next() call %d = %q, want %q", i+1, got, want)
		}
	}
	if got := g.Issued(); got != 3 {
		t.Errorf("Issued() = %d, want 3", got)
	}
}

func TestZeroGenerator(t *testing.T) {
	var g ids.Generator
	if got := g.Issued(); got != 0 {
		t.Errorf("zero Generator Issued() = %d, want 0", got)
	}
	for i, want := range []string{"1", "2"} {
		if got := g.Next(); got != want {
			t.Errorf("zero Generator Next() call %d = %q, want %q", i+1, got, want)
		}
	}
}

func TestGeneratorsAreIndependent(t *testing.T) {
	a, b := ids.New("x"), ids.New("x")
	if a == nil || b == nil {
		t.Fatalf("New(%q) = nil, want a generator", "x")
	}
	a.Next()
	a.Next()
	if got := b.Next(); got != "x1" {
		t.Errorf("second generator Next() = %q, want %q", got, "x1")
	}
	if got := a.Next(); got != "x3" {
		t.Errorf("first generator third Next() = %q, want %q", got, "x3")
	}
	// The deprecated package-level Next must not touch other generators.
	ids.Next()
	if got := b.Next(); got != "x2" {
		t.Errorf("second generator Next() after ids.Next() = %q, want %q", got, "x2")
	}
}

func TestDeprecatedNextStillWorks(t *testing.T) {
	first, err1 := strconv.Atoi(ids.Next())
	second, err2 := strconv.Atoi(ids.Next())
	if err1 != nil || err2 != nil || second != first+1 {
		t.Errorf("two ids.Next() calls = %d, %d, want consecutive numbers", first, second)
	}
}

func parseIDs(t *testing.T) *ast.File {
	t.Helper()
	f, err := parser.ParseFile(token.NewFileSet(), "ids.go", nil, parser.ParseComments)
	if err != nil {
		t.Fatalf("parse ids.go: %v", err)
	}
	return f
}

func TestNextIsDeprecated(t *testing.T) {
	for _, d := range parseIDs(t).Decls {
		fn, ok := d.(*ast.FuncDecl)
		if !ok || fn.Recv != nil || fn.Name.Name != "Next" {
			continue
		}
		doc := fn.Doc.Text()
		deprecated := false
		for para := range strings.SplitSeq(doc, "\n\n") {
			if strings.HasPrefix(para, "Deprecated: ") && strings.Contains(para, "Generator") {
				deprecated = true
			}
		}
		if !deprecated {
			t.Errorf("ids.Next doc comment = %q, want a paragraph starting with %q that mentions Generator", doc, "Deprecated: ")
		}
		return
	}
	t.Fatalf("ids.go has no func Next, want the deprecated package-level Next kept")
}

func TestNoIntCounter(t *testing.T) {
	for _, d := range parseIDs(t).Decls {
		gd, ok := d.(*ast.GenDecl)
		if !ok || gd.Tok != token.VAR {
			continue
		}
		for _, spec := range gd.Specs {
			vs := spec.(*ast.ValueSpec)
			if id, ok := vs.Type.(*ast.Ident); ok && id.Name == "int" {
				t.Errorf("package-level var %s int = present, want removed (use a default Generator)", vs.Names[0].Name)
			}
		}
	}
}

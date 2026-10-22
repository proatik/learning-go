package grades_test

import (
	"go/ast"
	"go/parser"
	"go/token"
	"slices"
	"testing"

	grades "learngo/course/04-generics-packages-testing/04-packages/exercises/01-split-package"
)

func TestReport(t *testing.T) {
	tests := []struct {
		in   []int
		want string
	}{
		{nil, "n=0"},
		{[]int{}, "n=0"},
		{[]int{70, 80, 95}, "n=3 mean=81.67 min=70 max=95"},
		{[]int{100}, "n=1 mean=100.00 min=100 max=100"},
		{[]int{0, 0}, "n=2 mean=0.00 min=0 max=0"},
		{[]int{-5, 5, 10, 2}, "n=4 mean=3.00 min=-5 max=10"},
	}
	for _, tt := range tests {
		if got := grades.Report(tt.in); got != tt.want {
			t.Errorf("Report(%v) = %q, want %q", tt.in, got, tt.want)
		}
	}
}

// TestHelpersMoved checks the package structure: grades.go keeps only Report.
func TestHelpersMoved(t *testing.T) {
	f, err := parser.ParseFile(token.NewFileSet(), "grades.go", nil, 0)
	if err != nil {
		t.Fatalf("parse grades.go: %v", err)
	}
	var funcs []string
	for _, d := range f.Decls {
		if fn, ok := d.(*ast.FuncDecl); ok {
			funcs = append(funcs, fn.Name.Name)
		}
	}
	if want := []string{"Report"}; !slices.Equal(funcs, want) {
		t.Errorf("functions in grades.go = %v, want %v", funcs, want)
	}
}

package build

import (
	"errors"
	"slices"
	"testing"

	"learngo/course/04-generics-packages-testing/06-module-challenge/exercises/07-build-graph/graph"
)

var project = []Target{
	{Name: "app", Deps: []string{"lib", "ui"}},
	{Name: "ui", Deps: []string{"lib", "assets"}},
	{Name: "lib", Deps: []string{"core"}},
	{Name: "core"},
	{Name: "assets"},
	{Name: "docs", Deps: []string{"core"}},
	{Name: "tests", Deps: []string{"app", "lib"}},
}

func TestPlan(t *testing.T) {
	tests := []struct {
		name    string
		targets []Target
		want    []string
	}{
		{"project", project, []string{"assets", "core", "docs", "lib", "ui", "app", "tests"}},
		{"independent targets are alphabetical", []Target{{Name: "c"}, {Name: "a"}, {Name: "b"}}, []string{"a", "b", "c"}},
		{"chain", []Target{{Name: "z", Deps: []string{"y"}}, {Name: "y", Deps: []string{"x"}}, {Name: "x"}}, []string{"x", "y", "z"}},
		{"repeated dependency", []Target{{Name: "b", Deps: []string{"a", "a"}}, {Name: "a"}}, []string{"a", "b"}},
		{"nothing", nil, []string{}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Plan(tt.targets)
			if err != nil || !slices.Equal(got, tt.want) {
				t.Errorf("Plan(%v) = %q, %v, want %q, <nil>", tt.targets, got, err, tt.want)
			}
		})
	}
}

func TestPlanUnknownDep(t *testing.T) {
	targets := []Target{{Name: "app", Deps: []string{"lib"}}, {Name: "lib", Deps: []string{"core"}}}
	got, err := Plan(targets)
	if got != nil || !errors.Is(err, ErrUnknownDep) {
		t.Fatalf("Plan(lib needs missing core) = %q, %v, want nil, %v", got, err, ErrUnknownDep)
	}
	if want := `build: unknown dependency: "lib" needs "core"`; err.Error() != want {
		t.Errorf("Plan(lib needs missing core) error = %q, want %q", err.Error(), want)
	}
}

func TestPlanCycle(t *testing.T) {
	targets := []Target{
		{Name: "a", Deps: []string{"b"}},
		{Name: "b", Deps: []string{"a"}},
		{Name: "c"},
	}
	got, err := Plan(targets)
	if got != nil || !errors.Is(err, graph.ErrCycle) {
		t.Fatalf("Plan(a ⇄ b) = %q, %v, want nil, %v", got, err, graph.ErrCycle)
	}
	var ce *graph.CycleError[string]
	if !errors.As(err, &ce) || !slices.Equal(ce.Nodes, []string{"a", "b"}) {
		t.Errorf("errors.As(%v, *graph.CycleError[string]) Nodes = %v, want [a b] (wrap with %%w)", err, ce)
	}
	if want := "build: graph: cycle among [a b]"; err.Error() != want {
		t.Errorf("Plan(a ⇄ b) error = %q, want %q", err.Error(), want)
	}
}

func TestAffected(t *testing.T) {
	tests := []struct {
		changed string
		want    []string
	}{
		{"core", []string{"app", "core", "docs", "lib", "tests", "ui"}},
		{"assets", []string{"app", "assets", "tests", "ui"}},
		{"tests", []string{"tests"}},
		{"unknown", []string{}},
	}
	for _, tt := range tests {
		got, err := Affected(project, tt.changed)
		if err != nil || got == nil || !slices.Equal(got, tt.want) {
			t.Errorf("Affected(project, %q) = %#v, %v, want %q, <nil>", tt.changed, got, err, tt.want)
		}
	}
	if _, err := Affected([]Target{{Name: "x", Deps: []string{"y"}}}, "x"); !errors.Is(err, ErrUnknownDep) {
		t.Errorf("Affected(x needs missing y) error = %v, want %v", err, ErrUnknownDep)
	}
}

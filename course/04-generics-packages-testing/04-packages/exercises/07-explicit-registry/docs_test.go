package docs_test

import (
	"errors"
	"testing"

	docs "learngo/course/04-generics-packages-testing/04-packages/exercises/07-explicit-registry"
	"learngo/course/04-generics-packages-testing/04-packages/exercises/07-explicit-registry/render"
)

func TestDefault(t *testing.T) {
	r, err := docs.Default()
	if err != nil || r == nil {
		t.Fatalf("Default() = %v, %v, want a registry, <nil>", r, err)
	}
	got := r.Names()
	want := []string{"markdown", "numbered", "plain"}
	if len(got) != len(want) || got[0] != want[0] || got[1] != want[1] || got[2] != want[2] {
		t.Errorf("Default().Names() = %q, want %q", got, want)
	}
	tests := []struct {
		format string
		want   string
	}{
		{"plain", "Go: maps, slices"},
		{"markdown", "# Go\n\n- maps\n- slices"},
		{"numbered", "Go\n1. maps\n2. slices"},
	}
	for _, tt := range tests {
		out, err := r.Render(tt.format, "Go", []string{"maps", "slices"})
		if out != tt.want || err != nil {
			t.Errorf("Render(%q, %q, [maps slices]) = %q, %v, want %q, <nil>", tt.format, "Go", out, err, tt.want)
		}
	}
}

func TestDefaultIsFresh(t *testing.T) {
	a, err1 := docs.Default()
	b, err2 := docs.Default()
	if a == nil || b == nil || err1 != nil || err2 != nil {
		t.Fatalf("Default() twice = %v, %v / %v, %v, want two registries", a, err1, b, err2)
	}
	if a == b {
		t.Fatalf("Default() returned the same *Registry twice, want a new one per call")
	}
	if err := a.Register("extra", func(title string, _ []string) string { return title }); err != nil {
		t.Fatalf("Register(%q) error = %v, want <nil>", "extra", err)
	}
	if _, err := b.Render("extra", "T", nil); !errors.Is(err, render.ErrUnknown) {
		t.Errorf("other registry Render(%q) error = %v, want %v", "extra", err, render.ErrUnknown)
	}
}

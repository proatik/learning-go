package wordcount

import (
	"maps"
	"testing"
)

func TestWordCount(t *testing.T) {
	tests := []struct {
		in   string
		want map[string]int
	}{
		{"the cat the hat", map[string]int{"the": 2, "cat": 1, "hat": 1}},
		{"Go go GO", map[string]int{"go": 3}},
		{"  spaced\tout\n words  ", map[string]int{"spaced": 1, "out": 1, "words": 1}},
		{"one", map[string]int{"one": 1}},
		{"go! go Go!", map[string]int{"go!": 2, "go": 1}},
		{"a b a b a", map[string]int{"a": 3, "b": 2}},
		{"Héllo HÉLLO wörld", map[string]int{"héllo": 2, "wörld": 1}},
		{"", map[string]int{}},
		{" \t\n ", map[string]int{}},
	}
	for _, tt := range tests {
		got := WordCount(tt.in)
		if got == nil {
			t.Errorf("WordCount(%q) = nil, want an initialized map", tt.in)
			continue
		}
		if !maps.Equal(got, tt.want) {
			t.Errorf("WordCount(%q) = %v, want %v", tt.in, got, tt.want)
		}
	}
}

func TestWordCountResultIsWritable(t *testing.T) {
	got := WordCount("")
	if got == nil {
		t.Fatal("WordCount(\"\") = nil; writing to it would panic")
	}
	got["new"] = 1 // must not panic
}

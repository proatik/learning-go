package topk

import (
	"maps"
	"slices"
	"testing"
)

func TestCounts(t *testing.T) {
	tests := []struct {
		text string
		want map[string]int
	}{
		{"Go go GO! gopher", map[string]int{"go": 3, "gopher": 1}},
		{"don't stop", map[string]int{"don": 1, "t": 1, "stop": 1}},
		{"a1 a1, b2;b2 b2", map[string]int{"a1": 2, "b2": 3}},
		{"Ünïcode ünïcode café", map[string]int{"ünïcode": 2, "café": 1}},
		{"  \t\n ", map[string]int{}},
		{"", map[string]int{}},
		{"...!!!", map[string]int{}},
	}
	for _, tt := range tests {
		got := Counts(tt.text)
		if got == nil {
			t.Errorf("Counts(%q) = nil, want a non-nil map", tt.text)
			continue
		}
		if !maps.Equal(got, tt.want) {
			t.Errorf("Counts(%q) = %v, want %v", tt.text, got, tt.want)
		}
	}
}

func TestTopK(t *testing.T) {
	tests := []struct {
		text string
		k    int
		want []WordCount
	}{
		{"b a b c a b", 2, []WordCount{{"b", 3}, {"a", 2}}},
		{"b a b c a b", 3, []WordCount{{"b", 3}, {"a", 2}, {"c", 1}}},
		{"b a b c a b", 10, []WordCount{{"b", 3}, {"a", 2}, {"c", 1}}},
		{"x y z", 2, []WordCount{{"x", 1}, {"y", 1}}},
		{"z y x z y", 2, []WordCount{{"y", 2}, {"z", 2}}},
		{"one", 5, []WordCount{{"one", 1}}},
		{"The cat and the hat. THE END, and done.", 3, []WordCount{{"the", 3}, {"and", 2}, {"cat", 1}}},
		{"apple Apple APPLE banana Banana cherry", 1, []WordCount{{"apple", 3}}},
		{"", 3, []WordCount{}},
	}
	for _, tt := range tests {
		got := TopK(tt.text, tt.k)
		if !slices.Equal(got, tt.want) {
			t.Errorf("TopK(%q, %d) = %v, want %v", tt.text, tt.k, got, tt.want)
		}
	}
}

func TestTopKNonPositive(t *testing.T) {
	for _, k := range []int{0, -1, -100} {
		if got := TopK("a b c", k); len(got) != 0 {
			t.Errorf("TopK(\"a b c\", %d) = %v, want empty", k, got)
		}
	}
}

func TestTopKLargeTieBlock(t *testing.T) {
	// 26 words each appearing once, plus "q" once more.
	text := "q "
	for c := 'a'; c <= 'z'; c++ {
		text += string(c) + " "
	}
	got := TopK(text, 4)
	want := []WordCount{{"q", 2}, {"a", 1}, {"b", 1}, {"c", 1}}
	if !slices.Equal(got, want) {
		t.Fatalf("TopK(alphabet+q, 4) = %v, want %v", got, want)
	}
}

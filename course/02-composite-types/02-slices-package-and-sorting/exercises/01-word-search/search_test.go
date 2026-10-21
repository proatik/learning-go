package search

import "testing"

func TestHasWord(t *testing.T) {
	tests := []struct {
		words []string
		w     string
		want  bool
	}{
		{[]string{"go", "zig"}, "zig", true},
		{[]string{"go", "zig"}, "go", true},
		{[]string{"go", "zig"}, "rust", false},
		{[]string{"Go"}, "go", false},
		{[]string{"go"}, "", false},
		{[]string{""}, "", true},
		{nil, "go", false},
	}
	for _, tt := range tests {
		if got := HasWord(tt.words, tt.w); got != tt.want {
			t.Errorf("HasWord(%q, %q) = %v, want %v", tt.words, tt.w, got, tt.want)
		}
	}
}

func TestPosition(t *testing.T) {
	tests := []struct {
		words []string
		w     string
		want  int
	}{
		{[]string{"a", "b", "a"}, "a", 0},
		{[]string{"a", "b", "a"}, "b", 1},
		{[]string{"x", "y", "z"}, "z", 2},
		{[]string{"a", "b"}, "c", -1},
		{[]string{"A"}, "a", -1},
		{nil, "a", -1},
	}
	for _, tt := range tests {
		if got := Position(tt.words, tt.w); got != tt.want {
			t.Errorf("Position(%q, %q) = %d, want %d", tt.words, tt.w, got, tt.want)
		}
	}
}

func TestFirstLonger(t *testing.T) {
	tests := []struct {
		words []string
		n     int
		want  int
	}{
		{[]string{"go", "rust", "haskell"}, 3, 1},
		{[]string{"go", "rust", "haskell"}, 4, 2},
		{[]string{"go", "rust", "haskell"}, 1, 0},
		{[]string{"go", "rust", "haskell"}, 7, -1},
		{[]string{"go"}, 5, -1},
		{[]string{"", "a"}, 0, 1},
		{[]string{"héllo"}, 5, 0}, // 6 bytes: é is 2 bytes in UTF-8
		{nil, 0, -1},
	}
	for _, tt := range tests {
		if got := FirstLonger(tt.words, tt.n); got != tt.want {
			t.Errorf("FirstLonger(%q, %d) = %d, want %d", tt.words, tt.n, got, tt.want)
		}
	}
}

func TestAnyEmpty(t *testing.T) {
	tests := []struct {
		words []string
		want  bool
	}{
		{[]string{"a", "", "b"}, true},
		{[]string{""}, true},
		{[]string{"a", " ", "b"}, false},
		{[]string{"a"}, false},
		{nil, false},
	}
	for _, tt := range tests {
		if got := AnyEmpty(tt.words); got != tt.want {
			t.Errorf("AnyEmpty(%q) = %v, want %v", tt.words, got, tt.want)
		}
	}
}

func TestSameWords(t *testing.T) {
	tests := []struct {
		a, b []string
		want bool
	}{
		{[]string{"a", "b"}, []string{"a", "b"}, true},
		{[]string{"a", "b"}, []string{"b", "a"}, false},
		{[]string{"a"}, []string{"a", "a"}, false},
		{[]string{"a", "b"}, []string{"a"}, false},
		{[]string{"Go"}, []string{"go"}, false},
		{nil, []string{}, true},
		{nil, nil, true},
		{[]string{""}, nil, false},
	}
	for _, tt := range tests {
		if got := SameWords(tt.a, tt.b); got != tt.want {
			t.Errorf("SameWords(%q, %q) = %v, want %v", tt.a, tt.b, got, tt.want)
		}
	}
}

package wrap

import (
	"fmt"
	"testing"
)

func TestWrap(t *testing.T) {
	tests := []struct {
		text  string
		width int
		want  string
	}{
		{"The quick brown fox jumps over the lazy dog", 10, "The quick\nbrown fox\njumps over\nthe lazy\ndog"},
		{"a b c d", 3, "a b\nc d"},
		{"a b c d", 1, "a\nb\nc\nd"},
		{"a b c d", 7, "a b c d"},
		{"a b c d", 6, "a b c\nd"},
		{"  lots   of\n\tspace  ", 20, "lots of space"},
		{"supercalifragilistic is long", 10, "supercalifragilistic\nis long"},
		{"is supercalifragilistic long", 10, "is\nsupercalifragilistic\nlong"},
		{"héllo wörld", 5, "héllo\nwörld"},
		{"héllo wörld", 11, "héllo wörld"},
		{"世界 世界 世界", 5, "世界 世界\n世界"},
		{"one", 100, "one"},
		{"", 10, ""},
		{" \n\t ", 10, ""},
	}
	for _, tt := range tests {
		if got := Wrap(tt.text, tt.width); got != tt.want {
			t.Errorf("Wrap(%q, %d) = %q, want %q", tt.text, tt.width, got, tt.want)
		}
	}
}

func ExampleWrap() {
	fmt.Println(Wrap("Go is an open source programming language that makes it simple to build secure, scalable systems.", 24))
	// Output:
	// Go is an open source
	// programming language
	// that makes it simple to
	// build secure, scalable
	// systems.
}

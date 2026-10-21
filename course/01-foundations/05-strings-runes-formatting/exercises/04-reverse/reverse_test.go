package reverse

import (
	"testing"
	"unicode/utf8"
)

func TestReverse(t *testing.T) {
	tests := []struct {
		s    string
		want string
	}{
		{"Hello", "olleH"},
		{"héllo, 世界", "界世 ,olléh"},
		{"🙂👍", "👍🙂"},
		{"a", "a"},
		{"ab", "ba"},
		{"racecar", "racecar"},
		{"Go 1.27!", "!72.1 oG"},
		{"", ""},
	}
	for _, tt := range tests {
		got := Reverse(tt.s)
		if got != tt.want {
			t.Errorf("Reverse(%q) = %q, want %q", tt.s, got, tt.want)
		}
		if !utf8.ValidString(got) {
			t.Errorf("Reverse(%q) returned invalid UTF-8 %q", tt.s, got)
		}
	}
}

func TestReverseTwiceIsIdentity(t *testing.T) {
	for _, s := range []string{"héllo, 世界", "naïve café", "日本語"} {
		if got := Reverse(Reverse(s)); got != s {
			t.Errorf("Reverse(Reverse(%q)) = %q, want %q", s, got, s)
		}
	}
}

func TestReverseWords(t *testing.T) {
	tests := []struct {
		s    string
		want string
	}{
		{"one two three", "three two one"},
		{"  Go   is\tfun ", "fun is Go"},
		{"世界 héllo", "héllo 世界"},
		{"single", "single"},
		{"a\nb\nc", "c b a"},
		{"", ""},
		{"   ", ""},
	}
	for _, tt := range tests {
		if got := ReverseWords(tt.s); got != tt.want {
			t.Errorf("ReverseWords(%q) = %q, want %q", tt.s, got, tt.want)
		}
	}
}

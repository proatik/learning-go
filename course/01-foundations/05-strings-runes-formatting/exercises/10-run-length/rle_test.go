package rle

import "testing"

func TestEncode(t *testing.T) {
	tests := []struct {
		s    string
		want string
	}{
		{"aaabcc", "a3b1c2"},
		{"abc", "a1b1c1"},
		{"wwwwwwwwwwwwb", "w12b1"},
		{"ééé世世", "é3世2"},
		{"aabbaa", "a2b2a2"},
		{"a", "a1"},
		{"  !!!", " 2!3"},
		{"🙂🙂🙂🙂", "🙂4"},
		{"AAaa", "A2a2"},
		{"", ""},
	}
	for _, tt := range tests {
		if got := Encode(tt.s); got != tt.want {
			t.Errorf("Encode(%q) = %q, want %q", tt.s, got, tt.want)
		}
	}
}

func TestDecode(t *testing.T) {
	tests := []struct {
		s    string
		want string
	}{
		{"a3b1c2", "aaabcc"},
		{"w12b1", "wwwwwwwwwwwwb"},
		{"é3世2", "ééé世世"},
		{" 2!3", "  !!!"},
		{"a1", "a"},
		{"x007", "xxxxxxx"},
		{"", ""},
	}
	for _, tt := range tests {
		got, err := Decode(tt.s)
		if err != nil || got != tt.want {
			t.Errorf("Decode(%q) = (%q, %v), want (%q, nil)", tt.s, got, err, tt.want)
		}
	}
}

func TestDecodeErrors(t *testing.T) {
	for _, s := range []string{"3a", "a3b", "a0", "ab2", "a", "é", "a2 ", "12", "a00"} {
		got, err := Decode(s)
		if err == nil || got != "" {
			t.Errorf("Decode(%q) = (%q, %v), want (\"\", error)", s, got, err)
		}
	}
}

func TestRoundTrip(t *testing.T) {
	for _, s := range []string{"aaabcc", "héllo, wörld!!", "世世界界界", "   ", "zzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzz"} {
		got, err := Decode(Encode(s))
		if err != nil || got != s {
			t.Errorf("Decode(Encode(%q)) = (%q, %v), want (%q, nil)", s, got, err, s)
		}
	}
}

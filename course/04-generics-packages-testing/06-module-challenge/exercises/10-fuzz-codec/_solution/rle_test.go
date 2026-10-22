package rle

import (
	"errors"
	"fmt"
	"strings"
	"testing"
	"unicode/utf8"
)

// FuzzRoundTrip checks that Decode(Encode(s)) == s.
func FuzzRoundTrip(f *testing.F) {
	seeds := []string{
		"", "a", "aaab", "abc",
		strings.Repeat("x", 12), // a run of 10+ needs a two-digit count
		"aAAa",                  // runs that differ only in case
		"héé",                   // a multi-byte rune
		"::a:",                  // the separator as text
		"112233",                // digits as text
		"日本日本",
	}
	for _, s := range seeds {
		f.Add(s)
	}
	f.Fuzz(func(t *testing.T, s string) {
		if !utf8.ValidString(s) {
			t.Skip("Encode is defined for valid UTF-8 only")
		}
		enc := Encode(s)
		got, err := Decode(enc)
		if err != nil || got != s {
			t.Errorf("Decode(Encode(%q)) = %q, %v, want %q, <nil> (Encode gave %q)", s, got, err, s, enc)
		}
	})
}

// FuzzDecode feeds Decode arbitrary input: it must not panic, and whatever
// it accepts must survive a round trip.
func FuzzDecode(f *testing.F) {
	for _, s := range []string{"", "3:a", "3:", ":", "x", "2:a3", "1:\xff", "99999999999:a"} {
		f.Add(s)
	}
	f.Fuzz(func(t *testing.T, s string) {
		out, err := Decode(s)
		if err != nil {
			if !errors.Is(err, ErrMalformed) {
				t.Errorf("Decode(%q) error = %v, want %v", s, err, ErrMalformed)
			}
			return
		}
		if again, err := Decode(Encode(out)); err != nil || again != out {
			t.Errorf("Decode(Encode(%q)) = %q, %v, want %q, <nil>", out, again, err, out)
		}
	})
}

func TestDecodeErrors(t *testing.T) {
	tests := []string{"a", ":a", "3:", "0:a", "-1:a", "x:a", "2:a3", "1:\xff", "1000001:a"}
	for _, in := range tests {
		t.Run(fmt.Sprintf("%q", in), func(t *testing.T) {
			if got, err := Decode(in); !errors.Is(err, ErrMalformed) {
				t.Errorf("Decode(%q) = %q, %v, want error %v", in, got, err, ErrMalformed)
			}
		})
	}
}

func ExampleEncode() {
	fmt.Println(Encode("aaab"))
	fmt.Println(Encode("héé"))
	s, err := Decode("2:x1::")
	fmt.Println(s, err)
	// Output:
	// 3:a1:b
	// 1:h2:é
	// xx: <nil>
}

func BenchmarkEncode(b *testing.B) {
	s := strings.Repeat("aaabccdddd", 100)
	for b.Loop() {
		Encode(s)
	}
}

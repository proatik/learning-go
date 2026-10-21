package textstats

import (
	"bytes"
	"errors"
	"io"
	"strings"
	"testing"
)

var cases = []struct {
	in      string
	lines   int
	words   int
	longest string
}{
	{"", 0, 0, ""},
	{"one two\nthree", 2, 3, "one two"},
	{"a\n", 1, 1, "a"},
	{"\n\n", 2, 0, ""},
	{"  hello   world\n\tgo ", 2, 3, "  hello   world"},
	{"hi\r\nthere\r\n", 2, 2, "there"},
	{"日本語\nabcd", 2, 2, "abcd"},
	{"abc\nxyz\nab", 3, 3, "abc"},
	{"no newline at all", 1, 4, "no newline at all"},
	{"a b c", 1, 3, "a b c"},
	{"x\n\n\ny\n", 4, 2, "x"},
	{"   \n\t\n", 2, 0, "   "},
	{"ééé\nabcd", 2, 2, "abcd"},
	{"ééééé\nabcd", 2, 2, "ééééé"},
}

func TestLineCount(t *testing.T) {
	for _, tt := range cases {
		got, err := LineCount(strings.NewReader(tt.in))
		if err != nil {
			t.Errorf("LineCount(%q) error = %v, want nil", tt.in, err)
		}
		if got != tt.lines {
			t.Errorf("LineCount(%q) = %d, want %d", tt.in, got, tt.lines)
		}
	}
}

func TestWordCount(t *testing.T) {
	for _, tt := range cases {
		got, err := WordCount(strings.NewReader(tt.in))
		if err != nil {
			t.Errorf("WordCount(%q) error = %v, want nil", tt.in, err)
		}
		if got != tt.words {
			t.Errorf("WordCount(%q) = %d, want %d", tt.in, got, tt.words)
		}
	}
}

func TestLongest(t *testing.T) {
	for _, tt := range cases {
		got, err := Longest(strings.NewReader(tt.in))
		if err != nil {
			t.Errorf("Longest(%q) error = %v, want nil", tt.in, err)
		}
		if got != tt.longest {
			t.Errorf("Longest(%q) = %q, want %q", tt.in, got, tt.longest)
		}
	}
}

func TestOtherReaders(t *testing.T) {
	// Any io.Reader works: a bytes.Buffer, or several readers glued together.
	buf := bytes.NewBufferString("alpha beta\ngamma\n")
	if got, _ := WordCount(buf); got != 3 {
		t.Errorf("WordCount(bytes.Buffer %q) = %d, want 3", "alpha beta\ngamma\n", got)
	}
	mr := io.MultiReader(strings.NewReader("first li"), strings.NewReader("ne\nsecond\n"))
	if got, _ := LineCount(mr); got != 2 {
		t.Errorf("LineCount(MultiReader %q + %q) = %d, want 2", "first li", "ne\nsecond\n", got)
	}
	long := strings.Repeat("x", 5000) + "\nshort\n"
	if got, _ := Longest(strings.NewReader(long)); len(got) != 5000 {
		t.Errorf("Longest(5000 x's + %q) has length %d, want 5000", "\nshort\n", len(got))
	}
}

var errDisk = errors.New("disk error")

// failingReader returns its data, then fails instead of reporting io.EOF.
type failingReader struct {
	data string
	done bool
}

func (f *failingReader) Read(p []byte) (int, error) {
	if f.done {
		return 0, errDisk
	}
	f.done = true
	return copy(p, f.data), nil
}

func TestReaderError(t *testing.T) {
	in := "one two\nthree\n"
	if _, err := LineCount(&failingReader{data: in}); !errors.Is(err, errDisk) {
		t.Errorf("LineCount(reader failing after %q) error = %v, want %v", in, err, errDisk)
	}
	if _, err := WordCount(&failingReader{data: in}); !errors.Is(err, errDisk) {
		t.Errorf("WordCount(reader failing after %q) error = %v, want %v", in, err, errDisk)
	}
	if _, err := Longest(&failingReader{data: in}); !errors.Is(err, errDisk) {
		t.Errorf("Longest(reader failing after %q) error = %v, want %v", in, err, errDisk)
	}
}

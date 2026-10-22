package lines_test

import (
	"errors"
	"fmt"
	"io"
	"iter"
	"slices"
	"strings"
	"testing"

	lines "learngo/course/04-generics-packages-testing/03-iterators/exercises/07-line-reader"
)

// collect2 ranges over seq and returns each pair formatted as "k:v",
// breaking after max pairs when max > 0 (and always after 10000). A panic
// during the loop is returned instead of crashing.
func collect2[K, V any](seq iter.Seq2[K, V], max int) (got []string, panicked any) {
	defer func() { panicked = recover() }()
	for k, v := range seq {
		got = append(got, fmt.Sprintf("%v:%v", k, v))
		if len(got) == max || len(got) == 10000 {
			break
		}
	}
	return got, nil
}

// oneLineReader returns one line of its input per Read call and counts
// the calls.
type oneLineReader struct {
	lines []string
	reads int
	err   error // returned once the lines run out (io.EOF if nil)
}

func (r *oneLineReader) Read(p []byte) (int, error) {
	r.reads++
	if len(r.lines) == 0 {
		if r.err != nil {
			return 0, r.err
		}
		return 0, io.EOF
	}
	n := copy(p, r.lines[0])
	r.lines[0] = r.lines[0][n:]
	if r.lines[0] == "" {
		r.lines = r.lines[1:]
	}
	return n, nil
}

func TestAll(t *testing.T) {
	tests := []struct {
		in   string
		want []string
	}{
		{"host = a\n\n  # comment\n  port = 80  \n", []string{"1:host = a", "4:port = 80"}},
		{"", nil},
		{"\n\n\n", nil},
		{"a\r\nb", []string{"1:a", "2:b"}},
		{"#only\n# comments", nil},
		{"x # not a comment\n\t#tabbed\n日本", []string{"1:x # not a comment", "3:日本"}},
		{"one", []string{"1:one"}},
	}
	for _, tt := range tests {
		lr := lines.NewLineReader(strings.NewReader(tt.in))
		seq := lr.All()
		if seq == nil {
			t.Fatalf("NewLineReader(%q).All() = nil, want a non-nil iterator", tt.in)
		}
		got, p := collect2(seq, -1)
		if p != nil {
			t.Fatalf("range over NewLineReader(%q).All() panicked: %v", tt.in, p)
		}
		if !slices.Equal(got, tt.want) {
			t.Errorf("NewLineReader(%q).All() = %q, want %q", tt.in, got, tt.want)
		}
		if err := lr.Err(); err != nil {
			t.Errorf("NewLineReader(%q) Err() = %v, want nil", tt.in, err)
		}
	}
}

func TestAllIsLazy(t *testing.T) {
	src := make([]string, 1000)
	for i := range src {
		src[i] = fmt.Sprintf("line %d\n", i+1)
	}
	r := &oneLineReader{lines: src}
	lr := lines.NewLineReader(r)
	seq := lr.All()
	if r.reads != 0 {
		t.Errorf("All() before ranging: Read calls = %d, want 0", r.reads)
	}
	got, p := collect2(seq, 2)
	if p != nil {
		t.Fatalf("break after 2 lines of a 1000-line input panicked: %v (stop when yield returns false)", p)
	}
	if want := []string{"1:line 1", "2:line 2"}; !slices.Equal(got, want) {
		t.Errorf("first 2 pairs of a 1000-line input = %q, want %q", got, want)
	}
	if r.reads > 3 {
		t.Errorf("break after 2 lines of a 1000-line input (one line per Read): Read calls = %d, want at most 3", r.reads)
	}
}

func TestSingleUse(t *testing.T) {
	in := "host = a\n\n  # comment\n  port = 80  \nuser = me\n"
	lr := lines.NewLineReader(strings.NewReader(in))
	first, p := collect2(lr.All(), 1)
	if p != nil {
		t.Fatalf("break after 1 pair panicked: %v (stop when yield returns false)", p)
	}
	rest, p := collect2(lr.All(), -1)
	if p != nil {
		t.Fatalf("second range over All() panicked: %v", p)
	}
	if want := []string{"1:host = a"}; !slices.Equal(first, want) {
		t.Errorf("first range over NewLineReader(%q).All() with a break after 1 = %q, want %q", in, first, want)
	}
	if want := []string{"4:port = 80", "5:user = me"}; !slices.Equal(rest, want) {
		t.Errorf("second range over the same All() = %q, want %q (continue the numbering)", rest, want)
	}
	again, _ := collect2(lr.All(), -1)
	if len(again) != 0 {
		t.Errorf("third range after the input ended = %q, want []", again)
	}
}

var errDisk = errors.New("disk on fire")

func TestReadError(t *testing.T) {
	r := &oneLineReader{lines: []string{"x\n", "\n", "y\n"}, err: errDisk}
	lr := lines.NewLineReader(r)
	got, p := collect2(lr.All(), -1)
	if p != nil {
		t.Fatalf("range over a failing reader panicked: %v", p)
	}
	if want := []string{"1:x", "3:y"}; !slices.Equal(got, want) {
		t.Errorf("All() on a reader failing after %q = %q, want %q", "x\n\ny\n", got, want)
	}
	if err := lr.Err(); !errors.Is(err, errDisk) {
		t.Errorf("Err() after a failing reader = %v, want %v", err, errDisk)
	}
}

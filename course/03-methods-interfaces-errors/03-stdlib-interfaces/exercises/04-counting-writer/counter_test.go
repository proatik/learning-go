package counter

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"strings"
	"testing"
)

var _ io.Writer = (*CountingWriter)(nil)

func TestCountingWriterForwards(t *testing.T) {
	tests := []struct {
		writes []string
		want   int64
	}{
		{[]string{"abc", "de"}, 5},
		{nil, 0},
		{[]string{""}, 0},
		{[]string{"héllo"}, 6},
		{[]string{"line 1\n", "line 2\n", "\n"}, 15},
		{[]string{"日本語"}, 9},
	}
	for _, tt := range tests {
		var buf bytes.Buffer
		c := &CountingWriter{W: &buf}
		for _, s := range tt.writes {
			n, err := c.Write([]byte(s))
			if n != len(s) || err != nil {
				t.Errorf("Write(%q) = %d, %v, want %d, <nil>", s, n, err, len(s))
			}
		}
		if c.N != tt.want {
			t.Errorf("N after writes %q = %d, want %d", tt.writes, c.N, tt.want)
		}
		if got, want := buf.String(), strings.Join(tt.writes, ""); got != want {
			t.Errorf("wrapped writer after writes %q = %q, want %q", tt.writes, got, want)
		}
	}
}

func TestCountingWriterWithFprintf(t *testing.T) {
	var sb strings.Builder
	c := &CountingWriter{W: &sb}
	fmt.Fprintf(c, "hi %d", 42)
	fmt.Fprintln(c, "!")
	if got, want := sb.String(), "hi 42!\n"; got != want {
		t.Errorf("Fprintf through CountingWriter = %q, want %q", got, want)
	}
	if c.N != 7 {
		t.Errorf("N after Fprintf(\"hi %%d\", 42) + Fprintln(\"!\") = %d, want 7", c.N)
	}
}

func TestCountingWriterNilWriter(t *testing.T) {
	c := &CountingWriter{}
	n, err := io.WriteString(c, "héllo")
	if n != 6 || err != nil {
		t.Errorf("WriteString(nil-W counter, %q) = %d, %v, want 6, <nil>", "héllo", n, err)
	}
	io.WriteString(c, "xyz")
	if c.N != 9 {
		t.Errorf("N after %q and %q with nil W = %d, want 9", "héllo", "xyz", c.N)
	}
}

var errBroken = errors.New("broken pipe")

// shortWriter accepts at most max bytes per call and then reports an error.
type shortWriter struct {
	max int
	got []byte
}

func (s *shortWriter) Write(p []byte) (int, error) {
	if len(p) <= s.max {
		s.got = append(s.got, p...)
		return len(p), nil
	}
	s.got = append(s.got, p[:s.max]...)
	return s.max, errBroken
}

func TestCountingWriterShortWrite(t *testing.T) {
	sw := &shortWriter{max: 3}
	c := &CountingWriter{W: sw}
	n, err := c.Write([]byte("hello"))
	if n != 3 || !errors.Is(err, errBroken) {
		t.Errorf("Write(%q) on a 3-byte writer = %d, %v, want 3, %v", "hello", n, err, errBroken)
	}
	if c.N != 3 {
		t.Errorf("N after short Write(%q) = %d, want 3", "hello", c.N)
	}
	c.Write([]byte("ok"))
	if c.N != 5 {
		t.Errorf("N after Write(%q) = %d, want 5", "ok", c.N)
	}
}

func TestMeasure(t *testing.T) {
	var buf bytes.Buffer
	n, err := Measure(&buf, func(w io.Writer) error {
		_, err := fmt.Fprintf(w, "name=%s age=%d\n", "Ann", 31)
		return err
	})
	if n != 16 || err != nil {
		t.Errorf("Measure(Fprintf %q) = %d, %v, want 16, <nil>", "name=Ann age=31\n", n, err)
	}
	if got, want := buf.String(), "name=Ann age=31\n"; got != want {
		t.Errorf("Measure wrote %q, want %q", got, want)
	}

	n, err = Measure(&buf, func(w io.Writer) error { return nil })
	if n != 0 || err != nil {
		t.Errorf("Measure(no writes) = %d, %v, want 0, <nil>", n, err)
	}

	errStop := errors.New("stop")
	n, err = Measure(io.Discard, func(w io.Writer) error {
		io.WriteString(w, "1234")
		return errStop
	})
	if n != 4 || !errors.Is(err, errStop) {
		t.Errorf("Measure(write 4 bytes then fail) = %d, %v, want 4, %v", n, err, errStop)
	}

	n, err = Measure(&shortWriter{max: 2}, func(w io.Writer) error {
		_, err := io.WriteString(w, "abcdef")
		return err
	})
	if n != 2 || !errors.Is(err, errBroken) {
		t.Errorf("Measure(6 bytes to a 2-byte writer) = %d, %v, want 2, %v", n, err, errBroken)
	}
}

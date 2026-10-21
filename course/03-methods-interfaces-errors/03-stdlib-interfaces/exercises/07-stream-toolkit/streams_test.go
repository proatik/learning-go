package streams

import (
	"bytes"
	"errors"
	"io"
	"strings"
	"testing"
)

func TestHead(t *testing.T) {
	tests := []struct {
		in   string
		n    int64
		want string
	}{
		{"0123456789", 4, "0123"},
		{"abc", 10, "abc"},
		{"abc", 3, "abc"},
		{"abc", 0, ""},
		{"abc", -5, ""},
		{"", 3, ""},
		{"héllo", 3, "hé"},
		{strings.Repeat("x", 10000), 8000, strings.Repeat("x", 8000)},
	}
	for _, tt := range tests {
		got, err := Head(strings.NewReader(tt.in), tt.n)
		if err != nil {
			t.Errorf("Head(%.20q, %d) error = %v, want nil", tt.in, tt.n, err)
		}
		if got != tt.want {
			t.Errorf("Head(%.20q, %d) = %.20q (len %d), want %.20q (len %d)", tt.in, tt.n, got, len(got), tt.want, len(tt.want))
		}
	}
}

func TestHeadLeavesRest(t *testing.T) {
	r := strings.NewReader("first|second")
	Head(r, 6)
	rest, _ := io.ReadAll(r)
	if got, want := string(rest), "second"; got != want {
		t.Errorf("after Head(r, 6) on %q, rest of r = %q, want %q", "first|second", got, want)
	}
}

// spyReader records whether it has been read from.
type spyReader struct {
	r    io.Reader
	used bool
}

func (s *spyReader) Read(p []byte) (int, error) {
	s.used = true
	return s.r.Read(p)
}

func TestJoin(t *testing.T) {
	tests := []struct {
		sep   string
		parts []string
		want  string
	}{
		{", ", []string{"a", "b", "c"}, "a, b, c"},
		{"-", nil, ""},
		{"-", []string{"only"}, "only"},
		{"", []string{"ab", "cd"}, "abcd"},
		{"\n", []string{"", "", ""}, "\n\n"},
		{" | ", []string{"日本", "語"}, "日本 | 語"},
	}
	for _, tt := range tests {
		var rs []io.Reader
		for _, p := range tt.parts {
			rs = append(rs, strings.NewReader(p))
		}
		r := Join(tt.sep, rs...)
		if r == nil {
			t.Fatalf("Join(%q, %q) returned a nil reader", tt.sep, tt.parts)
		}
		b, err := io.ReadAll(r)
		if err != nil {
			t.Errorf("ReadAll(Join(%q, %q)) error = %v, want nil", tt.sep, tt.parts, err)
		}
		if got := string(b); got != tt.want {
			t.Errorf("ReadAll(Join(%q, %q)) = %q, want %q", tt.sep, tt.parts, got, tt.want)
		}
	}
}

func TestJoinIsLazy(t *testing.T) {
	spy := &spyReader{r: strings.NewReader("x")}
	r := Join("+", strings.NewReader("a"), spy)
	if r == nil {
		t.Fatalf("Join returned a nil reader")
	}
	if spy.used {
		t.Errorf("Join(...) read from a part before the caller read: it must only build the reader")
	}
	b, _ := io.ReadAll(r)
	if got, want := string(b), "a+x"; got != want {
		t.Errorf("ReadAll(Join(%q, %q, %q)) = %q, want %q", "+", "a", "x", got, want)
	}
}

var errFull = errors.New("disk full")

type failWriter struct{}

func (failWriter) Write(p []byte) (int, error) { return 0, errFull }

func TestBroadcast(t *testing.T) {
	var b1, b2 bytes.Buffer
	var sb strings.Builder
	if err := Broadcast("hi\n", &b1, &b2, &sb); err != nil {
		t.Errorf("Broadcast(%q, 3 buffers) error = %v, want nil", "hi\n", err)
	}
	for i, got := range []string{b1.String(), b2.String(), sb.String()} {
		if got != "hi\n" {
			t.Errorf("Broadcast(%q, ...) writer %d holds %q, want %q", "hi\n", i, got, "hi\n")
		}
	}
	if err := Broadcast("nobody"); err != nil {
		t.Errorf("Broadcast(%q) with no writers error = %v, want nil", "nobody", err)
	}
	var b3 bytes.Buffer
	if err := Broadcast("x", &b3, failWriter{}); !errors.Is(err, errFull) {
		t.Errorf("Broadcast(%q, buffer, failing writer) error = %v, want %v", "x", err, errFull)
	}
	if b3.String() != "x" {
		t.Errorf("Broadcast(%q, buffer, failing writer) buffer holds %q, want %q", "x", b3.String(), "x")
	}
}

var errNet = errors.New("connection reset")

type brokenReader struct {
	data string
	done bool
}

func (b *brokenReader) Read(p []byte) (int, error) {
	if b.done {
		return 0, errNet
	}
	b.done = true
	return copy(p, b.data), nil
}

func TestReadAndLog(t *testing.T) {
	tests := []string{"data", "", "line 1\nline 2\n", strings.Repeat("abc", 5000)}
	for _, in := range tests {
		var log bytes.Buffer
		got, err := ReadAndLog(strings.NewReader(in), &log)
		if err != nil {
			t.Errorf("ReadAndLog(%.20q) error = %v, want nil", in, err)
		}
		if got != in {
			t.Errorf("ReadAndLog(%.20q) = %.20q, want %.20q", in, got, in)
		}
		if log.String() != in {
			t.Errorf("ReadAndLog(%.20q) log = %.20q, want %.20q", in, log.String(), in)
		}
	}

	var log bytes.Buffer
	got, err := ReadAndLog(&brokenReader{data: "part"}, &log)
	if !errors.Is(err, errNet) {
		t.Errorf("ReadAndLog(reader failing after %q) error = %v, want %v", "part", err, errNet)
	}
	if got != "part" || log.String() != "part" {
		t.Errorf("ReadAndLog(reader failing after %q) = %q, log %q, want %q, log %q", "part", got, log.String(), "part", "part")
	}
}

// callCounter counts Write calls.
type callCounter struct {
	bytes.Buffer
	calls int
}

func (c *callCounter) Write(p []byte) (int, error) {
	c.calls++
	return c.Buffer.Write(p)
}

func TestWriteLines(t *testing.T) {
	tests := []struct {
		lines []string
		want  string
	}{
		{[]string{"x", "y"}, "x\ny\n"},
		{nil, ""},
		{[]string{""}, "\n"},
		{[]string{"héllo", "wörld"}, "héllo\nwörld\n"},
	}
	for _, tt := range tests {
		var buf bytes.Buffer
		if err := WriteLines(&buf, tt.lines); err != nil {
			t.Errorf("WriteLines(%q) error = %v, want nil", tt.lines, err)
		}
		if got := buf.String(); got != tt.want {
			t.Errorf("WriteLines(%q) = %q, want %q", tt.lines, got, tt.want)
		}
	}
}

func TestWriteLinesIsBuffered(t *testing.T) {
	lines := make([]string, 300)
	for i := range lines {
		lines[i] = "row"
	}
	c := &callCounter{}
	if err := WriteLines(c, lines); err != nil {
		t.Errorf("WriteLines(300 lines) error = %v, want nil", err)
	}
	if got, want := c.String(), strings.Repeat("row\n", 300); got != want {
		t.Errorf("WriteLines(300 x %q) wrote %d bytes, want %d (did you Flush?)", "row", len(got), len(want))
	}
	if c.calls > 2 || c.calls == 0 {
		t.Errorf("WriteLines(300 lines) made %d Write calls, want 1 or 2 (use a bufio.Writer)", c.calls)
	}
}

func TestWriteLinesError(t *testing.T) {
	if err := WriteLines(failWriter{}, []string{"a", "b"}); !errors.Is(err, errFull) {
		t.Errorf("WriteLines(failing writer, %q) error = %v, want %v", []string{"a", "b"}, err, errFull)
	}
}

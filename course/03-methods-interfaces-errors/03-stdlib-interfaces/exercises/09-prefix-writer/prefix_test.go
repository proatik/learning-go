package prefix

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"strings"
	"testing"
)

var _ io.Writer = (*PrefixWriter)(nil)

// reference prefixes every line of s the simple way, for comparison.
func reference(s, prefix string) string {
	if s == "" {
		return ""
	}
	lines := strings.SplitAfter(s, "\n")
	var sb strings.Builder
	for _, l := range lines {
		if l != "" {
			sb.WriteString(prefix + l)
		}
	}
	return sb.String()
}

func writeAll(t *testing.T, w io.Writer, chunks []string) {
	t.Helper()
	for _, c := range chunks {
		n, err := w.Write([]byte(c))
		if n != len(c) || err != nil {
			t.Errorf("Write(%q) = %d, %v, want %d, <nil>", c, n, err, len(c))
		}
	}
}

func TestPrefixWriter(t *testing.T) {
	tests := []struct {
		prefix string
		writes []string
		want   string
	}{
		{"> ", []string{"hello\n"}, "> hello\n"},
		{"> ", []string{"hel", "lo\nwor", "ld\n"}, "> hello\n> world\n"},
		{"> ", []string{"a\nb"}, "> a\n> b"},
		{"> ", []string{"a\n"}, "> a\n"},
		{"> ", []string{"\n\n"}, "> \n> \n"},
		{"> ", []string{""}, ""},
		{"> ", nil, ""},
		{"> ", []string{"a", "", "b\n", "", "c"}, "> ab\n> c"},
		{"> ", []string{"x\n", "\n", "y"}, "> x\n> \n> y"},
		{"", []string{"a\nb\n"}, "a\nb\n"},
		{"日本: ", []string{"un\ndeux\n"}, "日本: un\n日本: deux\n"},
		{"# ", []string{"no newline", " still same line"}, "# no newline still same line"},
	}
	for _, tt := range tests {
		var buf bytes.Buffer
		pw := NewPrefixWriter(&buf, tt.prefix)
		if pw == nil {
			t.Fatalf("NewPrefixWriter(&buf, %q) = nil", tt.prefix)
		}
		writeAll(t, pw, tt.writes)
		if got := buf.String(); got != tt.want {
			t.Errorf("PrefixWriter(%q) after writes %q = %q, want %q", tt.prefix, tt.writes, got, tt.want)
		}
	}
}

func TestPrefixWriterAnySplit(t *testing.T) {
	inputs := []string{
		"first line\nsecond line\n",
		"a\n\nb\n\n\nc",
		"\nstarts with newline",
		"one",
	}
	for _, in := range inputs {
		want := reference(in, "| ")
		// Every way of cutting the input into two writes.
		for cut := 0; cut <= len(in); cut++ {
			var buf bytes.Buffer
			pw := NewPrefixWriter(&buf, "| ")
			if pw == nil {
				t.Fatalf("NewPrefixWriter returned nil")
			}
			writeAll(t, pw, []string{in[:cut], in[cut:]})
			if got := buf.String(); got != want {
				t.Errorf("PrefixWriter(%q) after writes %q = %q, want %q", "| ", []string{in[:cut], in[cut:]}, got, want)
			}
		}
		// One byte at a time.
		var buf bytes.Buffer
		pw := NewPrefixWriter(&buf, "| ")
		for i := range len(in) {
			pw.Write([]byte{in[i]})
		}
		if got := buf.String(); got != want {
			t.Errorf("PrefixWriter(%q) writing %q one byte at a time = %q, want %q", "| ", in, got, want)
		}
	}
}

func TestPrefixWriterFprintf(t *testing.T) {
	var sb strings.Builder
	pw := NewPrefixWriter(&sb, "  ")
	fmt.Fprintf(pw, "x=%d\ny=", 1)
	fmt.Fprintf(pw, "%d\n", 2)
	fmt.Fprintln(pw, "done")
	if got, want := sb.String(), "  x=1\n  y=2\n  done\n"; got != want {
		t.Errorf("Fprintf through PrefixWriter(%q) = %q, want %q", "  ", got, want)
	}
}

func TestPrefixWriterDoesNotModifyInput(t *testing.T) {
	var buf bytes.Buffer
	pw := NewPrefixWriter(&buf, ">>")
	p := []byte("ab\ncd\n")
	pw.Write(p)
	if string(p) != "ab\ncd\n" {
		t.Errorf("Write(%q) changed its argument to %q", "ab\ncd\n", p)
	}
}

var errClosed = errors.New("writer closed")

type failWriter struct{}

func (failWriter) Write(p []byte) (int, error) { return 0, errClosed }

func TestPrefixWriterError(t *testing.T) {
	pw := NewPrefixWriter(failWriter{}, "> ")
	n, err := pw.Write([]byte("hello\n"))
	if n != 0 || !errors.Is(err, errClosed) {
		t.Errorf("Write(%q) to a failing writer = %d, %v, want 0, %v", "hello\n", n, err, errClosed)
	}
}

func TestNewLogger(t *testing.T) {
	var console, file bytes.Buffer
	l := NewLogger("[app] ", &console, &file)
	if l == nil {
		t.Fatalf("NewLogger(%q, ...) = nil", "[app] ")
	}
	fmt.Fprintln(l, "started")
	fmt.Fprint(l, "pid=")
	fmt.Fprint(l, "42\nready\n")
	want := "[app] started\n[app] pid=42\n[app] ready\n"
	if got := console.String(); got != want {
		t.Errorf("NewLogger(%q) console = %q, want %q", "[app] ", got, want)
	}
	if got := file.String(); got != want {
		t.Errorf("NewLogger(%q) file = %q, want %q", "[app] ", got, want)
	}
}

func TestNewLoggerError(t *testing.T) {
	var buf bytes.Buffer
	l := NewLogger("! ", &buf, failWriter{})
	if l == nil {
		t.Fatalf("NewLogger returned nil")
	}
	if _, err := io.WriteString(l, "boom\n"); !errors.Is(err, errClosed) {
		t.Errorf("WriteString(logger with a failing writer, %q) error = %v, want %v", "boom\n", err, errClosed)
	}
}

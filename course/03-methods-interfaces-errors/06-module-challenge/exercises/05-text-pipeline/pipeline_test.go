package pipeline_test

import (
	"bytes"
	"errors"
	"io"
	"strings"
	"testing"
	"testing/iotest"

	pipeline "learngo/course/03-methods-interfaces-errors/06-module-challenge/exercises/05-text-pipeline"
)

var errBoom = errors.New("boom")

// readAll reads r to the end, failing cleanly if the stage returned nil.
func readAll(t *testing.T, name string, r io.Reader) (string, error) {
	t.Helper()
	if r == nil {
		t.Fatalf("%s(...) = nil, want an io.Reader", name)
	}
	b, err := io.ReadAll(r)
	return string(b), err
}

type stageCase struct {
	name  string
	stage pipeline.Stage
	in    string
	want  string
}

var stageCases = []stageCase{
	{"Upper", pipeline.Upper, "", ""},
	{"Upper", pipeline.Upper, "hello, World!", "HELLO, WORLD!"},
	{"Upper", pipeline.Upper, "go 1.27 rocks\n", "GO 1.27 ROCKS\n"},
	{"Upper", pipeline.Upper, "héllo wörld", "HéLLO WöRLD"}, // non-ASCII bytes untouched
	{"Rot13", pipeline.Rot13, "", ""},
	{"Rot13", pipeline.Rot13, "Hello, World!", "Uryyb, Jbeyq!"},
	{"Rot13", pipeline.Rot13, "abcxyzABCXYZ", "nopklmNOPKLM"},
	{"Rot13", pipeline.Rot13, "Grüße 42", "Teüßr 42"},
	{"StripDigits", pipeline.StripDigits, "", ""},
	{"StripDigits", pipeline.StripDigits, "a1b22c333", "abc"},
	{"StripDigits", pipeline.StripDigits, "0123456789", ""},
	{"StripDigits", pipeline.StripDigits, "no digits", "no digits"},
	{"StripDigits", pipeline.StripDigits, "٣ is not ASCII 3", "٣ is not ASCII "},
}

func TestStages(t *testing.T) {
	for _, tt := range stageCases {
		got, err := readAll(t, tt.name, tt.stage(strings.NewReader(tt.in)))
		if err != nil || got != tt.want {
			t.Errorf("%s(%q) = %q, %v, want %q, nil", tt.name, tt.in, got, err, tt.want)
		}
	}
}

// TestReaderContract uses iotest.TestReader, which reads with many buffer
// sizes and checks the io.Reader rules.
func TestReaderContract(t *testing.T) {
	for _, tt := range stageCases {
		r := tt.stage(strings.NewReader(tt.in))
		if r == nil {
			t.Fatalf("%s(...) = nil, want an io.Reader", tt.name)
		}
		if err := iotest.TestReader(r, []byte(tt.want)); err != nil {
			t.Errorf("%s(%q) breaks the io.Reader contract: %v", tt.name, tt.in, err)
		}
	}
}

func TestOneByteAtATime(t *testing.T) {
	for _, tt := range stageCases {
		got, err := readAll(t, tt.name, tt.stage(iotest.OneByteReader(strings.NewReader(tt.in))))
		if err != nil || got != tt.want {
			t.Errorf("%s(OneByteReader(%q)) = %q, %v, want %q, nil", tt.name, tt.in, got, err, tt.want)
		}
	}
}

func TestStripDigitsNeverReturnsZeroNil(t *testing.T) {
	in := "12345a678b90"
	r := pipeline.StripDigits(iotest.OneByteReader(strings.NewReader(in)))
	if r == nil {
		t.Fatalf("StripDigits(...) = nil, want an io.Reader")
	}
	var got []byte
	buf := make([]byte, 4)
	for range 100 {
		n, err := r.Read(buf)
		if n == 0 && err == nil {
			t.Fatalf("StripDigits(%q).Read = 0, nil, want n > 0 or an error (keep reading until you have data)", in)
		}
		got = append(got, buf[:n]...)
		if err == io.EOF {
			break
		}
		if err != nil {
			t.Fatalf("StripDigits(%q).Read error = %v, want nil or io.EOF", in, err)
		}
	}
	if string(got) != "ab" {
		t.Errorf("StripDigits(%q) = %q, want %q", in, got, "ab")
	}
}

func TestErrorsPropagate(t *testing.T) {
	for _, s := range []struct {
		name  string
		stage pipeline.Stage
	}{{"Upper", pipeline.Upper}, {"Rot13", pipeline.Rot13}, {"StripDigits", pipeline.StripDigits}} {
		_, err := readAll(t, s.name, s.stage(iotest.ErrReader(errBoom)))
		if !errors.Is(err, errBoom) {
			t.Errorf("%s(ErrReader(boom)) error = %v, want %v", s.name, err, errBoom)
		}
		// Data that arrives together with an error must still be delivered.
		src := io.MultiReader(strings.NewReader("ab1"), iotest.ErrReader(errBoom))
		got, err := readAll(t, s.name, s.stage(src))
		if !errors.Is(err, errBoom) || len(got) < 2 {
			t.Errorf("%s(\"ab1\" then boom) = %q, %v, want the data then %v", s.name, got, err, errBoom)
		}
	}
}

func TestCountingWriter(t *testing.T) {
	var buf bytes.Buffer
	cw := &pipeline.CountingWriter{W: &buf}
	var w io.Writer = cw // *CountingWriter is an io.Writer
	for _, s := range []string{"one\n", "two\nthree", "", "\n", "é\n"} {
		n, err := io.WriteString(w, s)
		if n != len(s) || err != nil {
			t.Fatalf("Write(%q) = %d, %v, want %d, nil", s, n, err, len(s))
		}
	}
	if buf.String() != "one\ntwo\nthree\né\n" {
		t.Errorf("forwarded data = %q, want %q", buf.String(), "one\ntwo\nthree\né\n")
	}
	if cw.Bytes != 17 || cw.Lines != 4 {
		t.Errorf("Bytes, Lines = %d, %d, want 17, 4", cw.Bytes, cw.Lines)
	}
}

func TestCountingWriterNilW(t *testing.T) {
	cw := &pipeline.CountingWriter{}
	n, err := cw.Write([]byte("a\nb\n"))
	if n != 4 || err != nil {
		t.Fatalf("Write(%q) with nil W = %d, %v, want 4, nil", "a\nb\n", n, err)
	}
	if cw.Bytes != 4 || cw.Lines != 2 {
		t.Errorf("Bytes, Lines = %d, %d, want 4, 2", cw.Bytes, cw.Lines)
	}
}

// shortWriter accepts at most limit bytes per call, then fails.
type shortWriter struct{ limit int }

func (s shortWriter) Write(p []byte) (int, error) {
	if len(p) <= s.limit {
		return len(p), nil
	}
	return s.limit, errBoom
}

func TestCountingWriterShortWrite(t *testing.T) {
	cw := &pipeline.CountingWriter{W: shortWriter{limit: 3}}
	n, err := cw.Write([]byte("a\nbc\nd"))
	if n != 3 || !errors.Is(err, errBoom) {
		t.Fatalf("Write(%q) = %d, %v, want 3, %v", "a\nbc\nd", n, err, errBoom)
	}
	if cw.Bytes != 3 || cw.Lines != 1 {
		t.Errorf("Bytes, Lines = %d, %d, want 3, 1 (count only what was written)", cw.Bytes, cw.Lines)
	}
}

func TestRun(t *testing.T) {
	limit5 := func(r io.Reader) io.Reader { return io.LimitReader(r, 5) }
	tests := []struct {
		in     string
		stages []pipeline.Stage
		want   string
	}{
		{"", nil, ""},
		{"plain text", nil, "plain text"},
		{"hello 2026", []pipeline.Stage{pipeline.Upper}, "HELLO 2026"},
		{"hello 2026", []pipeline.Stage{pipeline.StripDigits, pipeline.Upper}, "HELLO "},
		{"Secret 7", []pipeline.Stage{pipeline.Rot13, pipeline.Rot13}, "Secret 7"},
		{"Secret 7", []pipeline.Stage{pipeline.Upper, pipeline.Rot13}, "FRPERG 7"},
		{"a1b2c3d4e5", []pipeline.Stage{limit5, pipeline.StripDigits}, "abc"},   // order matters:
		{"a1b2c3d4e5", []pipeline.Stage{pipeline.StripDigits, limit5}, "abcde"}, // stages run left to right
		{"x9y", []pipeline.Stage{nil, pipeline.StripDigits, nil}, "xy"},
	}
	for _, tt := range tests {
		var out strings.Builder
		n, err := pipeline.Run(&out, strings.NewReader(tt.in), tt.stages...)
		if err != nil || out.String() != tt.want || n != int64(len(tt.want)) {
			t.Errorf("Run(%q, %d stages) = %q, %d, %v, want %q, %d, nil", tt.in, len(tt.stages), out.String(), n, err, tt.want, len(tt.want))
		}
	}
}

func TestRunWithCountingWriter(t *testing.T) {
	var a, b bytes.Buffer
	cw := &pipeline.CountingWriter{W: io.MultiWriter(&a, &b)}
	in := "line 1\nline 22\nline 333\n"
	n, err := pipeline.Run(cw, strings.NewReader(in), pipeline.StripDigits, pipeline.Upper)
	want := "LINE \nLINE \nLINE \n"
	if err != nil || n != int64(len(want)) {
		t.Fatalf("Run(%q) = %d, %v, want %d, nil", in, n, err, len(want))
	}
	if a.String() != want || b.String() != want {
		t.Errorf("Run(%q) wrote %q and %q, want %q twice", in, a.String(), b.String(), want)
	}
	if cw.Bytes != int64(len(want)) || cw.Lines != 3 {
		t.Errorf("Bytes, Lines = %d, %d, want %d, 3", cw.Bytes, cw.Lines, len(want))
	}
}

func TestRunErrors(t *testing.T) {
	var out strings.Builder
	src := io.MultiReader(strings.NewReader("ok "), iotest.ErrReader(errBoom))
	n, err := pipeline.Run(&out, src, pipeline.Upper)
	if !errors.Is(err, errBoom) || out.String() != "OK " || n != 3 {
		t.Errorf("Run(\"ok \" then boom) = %q, %d, %v, want %q, 3, %v", out.String(), n, err, "OK ", errBoom)
	}

	n, err = pipeline.Run(shortWriter{limit: 2}, strings.NewReader("abcdef"), pipeline.Rot13)
	if !errors.Is(err, errBoom) || n != 2 {
		t.Errorf("Run to a failing writer = %d, %v, want 2, %v", n, err, errBoom)
	}
}

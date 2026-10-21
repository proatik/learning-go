package transform

import (
	"errors"
	"io"
	"strings"
	"testing"
)

var (
	_ io.Reader = UpperReader{}
	_ io.Reader = Rot13Reader{}
)

func readAll(t *testing.T, r io.Reader) string {
	t.Helper()
	b, err := io.ReadAll(r)
	if err != nil {
		t.Errorf("io.ReadAll error = %v, want nil", err)
	}
	return string(b)
}

func TestUpperReader(t *testing.T) {
	tests := []struct{ in, want string }{
		{"Hello, Go!", "HELLO, GO!"},
		{"", ""},
		{"already UPPER", "ALREADY UPPER"},
		{"héllo", "HéLLO"},
		{"abc xyz 0123 `{@[", "ABC XYZ 0123 `{@["},
		{strings.Repeat("go", 3000), strings.Repeat("GO", 3000)},
	}
	for _, tt := range tests {
		if got := readAll(t, UpperReader{R: strings.NewReader(tt.in)}); got != tt.want {
			t.Errorf("ReadAll(UpperReader(%.20q)) = %.20q, want %.20q", tt.in, got, tt.want)
		}
	}
}

func TestRot13Reader(t *testing.T) {
	tests := []struct{ in, want string }{
		{"Lbh penpxrq gur pbqr!", "You cracked the code!"},
		{"abc XYZ 123", "nop KLM 123"},
		{"", ""},
		{"abcdefghijklmnopqrstuvwxyz", "nopqrstuvwxyzabcdefghijklm"},
		{"ABCDEFGHIJKLMNOPQRSTUVWXYZ", "NOPQRSTUVWXYZABCDEFGHIJKLM"},
		{"mMnN", "zZaA"},
		{"café ☕", "pnsé ☕"},
		{"@[`{", "@[`{"},
	}
	for _, tt := range tests {
		if got := readAll(t, Rot13Reader{R: strings.NewReader(tt.in)}); got != tt.want {
			t.Errorf("ReadAll(Rot13Reader(%q)) = %q, want %q", tt.in, got, tt.want)
		}
	}
}

func TestRot13Twice(t *testing.T) {
	in := "The Quick Brown Fox, 42!"
	if got := readAll(t, Rot13Reader{Rot13Reader{strings.NewReader(in)}}); got != in {
		t.Errorf("ReadAll(Rot13Reader(Rot13Reader(%q))) = %q, want %q", in, got, in)
	}
	if got, want := readAll(t, UpperReader{Rot13Reader{strings.NewReader("uryyb")}}), "HELLO"; got != want {
		t.Errorf("ReadAll(UpperReader(Rot13Reader(%q))) = %q, want %q", "uryyb", got, want)
	}
}

// smallReads calls Read with a 3-byte buffer and collects the results.
func smallReads(t *testing.T, r io.Reader) string {
	t.Helper()
	var sb strings.Builder
	buf := make([]byte, 3)
	for range 1000 {
		n, err := r.Read(buf)
		sb.Write(buf[:n])
		if err == io.EOF {
			return sb.String()
		}
		if err != nil {
			t.Fatalf("Read error = %v", err)
		}
	}
	t.Fatalf("Read never returned io.EOF")
	return ""
}

func TestSmallBuffer(t *testing.T) {
	if got, want := smallReads(t, UpperReader{strings.NewReader("streaming data")}), "STREAMING DATA"; got != want {
		t.Errorf("UpperReader(%q) read 3 bytes at a time = %q, want %q", "streaming data", got, want)
	}
	if got, want := smallReads(t, Rot13Reader{strings.NewReader("streaming data")}), "fgernzvat qngn"; got != want {
		t.Errorf("Rot13Reader(%q) read 3 bytes at a time = %q, want %q", "streaming data", got, want)
	}
}

var errNet = errors.New("connection reset")

// lastChunk returns all its data together with err in a single call.
type lastChunk struct {
	data string
	err  error
	done bool
}

func (l *lastChunk) Read(p []byte) (int, error) {
	if l.done {
		return 0, l.err
	}
	l.done = true
	return copy(p, l.data), l.err
}

func TestPassesThroughDataWithError(t *testing.T) {
	for _, wantErr := range []error{io.EOF, errNet} {
		p := make([]byte, 16)
		n, err := UpperReader{&lastChunk{data: "abc", err: wantErr}}.Read(p)
		if n != 3 || string(p[:n]) != "ABC" || err != wantErr {
			t.Errorf("UpperReader.Read on (\"abc\", %v) = %d, %q, %v, want 3, %q, %v", wantErr, n, p[:n], err, "ABC", wantErr)
		}
		n, err = Rot13Reader{&lastChunk{data: "abc", err: wantErr}}.Read(p)
		if n != 3 || string(p[:n]) != "nop" || err != wantErr {
			t.Errorf("Rot13Reader.Read on (\"abc\", %v) = %d, %q, %v, want 3, %q, %v", wantErr, n, p[:n], err, "nop", wantErr)
		}
	}
}

func TestOnlyTouchesValidBytes(t *testing.T) {
	p := []byte("zzzzzzzz")
	n, _ := UpperReader{strings.NewReader("ab")}.Read(p)
	if got, want := string(p), "ABzzzzzz"; n != 2 || got != want {
		t.Errorf("UpperReader(%q).Read(%q) = %d, buffer %q, want 2, buffer %q (only p[:n] may change)", "ab", "zzzzzzzz", n, got, want)
	}
	p = []byte("aaaaaaaa")
	n, _ = Rot13Reader{strings.NewReader("ab")}.Read(p)
	if got, want := string(p), "noaaaaaa"; n != 2 || got != want {
		t.Errorf("Rot13Reader(%q).Read(%q) = %d, buffer %q, want 2, buffer %q (only p[:n] may change)", "ab", "aaaaaaaa", n, got, want)
	}
}

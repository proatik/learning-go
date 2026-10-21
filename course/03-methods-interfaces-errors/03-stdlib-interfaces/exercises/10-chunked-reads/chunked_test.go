package chunked

import (
	"errors"
	"fmt"
	"io"
	"strings"
	"testing"
)

var _ io.Reader = (*ChunkReader)(nil)

// result is one Read outcome: the data returned and the error.
type result struct {
	data string
	err  error
}

func (r result) String() string { return fmt.Sprintf("(%q, %v)", r.data, r.err) }

// reads calls Read with a buffer of bufSize until io.EOF (at most 50 times).
func reads(c *ChunkReader, bufSize int) []result {
	var out []result
	buf := make([]byte, bufSize)
	for range 50 {
		n, err := c.Read(buf)
		if n < 0 || n > len(buf) {
			out = append(out, result{fmt.Sprintf("<invalid n=%d>", n), err})
			return out
		}
		out = append(out, result{string(buf[:n]), err})
		if err != nil {
			return out
		}
	}
	return out
}

func TestChunkReader(t *testing.T) {
	EOF := io.EOF
	tests := []struct {
		data    string
		size    int
		withEOF bool
		bufSize int
		want    []result
	}{
		{"abcdefg", 3, false, 100, []result{{"abc", nil}, {"def", nil}, {"g", nil}, {"", EOF}}},
		{"abcdefg", 3, true, 100, []result{{"abc", nil}, {"def", nil}, {"g", EOF}}},
		{"abcdef", 3, true, 100, []result{{"abc", nil}, {"def", EOF}}},
		{"abcdef", 3, false, 100, []result{{"abc", nil}, {"def", nil}, {"", EOF}}},
		{"abcde", 10, false, 2, []result{{"ab", nil}, {"cd", nil}, {"e", nil}, {"", EOF}}},
		{"abc", 0, false, 100, []result{{"a", nil}, {"b", nil}, {"c", nil}, {"", EOF}}},
		{"ab", -4, true, 100, []result{{"a", nil}, {"b", EOF}}},
		{"", 3, false, 100, []result{{"", EOF}}},
		{"", 3, true, 100, []result{{"", EOF}}},
		{"héllo", 2, false, 100, []result{{"h\xc3", nil}, {"\xa9l", nil}, {"lo", nil}, {"", EOF}}},
		{"xyz", 100, true, 100, []result{{"xyz", EOF}}},
	}
	for _, tt := range tests {
		c := &ChunkReader{Data: tt.data, Size: tt.size, EOFWithData: tt.withEOF}
		got := reads(c, tt.bufSize)
		if fmt.Sprint(got) != fmt.Sprint(tt.want) {
			t.Errorf("ChunkReader{%q, Size: %d, EOFWithData: %v} read with %d-byte buffer = %v, want %v",
				tt.data, tt.size, tt.withEOF, tt.bufSize, got, tt.want)
		}
	}
}

func TestChunkReaderAfterEOF(t *testing.T) {
	c := &ChunkReader{Data: "ab", Size: 5, EOFWithData: true}
	buf := make([]byte, 10)
	c.Read(buf)
	for i := range 3 {
		if n, err := c.Read(buf); n != 0 || err != io.EOF {
			t.Errorf("Read #%d after the data is used up = %d, %v, want 0, EOF", i+2, n, err)
		}
	}
}

func TestChunkReaderEmptyBuffer(t *testing.T) {
	c := &ChunkReader{Data: "abc", Size: 2}
	if n, err := c.Read(nil); n != 0 || err != nil {
		t.Errorf("Read(nil) with data left = %d, %v, want 0, <nil>", n, err)
	}
	got := reads(c, 10)
	want := []result{{"ab", nil}, {"c", nil}, {"", io.EOF}}
	if fmt.Sprint(got) != fmt.Sprint(want) {
		t.Errorf("reads after Read(nil) = %v, want %v (an empty read must not consume data)", got, want)
	}
}

func TestReadAllChunkedWithChunkReader(t *testing.T) {
	inputs := []string{"hello world", "", "x", strings.Repeat("0123456789", 300), "日本語テキスト"}
	for _, in := range inputs {
		for _, size := range []int{1, 2, 7, 512, 513, 5000} {
			for _, withEOF := range []bool{false, true} {
				got, err := ReadAllChunked(&ChunkReader{Data: in, Size: size, EOFWithData: withEOF})
				if err != nil {
					t.Errorf("ReadAllChunked(ChunkReader{%.15q, Size: %d, EOFWithData: %v}) error = %v, want nil", in, size, withEOF, err)
				}
				if string(got) != in {
					t.Errorf("ReadAllChunked(ChunkReader{%.15q, Size: %d, EOFWithData: %v}) = %.15q (len %d), want %.15q (len %d)",
						in, size, withEOF, got, len(got), in, len(in))
				}
			}
		}
	}
}

func TestReadAllChunkedStringsReader(t *testing.T) {
	in := strings.Repeat("abc\n", 1000)
	got, err := ReadAllChunked(strings.NewReader(in))
	if err != nil || string(got) != in {
		t.Errorf("ReadAllChunked(strings.Reader of %d bytes) = %d bytes, %v, want %d bytes, <nil>", len(in), len(got), err, len(in))
	}
}

// script replays a fixed list of Read results, then returns (0, EOF).
// It records a problem if it is read again after returning io.EOF.
type script struct {
	steps    []result
	i        int
	afterEOF bool
	sawEOF   bool
}

func (s *script) Read(p []byte) (int, error) {
	if s.sawEOF {
		s.afterEOF = true
	}
	if s.i >= len(s.steps) {
		s.sawEOF = true
		return 0, io.EOF
	}
	st := s.steps[s.i]
	s.i++
	if st.err == io.EOF {
		s.sawEOF = true
	}
	return copy(p, st.data), st.err
}

var errNet = errors.New("connection reset")

func TestReadAllChunkedScripts(t *testing.T) {
	EOF := io.EOF
	tests := []struct {
		name    string
		steps   []result
		want    string
		wantErr error
	}{
		{"data then EOF", []result{{"ab", nil}, {"cd", nil}, {"", EOF}}, "abcd", nil},
		{"data with EOF", []result{{"ab", nil}, {"cd", EOF}}, "abcd", nil},
		{"only EOF", []result{{"", EOF}}, "", nil},
		{"data with error", []result{{"ab", nil}, {"cd", errNet}}, "abcd", errNet},
		{"error then nothing", []result{{"ab", nil}, {"", errNet}}, "ab", errNet},
		{"immediate error", []result{{"", errNet}}, "", errNet},
		{"some (0, nil) reads", []result{{"a", nil}, {"", nil}, {"", nil}, {"b", nil}, {"", nil}, {"c", EOF}}, "abc", nil},
	}
	for _, tt := range tests {
		s := &script{steps: tt.steps}
		got, err := ReadAllChunked(s)
		if !errors.Is(err, tt.wantErr) || (tt.wantErr == nil && err != nil) {
			t.Errorf("ReadAllChunked(%s %v) error = %v, want %v", tt.name, tt.steps, err, tt.wantErr)
		}
		if string(got) != tt.want {
			t.Errorf("ReadAllChunked(%s %v) = %q, want %q", tt.name, tt.steps, got, tt.want)
		}
		if s.afterEOF {
			t.Errorf("ReadAllChunked(%s %v) called Read again after io.EOF", tt.name, tt.steps)
		}
	}
}

func TestReadAllChunked99EmptyReads(t *testing.T) {
	steps := []result{{"start-", nil}}
	for range 99 {
		steps = append(steps, result{"", nil})
	}
	steps = append(steps, result{"end", nil})
	for range 99 {
		steps = append(steps, result{"", nil})
	}
	steps = append(steps, result{"!", io.EOF})
	got, err := ReadAllChunked(&script{steps: steps})
	if err != nil || string(got) != "start-end!" {
		t.Errorf("ReadAllChunked(data, 99 x (0, nil), data, 99 x (0, nil), data+EOF) = %q, %v, want %q, <nil>", got, err, "start-end!")
	}
}

// stall returns some data once, then (0, nil) forever.
type stall struct {
	calls int
}

func (s *stall) Read(p []byte) (int, error) {
	s.calls++
	if s.calls == 1 {
		return copy(p, "partial"), nil
	}
	return 0, nil
}

func TestReadAllChunkedNoProgress(t *testing.T) {
	s := &stall{}
	got, err := ReadAllChunked(s)
	if !errors.Is(err, io.ErrNoProgress) {
		t.Errorf("ReadAllChunked(reader stuck on (0, nil)) error = %v, want %v", err, io.ErrNoProgress)
	}
	if string(got) != "partial" {
		t.Errorf("ReadAllChunked(reader stuck on (0, nil)) = %q, want %q", got, "partial")
	}
	if s.calls != 101 {
		t.Errorf("ReadAllChunked(reader stuck on (0, nil)) made %d Read calls, want 101 (1 with data + 100 empty)", s.calls)
	}
}

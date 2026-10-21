package pipeline

import (
	"bytes"
	"io"
)

// Stage transforms one stream into another.
type Stage func(io.Reader) io.Reader

// mapReader applies f to every byte read from r.
type mapReader struct {
	r io.Reader
	f func(byte) byte
}

// Read reads from the source and transforms the bytes in place.
func (m *mapReader) Read(p []byte) (int, error) {
	n, err := m.r.Read(p)
	for i := range p[:n] {
		p[i] = m.f(p[i])
	}
	return n, err
}

// Upper returns a reader that turns ASCII a–z into A–Z.
func Upper(r io.Reader) io.Reader {
	return &mapReader{r: r, f: func(b byte) byte {
		if 'a' <= b && b <= 'z' {
			return b - 'a' + 'A'
		}
		return b
	}}
}

// Rot13 returns a reader that rotates ASCII letters by 13 places.
func Rot13(r io.Reader) io.Reader {
	return &mapReader{r: r, f: func(b byte) byte {
		switch {
		case 'a' <= b && b <= 'z':
			return 'a' + (b-'a'+13)%26
		case 'A' <= b && b <= 'Z':
			return 'A' + (b-'A'+13)%26
		}
		return b
	}}
}

// stripReader drops ASCII digits from its source.
type stripReader struct {
	r io.Reader
}

// Read fills p with the non-digit bytes of the source. It never returns
// 0, nil for a non-empty p.
func (s *stripReader) Read(p []byte) (int, error) {
	if len(p) == 0 {
		return 0, nil
	}
	for {
		n, err := s.r.Read(p)
		k := 0
		for _, b := range p[:n] {
			if b < '0' || b > '9' {
				p[k] = b
				k++
			}
		}
		if k > 0 || err != nil {
			return k, err
		}
	}
}

// StripDigits returns a reader that removes the ASCII digits 0–9.
func StripDigits(r io.Reader) io.Reader {
	return &stripReader{r: r}
}

// CountingWriter forwards writes to W and counts bytes and newlines.
// A nil W discards the data but still counts it.
type CountingWriter struct {
	W     io.Writer
	Bytes int64
	Lines int
}

// Write forwards p to W and counts what was actually written.
func (c *CountingWriter) Write(p []byte) (int, error) {
	n, err := len(p), error(nil)
	if c.W != nil {
		n, err = c.W.Write(p)
	}
	c.Bytes += int64(n)
	c.Lines += bytes.Count(p[:n], []byte{'\n'})
	return n, err
}

// Run applies the stages to src in order (nil stages are skipped) and
// copies the result to dst. It returns the number of bytes copied.
func Run(dst io.Writer, src io.Reader, stages ...Stage) (int64, error) {
	r := src
	for _, st := range stages {
		if st != nil {
			r = st(r)
		}
	}
	return io.Copy(dst, r)
}

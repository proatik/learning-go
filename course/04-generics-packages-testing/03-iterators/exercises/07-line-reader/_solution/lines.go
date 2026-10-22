package lines

import (
	"bufio"
	"io"
	"iter"
	"strings"
)

// LineReader reads numbered, trimmed lines from an io.Reader,
// skipping blank lines and # comments.
type LineReader struct {
	sc   *bufio.Scanner
	line int // number of the last line read from sc
}

// NewLineReader returns a LineReader that reads from r.
func NewLineReader(r io.Reader) *LineReader {
	return &LineReader{sc: bufio.NewScanner(r)}
}

// Err returns the first read error, or nil if the input ended cleanly.
func (lr *LineReader) Err() error {
	return lr.sc.Err()
}

// All returns an iterator over (line number, trimmed line) pairs.
// Blank lines and lines starting with # are skipped but still counted.
// It returns a single-use iterator: ranging again continues where the
// previous loop stopped.
func (lr *LineReader) All() iter.Seq2[int, string] {
	return func(yield func(int, string) bool) {
		for lr.sc.Scan() {
			lr.line++
			text := strings.TrimSpace(lr.sc.Text())
			if text == "" || strings.HasPrefix(text, "#") {
				continue
			}
			if !yield(lr.line, text) {
				return
			}
		}
	}
}

package streams

import (
	"bufio"
	"io"
	"strings"
)

// Head returns the first n bytes of r (fewer if r is shorter).
func Head(r io.Reader, n int64) (string, error) {
	if n <= 0 {
		return "", nil
	}
	b, err := io.ReadAll(io.LimitReader(r, n))
	return string(b), err
}

// Join returns a reader that yields every part in order, with sep between them.
func Join(sep string, parts ...io.Reader) io.Reader {
	var rs []io.Reader
	for i, p := range parts {
		if i > 0 {
			rs = append(rs, strings.NewReader(sep))
		}
		rs = append(rs, p)
	}
	return io.MultiReader(rs...)
}

// Broadcast writes msg to every writer in ws and returns the first error.
func Broadcast(msg string, ws ...io.Writer) error {
	_, err := io.WriteString(io.MultiWriter(ws...), msg)
	return err
}

// ReadAndLog reads all of r, copying every byte to log, and returns what it read.
func ReadAndLog(r io.Reader, log io.Writer) (string, error) {
	b, err := io.ReadAll(io.TeeReader(r, log))
	return string(b), err
}

// WriteLines writes each line followed by "\n" to w through a bufio.Writer.
func WriteLines(w io.Writer, lines []string) error {
	bw := bufio.NewWriter(w)
	for _, line := range lines {
		bw.WriteString(line) // errors are sticky: Flush reports the first one
		bw.WriteByte('\n')
	}
	return bw.Flush()
}

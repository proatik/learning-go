package streams

import (
	"io"
	"strings"
)

// Head returns the first n bytes of r (fewer if r is shorter).
func Head(r io.Reader, n int64) (string, error) {
	// TODO
	return "", nil
}

// Join returns a reader that yields every part in order, with sep between them.
func Join(sep string, parts ...io.Reader) io.Reader {
	// TODO
	return strings.NewReader("")
}

// Broadcast writes msg to every writer in ws and returns the first error.
func Broadcast(msg string, ws ...io.Writer) error {
	// TODO
	return nil
}

// ReadAndLog reads all of r, copying every byte to log, and returns what it read.
func ReadAndLog(r io.Reader, log io.Writer) (string, error) {
	// TODO
	return "", nil
}

// WriteLines writes each line followed by "\n" to w through a bufio.Writer.
func WriteLines(w io.Writer, lines []string) error {
	// TODO
	return nil
}

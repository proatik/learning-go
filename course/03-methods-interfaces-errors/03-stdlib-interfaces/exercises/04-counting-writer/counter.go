package counter

import "io"

// CountingWriter wraps W and counts the bytes successfully written to it.
// A nil W discards the bytes but still counts them.
type CountingWriter struct {
	W io.Writer
	N int64
}

// Write forwards p to c.W and adds the number of bytes written to c.N.
func (c *CountingWriter) Write(p []byte) (int, error) {
	// TODO
	return 0, nil
}

// Measure calls write with a CountingWriter wrapping w and returns the
// number of bytes written along with write's error.
func Measure(w io.Writer, write func(io.Writer) error) (int64, error) {
	// TODO
	return 0, nil
}

package pipeline

import "io"

// Stage transforms one stream into another.
type Stage func(io.Reader) io.Reader

// Upper returns a reader that turns ASCII a–z into A–Z.
func Upper(r io.Reader) io.Reader {
	// TODO: return a value of your own reader type
	return nil
}

// Rot13 returns a reader that rotates ASCII letters by 13 places.
func Rot13(r io.Reader) io.Reader {
	// TODO
	return nil
}

// StripDigits returns a reader that removes the ASCII digits 0–9.
func StripDigits(r io.Reader) io.Reader {
	// TODO
	return nil
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
	// TODO
	return 0, nil
}

// Run applies the stages to src in order (nil stages are skipped) and
// copies the result to dst. It returns the number of bytes copied.
func Run(dst io.Writer, src io.Reader, stages ...Stage) (int64, error) {
	// TODO
	return 0, nil
}

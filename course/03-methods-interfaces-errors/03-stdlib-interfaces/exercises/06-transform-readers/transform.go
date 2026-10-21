package transform

import "io"

// UpperReader upper-cases the ASCII letters read from R.
type UpperReader struct {
	R io.Reader
}

// Read reads from u.R into p and upper-cases the ASCII letters in p[:n].
func (u UpperReader) Read(p []byte) (int, error) {
	// TODO
	return 0, io.EOF
}

// Rot13Reader applies the ROT13 cipher to the ASCII letters read from R.
type Rot13Reader struct {
	R io.Reader
}

// Read reads from r.R into p and applies ROT13 to the ASCII letters in p[:n].
func (r Rot13Reader) Read(p []byte) (int, error) {
	// TODO
	return 0, io.EOF
}

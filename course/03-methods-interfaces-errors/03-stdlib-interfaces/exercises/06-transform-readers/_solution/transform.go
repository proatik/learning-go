package transform

import "io"

// UpperReader upper-cases the ASCII letters read from R.
type UpperReader struct {
	R io.Reader
}

// Read reads from u.R into p and upper-cases the ASCII letters in p[:n].
func (u UpperReader) Read(p []byte) (int, error) {
	n, err := u.R.Read(p)
	for i, c := range p[:n] {
		if 'a' <= c && c <= 'z' {
			p[i] = c - ('a' - 'A')
		}
	}
	return n, err
}

// Rot13Reader applies the ROT13 cipher to the ASCII letters read from R.
type Rot13Reader struct {
	R io.Reader
}

func rot13(c byte) byte {
	switch {
	case 'a' <= c && c <= 'z':
		return 'a' + (c-'a'+13)%26
	case 'A' <= c && c <= 'Z':
		return 'A' + (c-'A'+13)%26
	}
	return c
}

// Read reads from r.R into p and applies ROT13 to the ASCII letters in p[:n].
func (r Rot13Reader) Read(p []byte) (int, error) {
	n, err := r.R.Read(p)
	for i := range p[:n] {
		p[i] = rot13(p[i])
	}
	return n, err
}

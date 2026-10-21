package prefix

import "io"

// PrefixWriter writes a prefix at the start of every line written through it.
type PrefixWriter struct {
	w       io.Writer
	prefix  []byte
	midLine bool // false when the next byte starts a new line
}

// NewPrefixWriter returns a PrefixWriter that writes to w.
func NewPrefixWriter(w io.Writer, prefix string) *PrefixWriter {
	return &PrefixWriter{w: w, prefix: []byte(prefix)}
}

// Write writes p to the wrapped writer, inserting the prefix before the
// first byte of every line. It returns len(p) on success.
func (pw *PrefixWriter) Write(p []byte) (int, error) {
	if len(p) == 0 {
		return 0, nil
	}
	out := make([]byte, 0, len(p)+len(pw.prefix))
	midLine := pw.midLine
	for _, c := range p {
		if !midLine {
			out = append(out, pw.prefix...)
			midLine = true
		}
		out = append(out, c)
		if c == '\n' {
			midLine = false
		}
	}
	if _, err := pw.w.Write(out); err != nil {
		return 0, err
	}
	pw.midLine = midLine
	return len(p), nil
}

// NewLogger returns a writer that prefixes every line and copies the result
// to all of ws.
func NewLogger(prefix string, ws ...io.Writer) io.Writer {
	return NewPrefixWriter(io.MultiWriter(ws...), prefix)
}

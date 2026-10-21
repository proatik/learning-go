package prefix

import "io"

// PrefixWriter writes a prefix at the start of every line written through it.
type PrefixWriter struct {
	// TODO: add the fields you need (the wrapped writer, the prefix, and
	// whether the next byte starts a new line).
}

// NewPrefixWriter returns a PrefixWriter that writes to w.
func NewPrefixWriter(w io.Writer, prefix string) *PrefixWriter {
	// TODO
	return &PrefixWriter{}
}

// Write writes p to the wrapped writer, inserting the prefix before the
// first byte of every line. It returns len(p) on success.
func (pw *PrefixWriter) Write(p []byte) (int, error) {
	// TODO
	return 0, nil
}

// NewLogger returns a writer that prefixes every line and copies the result
// to all of ws.
func NewLogger(prefix string, ws ...io.Writer) io.Writer {
	// TODO
	return NewPrefixWriter(io.Discard, prefix)
}

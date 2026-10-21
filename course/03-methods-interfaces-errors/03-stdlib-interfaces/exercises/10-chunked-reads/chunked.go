package chunked

import "io"

// ChunkReader returns Data at most Size bytes per Read.
type ChunkReader struct {
	Data        string
	Size        int  // max bytes per Read; values < 1 mean 1
	EOFWithData bool // true: the last chunk is returned together with io.EOF
	// TODO: add a field to remember how much has been read
}

// Read copies the next chunk of Data into p.
func (c *ChunkReader) Read(p []byte) (int, error) {
	// TODO
	return 0, io.EOF
}

// ReadAllChunked reads r until io.EOF and returns everything it read.
func ReadAllChunked(r io.Reader) ([]byte, error) {
	// TODO
	return nil, nil
}

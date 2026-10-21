package chunked

import "io"

// ChunkReader returns Data at most Size bytes per Read.
type ChunkReader struct {
	Data        string
	Size        int  // max bytes per Read; values < 1 mean 1
	EOFWithData bool // true: the last chunk is returned together with io.EOF
	pos         int  // bytes of Data already returned
}

// Read copies the next chunk of Data into p.
func (c *ChunkReader) Read(p []byte) (int, error) {
	if c.pos >= len(c.Data) {
		return 0, io.EOF
	}
	size := max(c.Size, 1)
	n := copy(p, c.Data[c.pos:min(c.pos+size, len(c.Data))])
	c.pos += n
	if c.EOFWithData && c.pos == len(c.Data) {
		return n, io.EOF
	}
	return n, nil
}

// maxEmptyReads is how many (0, nil) results ReadAllChunked tolerates in a row.
const maxEmptyReads = 100

// ReadAllChunked reads r until io.EOF and returns everything it read.
func ReadAllChunked(r io.Reader) ([]byte, error) {
	var data []byte
	buf := make([]byte, 512)
	empty := 0
	for {
		n, err := r.Read(buf)
		data = append(data, buf[:n]...) // 1. keep the bytes first
		if err == io.EOF {              // 2. then look at the error
			return data, nil
		}
		if err != nil {
			return data, err
		}
		if n > 0 {
			empty = 0
			continue
		}
		empty++
		if empty >= maxEmptyReads {
			return data, io.ErrNoProgress
		}
	}
}

# Short Reads

Real readers (networks, pipes) hand out data in small, unpredictable pieces. Build a reader that does exactly that, and a `ReadAll` that survives every legal `Read` behaviour.

```go
// ChunkReader returns Data at most Size bytes per Read.
type ChunkReader struct {
	Data        string
	Size        int  // max bytes per Read; values < 1 mean 1
	EOFWithData bool // true: the last chunk comes back together with io.EOF
	// unexported fields of your choice
}

func (c *ChunkReader) Read(p []byte) (int, error)

// ReadAllChunked reads r until io.EOF and returns everything it read.
// (Write the loop yourself: no io.ReadAll, io.Copy, bufio, …)
func ReadAllChunked(r io.Reader) ([]byte, error)
```

## Examples

`&ChunkReader{Data: "abcdefg", Size: 3}` read with a big buffer:

| Call | Returns |
|---|---|
| 1st `Read` | `3, nil` → `"abc"` |
| 2nd `Read` | `3, nil` → `"def"` |
| 3rd `Read` | `1, nil` → `"g"` |
| 4th `Read` | `0, io.EOF` |

The same with `EOFWithData: true`: the 3rd `Read` returns `1, io.EOF` → `"g"`.

| Call | Result |
|---|---|
| `ReadAllChunked(&ChunkReader{Data: "hello world", Size: 2})` | `[]byte("hello world"), nil` |
| `ReadAllChunked(strings.NewReader(""))` | empty, `nil` |

## Rules
- `ChunkReader.Read` copies `min(len(p), Size, bytes left)` bytes. It never returns more than `len(p)`; a `Read` with an empty `p` returns `0, nil` (or `0, io.EOF` once the data is used up).
- Once all data has been returned, every further `Read` returns `0, io.EOF`.
- `ReadAllChunked` must handle all of these:
  - short reads (`n < len(buf)` with `err == nil`): keep reading;
  - data **together with** an error: keep the `n` bytes, then look at the error;
  - `io.EOF`: stop and return the data with a **`nil`** error;
  - any other error: return the data read so far **and** that error;
  - `(0, nil)` means "nothing happened": try again, but after **100** of these in a row return the data so far and `io.ErrNoProgress`.
- Use a fixed-size buffer (e.g. 512 bytes) for each `Read` and append `buf[:n]` to the result.

# Transforming Readers

Write two **reader decorators**: each wraps an `io.Reader` and changes the bytes as they flow through.

```go
type UpperReader struct {
	R io.Reader
}

func (u UpperReader) Read(p []byte) (int, error) // ASCII a-z → A-Z

type Rot13Reader struct {
	R io.Reader
}

func (r Rot13Reader) Read(p []byte) (int, error) // rotate ASCII letters by 13
```

## Examples

| Code | Result |
|---|---|
| `io.ReadAll(UpperReader{strings.NewReader("Hello, Go!")})` | `"HELLO, GO!"` |
| `io.ReadAll(UpperReader{strings.NewReader("héllo")})` | `"HéLLO"` |
| `io.ReadAll(Rot13Reader{strings.NewReader("Lbh penpxrq gur pbqr!")})` | `"You cracked the code!"` |
| `io.ReadAll(Rot13Reader{strings.NewReader("abc XYZ 123")})` | `"nop KLM 123"` |
| `io.ReadAll(Rot13Reader{Rot13Reader{strings.NewReader("Go")}})` | `"Go"` |

## Rules
- Call the wrapped reader **once** per `Read`, transform only `p[:n]`, and return the wrapped reader's `n` and `err` **unchanged** (data may arrive together with `io.EOF` or another error).
- Only ASCII letters change; every other byte (digits, punctuation, UTF-8 bytes ≥ 0x80) passes through untouched.
- ROT13: `a`→`n` … `m`→`z`, `n`→`a` … `z`→`m`, same for upper case (case is preserved).
- Don't read everything up front: the decorators must work on streams of any size, for any `len(p)`.

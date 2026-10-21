# IP Addresses and Colors

Make two small types print nicely by implementing `fmt.Stringer`.

```go
type IPAddr [4]byte

func (ip IPAddr) String() string // "127.0.0.1"

type Color struct {
	R, G, B uint8
}

func (c Color) String() string // "#ff8000"
```

## Examples

| Expression | Result |
|---|---|
| `fmt.Sprint(IPAddr{127, 0, 0, 1})` | `"127.0.0.1"` |
| `fmt.Sprint(IPAddr{})` | `"0.0.0.0"` |
| `fmt.Sprintf("%v", []IPAddr{{8, 8, 8, 8}, {1, 1, 1, 1}})` | `"[8.8.8.8 1.1.1.1]"` |
| `fmt.Sprint(Color{255, 128, 0})` | `"#ff8000"` |
| `fmt.Sprint(&Color{1, 2, 3})` | `"#010203"` |

## Rules
- `IPAddr`: the four bytes in decimal, separated by `.`, no padding.
- `Color`: `#` followed by each component as **two lowercase hex digits**.
- Use **value receivers**, so both `T` and `*T` values print through `String`.

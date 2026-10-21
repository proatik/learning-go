# Counter

A counter whose zero value is ready to use, with pointer-receiver methods that change it.

```go
type Counter struct {
	n int
}

func (c *Counter) Inc()
func (c *Counter) Dec()
func (c *Counter) Add(delta int)
func (c *Counter) Reset()
func (c *Counter) Value() int

// CountWords returns a counter per distinct word.
func CountWords(words []string) map[string]*Counter
```

## Examples

| Code | `Value()` afterwards |
|---|---|
| `var c Counter` | `0` |
| `c.Inc(); c.Inc()` | `2` |
| `c.Add(5); c.Dec()` | `4` |
| `var d Counter; d.Dec()` | `0` (never below zero) |
| `c.Add(-10)` | `0` |
| `c.Reset()` | `0` |
| `CountWords([]string{"go", "is", "go"})["go"].Value()` | `2` |

## Rules
- The counter **never goes below 0**: `Dec` on 0 does nothing; `Add` with a negative `delta` clamps at 0.
- All methods use pointer receivers (consistency rule), including `Value`.
- `CountWords` returns a non-nil map (empty for no words) with one `*Counter` per distinct word.

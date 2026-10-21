# Safe Index

Three operations that normally **panic** — out-of-range index, write to a nil map, failed type assertion — wrapped so they return an `error` instead.

```go
func SafeIndex(s []string, i int) (v string, err error)
func SafeSet(m map[string]int, key string, val int) (err error)
func StringAt(row []any, i int) (s string, err error)
```

| Call | Result |
|---|---|
| `SafeIndex([]string{"a", "b"}, 1)` | `"b", nil` |
| `SafeIndex([]string{"a", "b"}, 5)` | `""`, `safe index: runtime error: index out of range [5] with length 2` |
| `SafeIndex([]string{"a"}, -1)` | `""`, `safe index: runtime error: index out of range [-1]` |
| `SafeSet(map[string]int{}, "x", 1)` | `nil` (and `m["x"] == 1`) |
| `SafeSet(nil, "x", 1)` | `safe set: assignment to entry in nil map` |
| `StringAt([]any{"id", 7}, 0)` | `"id", nil` |
| `StringAt([]any{"id", 7}, 1)` | `""`, `string at 1: interface conversion: interface {} is int, not string` |
| `StringAt([]any{"id", 7}, 2)` | `""`, `string at 2: runtime error: index out of range [2] with length 2` |

## Rules
- Do **not** check bounds, `m == nil` or use the comma-ok assertion: write the plain operation (`s[i]`, `m[key] = val`, `row[i].(string)`) and **recover** in a deferred function.
- Run-time panic values implement `error`: wrap them with **`%w`**, so `errors.As(err, &re)` with `var re runtime.Error` succeeds.
- Error prefixes: `safe index: `, `safe set: `, `string at <i>: `.
- On error, the string result is `""`.
- (In real code, prefer the explicit check or `v, ok := x.(T)`. This exercise practises recovery.)

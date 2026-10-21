# Mixed Sum

Data arrived as `[]any` (think decoded JSON or a spreadsheet row). Add up the numbers and pull out the strings.

```go
// Valuer is any type that knows its numeric value.
type Valuer interface {
	Value() float64
}

func Sum(vals []any) (total float64, skipped int)
func Strings(vals []any) []string
```

## Examples

| Call | Result |
|---|---|
| `Sum([]any{1, int64(2), 2.5})` | `5.5, 0` |
| `Sum([]any{1, "4", nil, true})` | `1, 3` |
| `Sum([]any{float32(0.5), Cents(250)})` (`Cents` is a `Valuer` worth `2.5`) | `3, 0` |
| `n := 7; Sum([]any{&n, (*int)(nil)})` | `7, 1` |
| `Sum(nil)` | `0, 0` |
| `Strings([]any{"a", 1, "b", nil})` | `["a" "b"]` |
| `Strings([]any{1, 2})` | empty (`nil` is fine) |

## Rules
- `Sum` adds, converted to `float64`:
  - `int`, `int64`, `float32`, `float64` values
  - `*int`: the pointed-to value; a **nil** `*int` is skipped
  - any `Valuer`: the result of `Value()`
- Everything else is **skipped** and counted in `skipped`: strings (even `"4"`), bools, `nil`, other types.
- `Strings` returns the `string` elements in input order. Use the comma-ok assertion `s, ok := v.(string)`.

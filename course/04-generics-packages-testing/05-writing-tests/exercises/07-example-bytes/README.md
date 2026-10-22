# An Example test

Write `ExampleFormatBytes` with an `// Output:` comment that catches every broken version of `FormatBytes`.

```go
// FormatBytes formats a size in bytes for humans, with binary units.
func FormatBytes(n int64) string
```

`FormatBytes` is already written (in the read-only checker). You write the example in `bytesize_test.go`.

## Examples

| Call | Result |
|---|---|
| `FormatBytes(512)` | `"512 B"` |
| `FormatBytes(1023)` | `"1023 B"` |
| `FormatBytes(1024)` | `"1.0 KiB"` |
| `FormatBytes(1536)` | `"1.5 KiB"` |
| `FormatBytes(5 * 1024 * 1024)` | `"5.0 MiB"` |
| `FormatBytes(3 << 30)` | `"3.0 GiB"` |

## Rules

- Below `1024`: the whole number and `B`.
- Otherwise divide by `1024` until the value is below `1024` (or the unit is `TiB`), and print it with **one decimal**: `KiB`, `MiB`, `GiB`, `TiB`.
- `n` is never negative.
- Print one size per line with `fmt.Println`, then end the function with `// Output:` and the exact lines.
- An example **without** an output comment is compiled but never run. The checker requires at least one `Example…` with output.
- **The checker** runs your example on the correct `FormatBytes` (it must pass), then once per planted bug (it must fail every time).

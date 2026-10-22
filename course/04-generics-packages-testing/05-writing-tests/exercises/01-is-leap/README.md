# Leap years

Write a table-driven `TestIsLeap` that catches every broken version of `IsLeap`.

```go
// IsLeap reports whether year is a leap year in the Gregorian calendar.
func IsLeap(year int) bool
```

`IsLeap` is already written (in the read-only checker). You write the test in `leap_test.go`.

## Examples

| Call | Result | Why |
|---|---|---|
| `IsLeap(2023)` | `false` | not divisible by 4 |
| `IsLeap(2024)` | `true` | divisible by 4 |
| `IsLeap(1900)` | `false` | divisible by 100, but not by 400 |
| `IsLeap(2000)` | `true` | divisible by 400 |

## Rules

- A year is leap if it is divisible by 4, **except** years divisible by 100, **unless** they are also divisible by 400.
- Use a table: a slice of `struct{ year int; want bool }` and one loop.
- Report a mismatch as `t.Errorf("IsLeap(%d) = %t, want %t", tt.year, got, tt.want)`.
- **The checker** runs your test on the correct `IsLeap` (it must pass), then once per planted bug (it must fail every time).
- A bug that survives is reported as `Bug "name" = not caught`.

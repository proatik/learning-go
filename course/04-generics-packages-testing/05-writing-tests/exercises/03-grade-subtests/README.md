# Named subtests

Write `TestGrade` with **named subtests** (`t.Run`) that catch every broken version of `Grade`.

```go
// Grade turns a score from 0 to 100 into a letter grade.
func Grade(score int) string
```

`Grade` is already written (in the read-only checker). You write the test in `grade_test.go`.

## Examples

| Score | Grade |
|---|---|
| `90` to `100` | `"A"` |
| `80` to `89` | `"B"` |
| `70` to `79` | `"C"` |
| `60` to `69` | `"D"` |
| `0` to `59` | `"F"` |
| below `0` or above `100` | `"invalid"` |

## Rules

- Use a table with a `name` field, and run each case with `t.Run(tt.name, ...)`.
- Names should say what the case checks, e.g. `"lowest A"` or `"above 100"`.
- **At least 6 subtests.** The checker runs your test with `-v` and counts the `=== RUN   TestGrade/...` lines.
- Report a mismatch as `t.Errorf("Grade(%d) = %q, want %q", tt.score, got, tt.want)`.
- **The checker** runs your test on the correct `Grade` (it must pass), then once per planted bug (it must fail every time).

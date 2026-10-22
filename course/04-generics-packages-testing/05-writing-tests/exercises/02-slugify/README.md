# Slugs

Write a table-driven `TestSlugify` that catches every broken version of `Slugify`.

```go
// Slugify turns a title into a URL slug.
func Slugify(s string) string
```

`Slugify` is already written (in the read-only checker). You write the test in `slug_test.go`.

## Examples

| Call | Result |
|---|---|
| `Slugify("Hello, World!")` | `"hello-world"` |
| `Slugify("Go 1.27")` | `"go-1-27"` |
| `Slugify("  a -- b  ")` | `"a-b"` |
| `Slugify("café")` | `"caf"` |
| `Slugify("")` | `""` |

## Rules

- ASCII letters are kept and lowercased. ASCII digits are kept.
- Every other character is a separator, including spaces, punctuation and non-ASCII letters.
- A **run** of separators between two kept characters becomes **one** `-`.
- Separators at the start or the end produce nothing: a slug never starts or ends with `-`.
- Report a mismatch as `t.Errorf("Slugify(%q) = %q, want %q", tt.in, got, tt.want)`.
- **The checker** runs your test on the correct `Slugify` (it must pass), then once per planted bug (it must fail every time).

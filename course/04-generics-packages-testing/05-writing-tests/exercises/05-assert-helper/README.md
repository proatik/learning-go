# A test helper

Write a generic `assertEqual` helper and use it in `TestInitials` to catch every broken version of `Initials`.

```go
// Initials returns the upper-cased first letter of every word in name.
func Initials(name string) string

// Yours, in initials_test.go:
func assertEqual[T comparable](t *testing.T, got, want T)
```

`Initials` is already written (in the read-only checker). You write the helper and the test in `initials_test.go`.

## Examples

| Call | Result |
|---|---|
| `Initials("ada lovelace")` | `"AL"` |
| `Initials("Grace Brewster Hopper")` | `"GBH"` |
| `Initials("  ken   thompson ")` | `"KT"` |
| `Initials("émile zola")` | `"ÉZ"` |
| `Initials("")` | `""` |

## Rules

- Words are separated by any run of whitespace (`strings.Fields`). Each word adds its first **rune**, upper-cased.
- `assertEqual` calls `t.Errorf("got %v, want %v", got, want)` when `got != want`.
- **The failure must point at the caller.** The checker reads the `initials_test.go:LINE:` of every failure. A line inside `assertEqual` fails the exercise: call `t.Helper()` first.
- **The checker** runs your test on the correct `Initials` (it must pass), then once per planted bug (it must fail every time).

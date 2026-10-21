# Word Search Drills

Warm up with the searching and comparing helpers of package `slices`.

```go
func HasWord(words []string, w string) bool
func Position(words []string, w string) int
func FirstLonger(words []string, n int) int
func AnyEmpty(words []string) bool
func SameWords(a, b []string) bool
```

| Function | Returns |
|---|---|
| `HasWord` | whether `w` is in `words` |
| `Position` | index of the **first** `w`, or `-1` |
| `FirstLonger` | index of the first word with **more than** `n` bytes, or `-1` |
| `AnyEmpty` | whether some word is `""` |
| `SameWords` | whether `a` and `b` hold the same words in the same order |

| Call | Result |
|---|---|
| `HasWord([]string{"go", "zig"}, "zig")` | `true` |
| `HasWord(nil, "go")` | `false` |
| `Position([]string{"a", "b", "a"}, "a")` | `0` |
| `Position([]string{"a", "b"}, "c")` | `-1` |
| `FirstLonger([]string{"go", "rust", "haskell"}, 3)` | `1` |
| `FirstLonger([]string{"go"}, 5)` | `-1` |
| `AnyEmpty([]string{"a", "", "b"})` | `true` |
| `SameWords([]string{"a", "b"}, []string{"b", "a"})` | `false` |
| `SameWords(nil, []string{})` | `true` |

## Rules
- Each body is **one line**: a single call to `slices.Contains`, `Index`, `IndexFunc`, `ContainsFunc` or `Equal`.
- No `for` loops.
- Matching is exact and case-sensitive: `"Go" != "go"`.

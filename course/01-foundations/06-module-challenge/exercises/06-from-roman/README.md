# From Roman Numerals

The reverse of the previous exercise, with **strict validation**. Implement:

```go
func FromRoman(s string) (int, error)
```

| Input | Output |
|---|---|
| `"III"` | `3, nil` |
| `"IV"` | `4, nil` |
| `"XLII"` | `42, nil` |
| `"MCMXCIV"` | `1994, nil` |
| `"MMMCMXCIX"` | `3999, nil` |
| `""` | error |
| `"IIII"` | error (not canonical: `IV`) |
| `"IC"` | error (not canonical: `XCIX`) |
| `"iv"` | error (lowercase) |
| `"MMMM"` | error (4000 is out of range) |

## Rules
- Accept **only** the canonical form that `ToRoman` from the previous exercise produces for `1`–`3999`. Everything else — empty input, lowercase, spaces, other characters, non-canonical forms (`IIII`, `VV`, `IIX`, `XM`, `MCMC`, ...) — returns `0` and a non-nil error.
- Values: `I`=1, `V`=5, `X`=10, `L`=50, `C`=100, `D`=500, `M`=1000. A symbol smaller than the one after it is subtracted (`IX` = 9).
- One clean strategy: compute the value with the subtract-if-smaller rule (rejecting unknown characters), then convert it back to Roman and require the result to equal `s`. You'll need to write that converter here too (an unexported helper is fine).

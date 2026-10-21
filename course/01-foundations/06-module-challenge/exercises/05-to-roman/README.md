# To Roman Numerals

Implement:

```go
func ToRoman(n int) (string, error)
```

| Symbol | `I` | `V` | `X` | `L` | `C` | `D` | `M` |
|---|---|---|---|---|---|---|---|
| Value | 1 | 5 | 10 | 50 | 100 | 500 | 1000 |

Numbers are written largest value first. Six **subtractive** pairs replace four-in-a-row: `IV` (4), `IX` (9), `XL` (40), `XC` (90), `CD` (400), `CM` (900).

| Input | Output |
|---|---|
| `1` | `"I"` |
| `4` | `"IV"` |
| `9` | `"IX"` |
| `14` | `"XIV"` |
| `40` | `"XL"` |
| `1994` | `"MCMXCIV"` |
| `2027` | `"MMXXVII"` |
| `3999` | `"MMMCMXCIX"` |
| `0`, `-5`, `4000` | `"", error` |

## Rules
- Valid range is `1`–`3999` inclusive; anything else returns `""` and a non-nil error.
- Use the canonical (shortest standard) form: never `IIII`, `VV` or `IC`.
- Uppercase letters only. Build the result with a `strings.Builder`.
- **Slice preview:** `[]int{1000, 900, 500}` is a list of values you can index (`vals[i]`), measure (`len(vals)`) and `range` over — a lookup table of value/symbol pairs makes this exercise short.

# Convert to Any Base

Implement your own version of `strconv.FormatInt`:

```go
func ToBase(n int, base int) (string, error)
```

| Call | Result |
|---|---|
| `ToBase(10, 2)` | `"1010", nil` |
| `ToBase(255, 16)` | `"ff", nil` |
| `ToBase(35, 36)` | `"z", nil` |
| `ToBase(0, 7)` | `"0", nil` |
| `ToBase(-255, 16)` | `"-ff", nil` |
| `ToBase(math.MinInt, 2)` | `"-1000…000"` (a `1` and 63 zeros on 64-bit), `nil` |
| `ToBase(5, 1)` | `"", error` |
| `ToBase(5, 37)` | `"", error` |

## Rules
- `base` must be `2`–`36`; otherwise return `""` and a non-nil error.
- Digits are `0`–`9` then lowercase `a`–`z` (`a` = 10, `z` = 35).
- Zero is `"0"`. Negative numbers get a single leading `-`. No leading zeros, no `+`.
- It must work for **every** `int`, including `math.MinInt` — careful: `-math.MinInt` overflows back to `math.MinInt`! (Hint: the magnitude of any `int` fits in a `uint`.)
- Do **not** use `strconv`'s formatting functions or `fmt` to produce the digits — the tests compare your output with `strconv.FormatInt`.

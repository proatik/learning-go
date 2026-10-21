# Look and Say

Each term of the **look-and-say** sequence describes the previous one by reading its runs of equal digits aloud: `"1211"` is "one 1, one 2, two 1s" → `"111221"`.

```
1 → 11 → 21 → 1211 → 111221 → 312211 → 13112221 → ...
```

Implement:

```go
func Next(s string) (string, error)
func Term(n int) (string, error)
```

| Call | Result |
|---|---|
| `Next("1")` | `"11", nil` |
| `Next("1211")` | `"111221", nil` |
| `Next("3333333333")` | `"103", nil` (ten 3s) |
| `Next("")` | error |
| `Next("12a")` | error |
| `Term(1)` | `"1", nil` |
| `Term(5)` | `"111221", nil` |
| `Term(0)` | error |

## Rules
- `Next`: each run of the same digit becomes `<count><digit>`, where `<count>` is written in decimal (a run of 12 becomes `"12"` followed by the digit).
- `Next` returns `"", error` if `s` is empty or contains any byte that is not an ASCII digit `0`–`9`.
- `Term(1)` is `"1"`; `Term(n)` applies `Next` to `Term(n-1)`. Return `"", error` if `n < 1`.
- Build strings with a `strings.Builder`, not repeated `+=`.

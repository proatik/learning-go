# Collatz and GCD: while-style Loops

Neither function knows in advance how many iterations it needs — a job for the condition-only `for cond { }` loop.

Implement:

```go
func CollatzSteps(n int) int // steps until n reaches 1; -1 if n < 1
func GCD(a, b int) int       // greatest common divisor (Euclid's algorithm)
```

## `CollatzSteps`
Repeat until `n == 1`: if `n` is even, `n = n / 2`; otherwise `n = 3*n + 1`. Count the steps.
`6 → 3 → 10 → 5 → 16 → 8 → 4 → 2 → 1` is **8** steps. `CollatzSteps(1)` is `0`. For `n < 1` return `-1` (the loop would never end).

## `GCD`
Euclid: while `b != 0`, replace `(a, b)` with `(b, a % b)`. When `b` is `0`, the answer is `a`.

| Call | Result |
|---|---|
| `GCD(48, 18)` | `6` |
| `GCD(17, 5)` | `1` |
| `GCD(0, 9)` | `9` |
| `GCD(0, 0)` | `0` |
| `GCD(-12, 18)` | `6` |

## Rules
- Use condition-only `for` loops (no three-part `for`, no `range`).
- The result of `GCD` is never negative: make both inputs non-negative first.
- Go can swap in one statement: `a, b = b, a%b`.

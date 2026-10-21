# Nested Loops with Labels

Implement:

```go
func FirstPair(target, limit int) (int, int) // first a*b == target
func CountPrimes(n int) int                  // number of primes p with 2 <= p <= n
```

> A function can return **two** values: declare `(int, int)` as the result type and write `return a, b`. The caller receives them with `x, y := FirstPair(12, 10)`, just like `n, err := strconv.Atoi(s)`.

## `FirstPair`
Search `a` from `1` up to `limit`, and for each `a`, search `b` from `a` up to `limit`. Return the **first** `(a, b)` found with `a*b == target`. If there is none, return `(0, 0)`.

| Call | Result | Why |
|---|---|---|
| `FirstPair(12, 20)` | `1, 12` | `a = 1` already works |
| `FirstPair(12, 10)` | `2, 6` | `12` is out of range, so `a = 1` fails |
| `FirstPair(49, 10)` | `7, 7` | `b` may equal `a` |
| `FirstPair(13, 10)` | `0, 0` | 13 is prime and `> 10` |
| `FirstPair(0, 10)` | `0, 0` | `a` and `b` are at least 1 |

## `CountPrimes`
A prime is an integer `>= 2` whose only divisors are 1 and itself. Testing divisors `d` with `d*d <= p` is enough.

| Call | Result |
|---|---|
| `CountPrimes(10)` | `4` (2, 3, 5, 7) |
| `CountPrimes(2)` | `1` |
| `CountPrimes(1)` | `0` |
| `CountPrimes(100)` | `25` |

## Rules
- Both functions use **two nested loops** with a **label** on the outer loop.
- `FirstPair`: save the pair in variables declared before the loops, then leave both loops with `break <label>`; have a single `return` at the end.
- `CountPrimes`: when a divisor is found, skip straight to the next candidate with `continue <label>`.

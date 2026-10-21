# Recursion Drills

Implement each function **recursively** (no loops). Every one needs a base case and a recursive case that moves toward it.

```go
func Factorial(n int) (int, error)
func Power(base, exp int) int
func GCD(a, b int) int
func SumDigits(n int) int
```

| Function | Definition | Examples |
|---|---|---|
| `Factorial` | `0! = 1`, `n! = n × (n-1)!` | `Factorial(5)` → `120, nil`; `Factorial(-1)` → `0, <error>`; `Factorial(21)` → `0, <error>` |
| `Power` | `base^0 = 1`; see below for `exp > 0` | `Power(2, 10)` → `1024`; `Power(-3, 3)` → `-27`; `Power(7, 0)` → `1` |
| `GCD` | Euclid: `gcd(a, 0) = a`, `gcd(a, b) = gcd(b, a % b)` | `GCD(48, 18)` → `6`; `GCD(0, 5)` → `5`; `GCD(0, 0)` → `0` |
| `SumDigits` | sum of the decimal digits | `SumDigits(1234)` → `10`; `SumDigits(0)` → `0`; `SumDigits(-56)` → `11` |

## Rules
- `Factorial`: negative `n` is an error; so is `n > 20`, because `21!` overflows a 64-bit `int`.
- `Power`: `exp >= 0` (you may assume it). Use **exponentiation by squaring**, so `Power(1, 1_000_000)` takes ~20 calls instead of a million:
  - `exp` even: `base^exp = (base^(exp/2))²`
  - `exp` odd: `base^exp = base × base^(exp-1)`
- `GCD`: `a, b >= 0`.
- `SumDigits`: for negative `n`, sum the digits of `-n`. Hint: `n % 10` is the last digit, `n / 10` drops it.

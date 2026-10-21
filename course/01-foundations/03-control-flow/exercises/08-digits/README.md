# Working with Digits

Peel digits off an integer with `%` and `/` — no strings allowed.

```go
n % 10  // last digit:          1234 % 10 == 4
n / 10  // drop the last digit: 1234 / 10 == 123
```

Implement:

```go
func DigitSum(n int) int        // sum of the decimal digits of |n|
func Reverse(n int) int         // digits in reverse order, keeping the sign
func IsPalindrome(n int) bool   // reads the same both ways
```

## Examples
| Call | Result |
|---|---|
| `DigitSum(1234)` | `10` |
| `DigitSum(-123)` | `6` |
| `DigitSum(0)` | `0` |
| `Reverse(1234)` | `4321` |
| `Reverse(120)` | `21` |
| `Reverse(-56)` | `-65` |
| `IsPalindrome(12321)` | `true` |
| `IsPalindrome(10)` | `false` |
| `IsPalindrome(-121)` | `false` (the `-` sign only appears on one side) |

## Rules
- Use a `for n != 0 { ... }` (or `for n > 0`) loop; do **not** convert to a string.
- Handle negative input: work on the absolute value and fix the sign at the end where needed.
- `IsPalindrome` should call `Reverse`.

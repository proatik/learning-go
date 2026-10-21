# Leap Years

Implement:

```go
func IsLeap(year int) bool
func DaysInYear(year int) int // 366 for leap years, otherwise 365
```

Gregorian rule:
1. A year divisible by **400** is a leap year.
2. Otherwise, a year divisible by **100** is **not** a leap year.
3. Otherwise, a year divisible by **4** is a leap year.
4. Every other year is not.

## Examples
| Call | Result | Why |
|---|---|---|
| `IsLeap(2024)` | `true` | divisible by 4 |
| `IsLeap(2023)` | `false` | not divisible by 4 |
| `IsLeap(1900)` | `false` | divisible by 100 but not 400 |
| `IsLeap(2000)` | `true` | divisible by 400 |
| `DaysInYear(2024)` | `366` | |
| `DaysInYear(2100)` | `365` | |

## Rules
- Write `IsLeap` as an `if` / `else if` chain that follows the rule **in the order above** (the order matters!).
- `DaysInYear` must call `IsLeap`.
- `%` gives the remainder: `year%4 == 0` means "divisible by 4".

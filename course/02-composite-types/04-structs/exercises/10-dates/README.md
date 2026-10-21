# Calendar Dates

Implement Gregorian calendar arithmetic on a small value type, **without** the `time` package.

```go
type Date struct {
	Year, Month, Day int // Month 1–12, Day 1–31
}

func IsLeap(year int) bool
func DaysInMonth(year, month int) int
func Valid(d Date) bool
func Compare(a, b Date) int           // -1, 0 or +1
func IsBefore(a, b Date) bool
func AddDays(d Date, n int) Date      // n may be negative
func DaysBetween(a, b Date) int       // b - a in days
func Format(d Date) string            // "YYYY-MM-DD"
```

| Call | Result |
|---|---|
| `IsLeap(2024)`, `IsLeap(1900)`, `IsLeap(2000)` | `true`, `false`, `true` |
| `DaysInMonth(2023, 2)`, `DaysInMonth(2024, 2)`, `DaysInMonth(2024, 13)` | `28`, `29`, `0` |
| `Valid(Date{2023, 2, 29})` | `false` |
| `Compare(Date{2024, 1, 31}, Date{2024, 2, 1})` | `-1` |
| `AddDays(Date{2024, 2, 28}, 1)` | `Date{2024, 2, 29}` |
| `AddDays(Date{2024, 12, 31}, 1)` | `Date{2025, 1, 1}` |
| `AddDays(Date{2024, 3, 1}, -1)` | `Date{2024, 2, 29}` |
| `DaysBetween(Date{2024, 1, 1}, Date{2025, 1, 1})` | `366` |
| `DaysBetween(Date{2024, 3, 1}, Date{2024, 2, 1})` | `-29` |
| `Format(Date{987, 6, 5})` | `"0987-06-05"` |

## Rules
- Leap year: divisible by 4, except centuries, except those divisible by 400.
- `DaysInMonth` returns `0` for a month outside 1–12.
- `Valid`: year 1–9999, month 1–12, day 1–`DaysInMonth`.
- `Compare` orders by year, then month, then day. It works for any `Date`, valid or not. `IsBefore(a, b)` is `Compare(a, b) < 0`. Because `Compare` has the shape `func(a, b Date) int`, `slices.SortFunc(dates, Compare)` sorts dates.
- `AddDays` and `DaysBetween` may assume valid input and results in years 1–9999. They must be fast for large `n` (hundreds of thousands of days).
- **Do not import `time`.** (The tests use it to check your answers.)

# Calendar Switches

Implement:

```go
func DayType(day string) string       // "weekday", "weekend" or "invalid"
func DaysInMonth(month, year int) int // month is 1 (January) … 12 (December)
```

## `DayType`
`day` is a three-letter English abbreviation, **case-sensitive**: `"Mon"`, `"Tue"`, `"Wed"`, `"Thu"`, `"Fri"` are weekdays; `"Sat"`, `"Sun"` are the weekend. Anything else (`"mon"`, `"Monday"`, `""`) is `"invalid"`.

## `DaysInMonth`
| Month(s) | Days |
|---|---|
| 1, 3, 5, 7, 8, 10, 12 | 31 |
| 4, 6, 9, 11 | 30 |
| 2 | 29 in a leap year, else 28 |
| anything else | 0 (invalid month) |

A leap year is divisible by 4, except years divisible by 100 that are not divisible by 400.

## Examples
| Call | Result |
|---|---|
| `DayType("Sat")` | `"weekend"` |
| `DayType("Wed")` | `"weekday"` |
| `DayType("sun")` | `"invalid"` |
| `DaysInMonth(4, 2023)` | `30` |
| `DaysInMonth(2, 2024)` | `29` |
| `DaysInMonth(2, 1900)` | `28` |
| `DaysInMonth(13, 2024)` | `0` |

## Rules
- Use a **value** `switch` with **several values per case** (`case "Sat", "Sun":`) in both functions.
- Use `default` for the invalid inputs.

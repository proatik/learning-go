# Weekday Enum with `iota`

`weekday.go` declares a named type and seven constants, but every constant is `0`.

```go
type Weekday int

const (
	Sunday Weekday = ... // 0
	Monday               // 1
	...
	Saturday             // 6
)
```

## Tasks
1. Rewrite the constant block with **`iota`** so `Sunday == 0`, `Monday == 1`, …, `Saturday == 6`. Write `= iota` only once.
2. Implement:

```go
func IsWeekend(d Weekday) bool // Saturday or Sunday
func Next(d Weekday) Weekday   // the following day; Saturday → Sunday
func Prev(d Weekday) Weekday   // the previous day;  Sunday → Saturday
```

| Call | Result |
|---|---|
| `IsWeekend(Saturday)` | `true` |
| `IsWeekend(Wednesday)` | `false` |
| `Next(Monday)` | `Tuesday` |
| `Next(Saturday)` | `Sunday` |
| `Prev(Sunday)` | `Saturday` |

## Rules
- Use `%` to wrap around: no `if` needed for `Next`/`Prev`.
- Careful: in Go `-1 % 7` is `-1`, not `6`. Find a way to stay non-negative.
- `||` is logical OR: `a || b`.

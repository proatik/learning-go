# Weekday Enum

Give an `iota` enum readable output with `String`, without falling into the recursion trap.

```go
type Weekday int

const (
	Sunday Weekday = iota
	Monday
	Tuesday
	Wednesday
	Thursday
	Friday
	Saturday
)

func (d Weekday) String() string // "Monday", or "Weekday(9)" if unknown
func (d Weekday) Short() string  // "Mon",    or "???" if unknown
func (d Weekday) Next() Weekday  // the following day; Saturday → Sunday
```

## Examples

| Expression | Result |
|---|---|
| `fmt.Sprint(Monday)` | `"Monday"` |
| `fmt.Sprintf("%v and %s", Saturday, Sunday)` | `"Saturday and Sunday"` |
| `fmt.Sprint(Weekday(9))` | `"Weekday(9)"` |
| `fmt.Sprint(Weekday(-1))` | `"Weekday(-1)"` |
| `Friday.Short()` | `"Fri"` |
| `Weekday(7).Short()` | `"???"` |
| `Saturday.Next()` | `Sunday` |
| `Weekday(9).Next()` | `Weekday(9)` |

## Rules
- Valid days are `0` (Sunday) through `6` (Saturday).
- `String` for an invalid day returns `Weekday(N)` with `N` in decimal. Don't call `%v`/`%s` on `d` inside `String`.
- `Short` returns the first three letters of the name.
- `Next` wraps from `Saturday` to `Sunday`; an invalid day is returned unchanged.

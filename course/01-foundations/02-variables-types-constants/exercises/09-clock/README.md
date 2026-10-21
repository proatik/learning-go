# Seconds to Clock Time

Practise integer `/` and `%`, including their sign rules.

Implement:

```go
func Split(total int) (h, m, s int) // hours, minutes, seconds
func Clock(total int) string        // "HH:MM:SS", with "-" for negative durations
func TimeOfDay(total int) string    // seconds after midnight, wrapped into one day
```

| Call | Result |
|---|---|
| `Split(3661)` | `1, 1, 1` |
| `Split(59)` | `0, 0, 59` |
| `Split(-3661)` | `-1, -1, -1` |
| `Clock(3661)` | `"01:01:01"` |
| `Clock(360000)` | `"100:00:00"` |
| `Clock(-90)` | `"-00:01:30"` |
| `TimeOfDay(3600)` | `"01:00:00"` |
| `TimeOfDay(86400 + 61)` | `"00:01:01"` (next day) |
| `TimeOfDay(-60)` | `"23:59:00"` (one minute before midnight) |

## Rules
- `Split` uses only `/` and `%`. Go's natural truncation gives the negative results above. Don't fight it.
- `Clock` pads each part to **at least** two digits with `%02d`. For a negative total, print `-` and then format the absolute value.
- `TimeOfDay` must always land in `00:00:00 … 23:59:59`. Remember `-60 % 86400 == -60` in Go. Turn it into a non-negative remainder first.
- Declare the magic numbers as constants (`secondsPerMinute`, `secondsPerHour`, `secondsPerDay`).

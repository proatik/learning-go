# Clock Arithmetic

A 24-hour wall clock has no date — after `23:59` comes `00:00`. Implement:

```go
func AddMinutes(h, m, delta int) string
func ParseClock(s string) (h, m int, err error)
```

## AddMinutes
Start at `h` hours and `m` minutes, add `delta` minutes (may be negative), and return the time as `"HH:MM"` (always two digits each).

| Call | Result |
|---|---|
| `AddMinutes(10, 0, 3)` | `"10:03"` |
| `AddMinutes(23, 59, 1)` | `"00:00"` |
| `AddMinutes(0, 0, -1)` | `"23:59"` |
| `AddMinutes(1, 0, -3*24*60)` | `"01:00"` |
| `AddMinutes(25, 160, 0)` | `"03:40"` |
| `AddMinutes(-1, 0, 0)` | `"23:00"` |

- `h` and `m` themselves may be out of range or negative: everything is just minutes (`h*60 + m + delta`) wrapped onto the clock.
- Remember that `%` can return a negative number in Go.

## ParseClock
Parse a strict `"HH:MM"` string.

| Input | Result |
|---|---|
| `"09:05"` | `9, 5, nil` |
| `"23:59"` | `23, 59, nil` |
| `"00:00"` | `0, 0, nil` |

Return a non-nil error (and `0, 0`) unless **all** of these hold:
- `s` is exactly 5 bytes long and `s[2] == ':'`;
- the other four bytes are ASCII digits `0`–`9` (so `"+1:00"`, `" 9:00"`, `"9:00"` are invalid);
- the hour is `0`–`23` and the minute is `0`–`59`.

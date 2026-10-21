# Human-Readable Duration

Implement:

```go
func HumanDuration(seconds int) string
```

It turns a number of seconds into days, hours, minutes and seconds, largest unit first.

| Input | Output |
|---|---|
| `0` | `"0s"` |
| `59` | `"59s"` |
| `60` | `"1m"` |
| `3605` | `"1h 5s"` |
| `3661` | `"1h 1m 1s"` |
| `86400` | `"1d"` |
| `1_000_000` | `"11d 13h 46m 40s"` |
| `-90` | `"-1m 30s"` |

## Rules
- Units: `d` (86400 s), `h` (3600 s), `m` (60 s), `s`. There is no larger unit: `100` days is `"100d"`.
- Parts are separated by **one space**; **omit** parts that are zero. No leading or trailing spaces.
- `0` is the only input that prints a `0` part: `"0s"`.
- A negative input gets a single leading `-` followed by the format of its absolute value.
- You may assume `seconds != math.MinInt`.
- Declare the unit sizes as constants.

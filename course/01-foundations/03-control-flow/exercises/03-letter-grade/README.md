# Letter Grade

Implement:

```go
func Grade(score int) string
```

| Score | Grade |
|---|---|
| 90–100 | `"A"` |
| 80–89 | `"B"` |
| 70–79 | `"C"` |
| 60–69 | `"D"` |
| 0–59 | `"F"` |
| below 0 or above 100 | `"invalid"` |

## Examples
| Call | Result |
|---|---|
| `Grade(95)` | `"A"` |
| `Grade(80)` | `"B"` |
| `Grade(59)` | `"F"` |
| `Grade(101)` | `"invalid"` |
| `Grade(-1)` | `"invalid"` |

## Rules
- Use a single **tagless** `switch { ... }` — no `if`.
- Cases are checked top to bottom, so each case only needs a lower bound once invalid scores are handled first.

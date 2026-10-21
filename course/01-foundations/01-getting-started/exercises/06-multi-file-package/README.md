# A Package in Two Files

Package `temperature` is split across two files:

- `format.go`: **already done**. `Describe` calls `CToF`.
- `convert.go`: **your job**. Implement:

```go
func CToF(c float64) float64 // F = C × 9/5 + 32
func FToC(f float64) float64 // C = (F − 32) × 5/9
```

Every file in a directory shares one package scope, so `format.go` can call functions from `convert.go` without importing anything.

## Examples
| Call | Result |
|---|---|
| `CToF(100)` | `212` |
| `FToC(32)` | `0` |
| `Describe(37)` | `"37.0°C = 98.6°F"` |

## Watch out
`5/9` written alone is a **constant expression of two untyped integers**, so it is evaluated with integer division and equals `0`. Make sure a `float64` operand is involved first, e.g. `(f - 32) * 5 / 9`.

# Temperature Types

Methods on named `float64` types: conversions, a predicate, and `String` for printing.

```go
type Celsius float64
type Fahrenheit float64

func (c Celsius) ToFahrenheit() Fahrenheit // F = C*9/5 + 32
func (f Fahrenheit) ToCelsius() Celsius    // C = (F-32)*5/9
func (c Celsius) IsFreezing() bool         // at or below 0°C
func (c Celsius) String() string           // "21.5°C"
func (f Fahrenheit) String() string        // "70.7°F"
```

## Examples

| Code | Result |
|---|---|
| `Celsius(100).ToFahrenheit()` | `212` |
| `Fahrenheit(32).ToCelsius()` | `0` |
| `Celsius(-40).ToFahrenheit()` | `-40` |
| `Celsius(0).IsFreezing()` | `true` |
| `Celsius(21.5).String()` | `"21.5°C"` |
| `Fahrenheit(70.7).String()` | `"70.7°F"` |
| `fmt.Println(Celsius(37))` | `37.0°C` |

## Rules
- `String` uses **one decimal place** followed by `°C` or `°F` (the degree sign is the rune `'°'`).
- `IsFreezing` is true for `c <= 0`.
- Don't call `fmt.Sprintf("%v", c)` inside `String` (infinite recursion): convert to `float64` first.

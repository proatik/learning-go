# Describe Anything

Write one function that accepts **any** value and describes it, using a type switch.

```go
func Describe(v any) string
```

## Examples

| Call | Result |
|---|---|
| `Describe(nil)` | `"nil"` |
| `Describe(42)` | `"int 42"` |
| `Describe(int64(-7))` | `"sized int -7"` |
| `Describe(uint8(255))` | `"unsigned 255"` |
| `Describe(3.5)` | `"float64 3.5"` |
| `Describe("héllo")` | `"string \"héllo\" with 5 runes"` |
| `Describe(true)` | `"bool true"` |
| `Describe([]int{1, 2, 3})` | `"[]int with 3 elements"` |
| `Describe(map[string]int{"a": 1})` | `"map[string]int with 1 keys"` |
| `Describe(Celsius(21))` (a type with `String() string` returning `"21°C"`) | `"stringer 21°C"` |
| `Describe(struct{}{})` | `"unknown struct {}"` |

## Rules
Check the cases in this order; the first match wins:

| Dynamic type | Output |
|---|---|
| nil interface | `nil` |
| `int` | `int <v>` |
| `int8`, `int16`, `int32`, `int64` | `sized int <v>` |
| `uint`, `uint8`, `uint16`, `uint32`, `uint64` | `unsigned <v>` |
| `float64` | `float64 <v>` (format with `%v`) |
| `string` | `string <%q of v> with <rune count> runes` |
| `bool` | `bool <v>` |
| `[]int` | `[]int with <len> elements` (a nil `[]int` has 0 elements: it is **not** the nil case) |
| `map[string]int` | `map[string]int with <len> keys` |
| `fmt.Stringer` | `stringer <result of String()>` |
| anything else | `unknown <%T of v>` |

- Use a single `switch x := v.(type)` statement.
- `<v>` means `fmt.Sprint`/`%v` formatting of the value.

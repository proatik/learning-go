# Zero Values

Implement:

```go
func PrintZeroValues()
```

Declare **one variable of each type** `int`, `int8`, `uint`, `float64`, `string`, `bool`, `rune`, `byte` — in that order, with `var` and **no value** — then print each one's type and value:

```
int = 0
int8 = 0
uint = 0
float64 = 0
string = ""
bool = false
int32 = 0
uint8 = 0
```

## Rules
- Use `var x T` declarations (a grouped `var ( ... )` block is nice). No literals like `0` or `""`.
- Print with `fmt.Printf("%T = %v\n", x, x)`. For the string use `%q` instead of `%v` so the empty value is visible.
- The last two lines are not typos: `rune` is an alias for `int32` and `byte` for `uint8`, and `%T` shows the real type.

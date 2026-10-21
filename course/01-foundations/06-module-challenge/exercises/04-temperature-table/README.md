# Temperature Table

Write a **program** that prints a Celsius → Fahrenheit → Kelvin conversion table.

```go
type Celsius float64
type Fahrenheit float64
type Kelvin float64

func CToF(c Celsius) Fahrenheit // F = C×9/5 + 32
func CToK(c Celsius) Kelvin     // K = C + 273.15
func main()
```

When run, the program must print exactly:

```
Celsius  Fahrenheit   Kelvin
-------  ----------  -------
  -40.0       -40.0   233.15
  -20.0        -4.0   253.15
    0.0        32.0   273.15
   20.0        68.0   293.15
   40.0       104.0   313.15
   60.0       140.0   333.15
   80.0       176.0   353.15
  100.0       212.0   373.15
8 rows
```

## Rules
- Columns are **right-aligned** with widths 7, 10 and 7, separated by **two spaces**. Lines have no trailing spaces.
- Celsius and Fahrenheit use 1 decimal place; Kelvin uses 2.
- Generate the rows with a loop from `-40` to `100` in steps of `20` — declare the start, end and step as constants. Don't hard-code the rows.
- The header line is printed with the same widths (`%7s`, `%10s`, `%7s`); the dashes line is built with `strings.Repeat`.
- The last line counts the rows you printed (don't hard-code `8`).
- `CToF` and `CToK` are tested separately, for any input.

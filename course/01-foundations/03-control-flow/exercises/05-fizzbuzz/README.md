# FizzBuzz

Implement:

```go
func FizzBuzz(n int) string // the FizzBuzz word for one number
func PrintFizzBuzz(n int)    // prints FizzBuzz(1) … FizzBuzz(n), one per line
```

## `FizzBuzz`
| `n` is divisible by | Result |
|---|---|
| both 3 and 5 | `"FizzBuzz"` |
| 3 only | `"Fizz"` |
| 5 only | `"Buzz"` |
| neither | the number itself, e.g. `"7"` |

Zero is divisible by everything, so `FizzBuzz(0)` is `"FizzBuzz"`. Negative numbers follow the same rules: `FizzBuzz(-3)` is `"Fizz"`, `FizzBuzz(-7)` is `"-7"`.

## `PrintFizzBuzz`
`PrintFizzBuzz(5)` prints:

```
1
2
Fizz
4
Buzz
```

If `n < 1` it prints nothing.

## Rules
- `FizzBuzz`: tagless `switch`. Convert a number to a string with `strconv.Itoa(n)`.
- `PrintFizzBuzz`: loop with `for i := range n` (remember it yields `0 … n-1`) and reuse `FizzBuzz`.

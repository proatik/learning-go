# defer: Order and Named Results

Two functions, two `defer` rules.

```go
func Stack(n int)
func ParseAge(s string) (age int, err error)
```

## 1. `Stack`: deferred calls run last-in, first-out

`Stack(n)` prints `push 1` … `push n`. For each push it also **defers** printing the matching `pop`. Because deferred calls run in LIFO order when `Stack` returns, the output for `Stack(3)` is:

```
push 1
push 2
push 3
pop 3
pop 2
pop 1
```

`Stack(0)` prints nothing. Write exactly **one** loop: the pops must come from `defer`, not from a second loop.

## 2. `ParseAge`: a deferred closure edits the results

`ParseAge` converts `s` with `strconv.Atoi` and accepts ages `0` to `150` inclusive.

| Call | `age` | `err.Error()` |
|---|---|---|
| `ParseAge("42")` | `42` | (nil) |
| `ParseAge("150")` | `150` | (nil) |
| `ParseAge("200")` | `0` | `ParseAge("200"): age out of range` |
| `ParseAge("-1")` | `0` | `ParseAge("-1"): age out of range` |
| `ParseAge("abc")` | `0` | `ParseAge("abc"): ` followed by the `strconv` error text |

Rules:
- The range error's own text is exactly `age out of range`.
- Add the `ParseAge("…"): ` prefix in **one place**: a `defer`red function literal, registered at the top of the function, that runs only when `err != nil`. Build it with `fmt.Errorf("ParseAge(%q): %v", s, err)`.
- The same deferred closure sets `age = 0` whenever there is an error.

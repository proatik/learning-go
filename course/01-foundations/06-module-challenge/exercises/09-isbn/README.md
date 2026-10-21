# ISBN-10 Check and Convert

An ISBN-10 has 9 digits plus a **check character** that is a digit or `X` (meaning 10). It is valid when

```
(d1×10 + d2×9 + d3×8 + d4×7 + d5×6 + d6×5 + d7×4 + d8×3 + d9×2 + d10×1) % 11 == 0
```

Implement:

```go
func IsValid(isbn string) bool
func ToISBN13(isbn string) (string, error)
```

| Call | Result |
|---|---|
| `IsValid("3-598-21508-8")` | `true` |
| `IsValid("3598215088")` | `true` |
| `IsValid("3-598-21507-X")` | `true` |
| `IsValid("3-598-21508-9")` | `false` (wrong check digit) |
| `IsValid("3-598-2X507-9")` | `false` (`X` only allowed last) |
| `ToISBN13("0-306-40615-2")` | `"9780306406157", nil` |
| `ToISBN13("3-598-21508-9")` | `"", error` |

## IsValid rules
- Hyphens (`-`) are ignored wherever they appear. Any other character that is not a digit (including spaces and lowercase `x`) makes the ISBN invalid.
- After removing hyphens there must be exactly **10** characters: 9 ASCII digits and then a digit or uppercase `X`.

## ToISBN13 rules
- If `isbn` is not a valid ISBN-10, return `""` and a non-nil error.
- Otherwise: take `"978"` + the first 9 digits (12 digits), then append an ISBN-13 check digit. Multiply the 12 digits alternately by **1, 3, 1, 3, ...**, sum them, and the check digit is `(10 - sum%10) % 10`.
- The result is 13 digits with no hyphens.

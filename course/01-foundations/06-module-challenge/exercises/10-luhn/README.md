# Luhn Checksum

The [Luhn algorithm](https://en.wikipedia.org/wiki/Luhn_algorithm) validates credit-card numbers, IMEIs and many other IDs. Implement:

```go
func Valid(number string) bool
func CheckDigit(payload string) (int, error)
```

## `Valid`: algorithm
1. Spaces are allowed and ignored. Any other non-digit character makes the number **invalid**.
2. After removing spaces, fewer than **2** digits is **invalid**.
3. Walk the digits **from right to left**. Double every **second** digit (the 2nd, 4th, … from the right). If doubling gives more than 9, subtract 9.
4. Sum all the digits. The number is valid when `sum % 10 == 0`.

Example: `79927398713`

| digit (right → left) | 3 | 1 | 7 | 8 | 9 | 3 | 7 | 2 | 9 | 9 | 7 |
|---|---|---|---|---|---|---|---|---|---|---|---|
| double every 2nd | 3 | **2** | 7 | **16→7** | 9 | **6** | 7 | **4** | 9 | **18→9** | 7 |

Sum = 3+2+7+7+9+6+7+4+9+9+7 = **70**, and 70 % 10 == 0, so the number is valid.

## `CheckDigit`
Returns the digit `d` (0–9) such that `payload + d` is valid. `payload` follows the same rules (spaces ignored, other characters → error) and must contain at least **1** digit, otherwise return an error.

| Call | Result |
|---|---|
| `Valid("4539 3195 0343 6467")` | `true` |
| `Valid("8273 1232 7352 0569")` | `false` |
| `Valid("059")` | `true` |
| `Valid("0")` | `false` (too short) |
| `Valid("055-444-285")` | `false` (`-` not allowed) |
| `CheckDigit("7992739871")` | `3, nil` |
| `CheckDigit("")` | `0, error` |

## Rules
- Work on the string directly: `range` over it, or index bytes from the end. No `strconv` needed: a digit byte `c` has the value `int(c - '0')`.
- Reuse your own code: `CheckDigit` can try all ten digits with `Valid`, or compute the digit directly.

# Number to Words

Spell out a number in English. Implement:

```go
func NumberToWords(n int) (string, error)
```

| Input | Output |
|---|---|
| `0` | `"zero"` |
| `7` | `"seven"` |
| `13` | `"thirteen"` |
| `40` | `"forty"` |
| `42` | `"forty-two"` |
| `100` | `"one hundred"` |
| `123` | `"one hundred twenty-three"` |
| `1000` | `"one thousand"` |
| `1001` | `"one thousand one"` |
| `20020` | `"twenty thousand twenty"` |
| `12345` | `"twelve thousand three hundred forty-five"` |
| `999999` | `"nine hundred ninety-nine thousand nine hundred ninety-nine"` |
| `-1`, `1_000_000` | `"", error` |

## Rules
- Valid range is `0`–`999_999`; anything else returns `""` and a non-nil error.
- All lowercase. Words are separated by a single space; tens and ones between 21 and 99 are joined by a hyphen (`"ninety-nine"`).
- No `"and"`, no commas, and no zero words except for `0` itself (`101` is `"one hundred one"`, `1000` is `"one thousand"`).
- Spelling: `forty` (not fourty), `fifteen`, `eighteen`, `nineteen`.
- **Slice preview:** a `[]string{"zero", "one", ...}` literal is an indexable lookup table: `ones[7]` is `"seven"`.
- Tip: write a helper that spells `1`–`999` and use it twice (thousands, then the rest).

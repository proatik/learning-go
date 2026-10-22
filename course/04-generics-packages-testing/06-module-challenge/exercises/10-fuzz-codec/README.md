# Bug Hunt: Fuzzing a Codec

A run-length codec is already written. Prove it correct with the **whole testing toolkit** in `rle_test.go`: fuzz tests with well-chosen seeds, a table test, an example and a benchmark. Together they must catch all **7** planted bugs.

```go
// Provided in check_test.go (read it in the Tests tab).
var ErrMalformed = errors.New("rle: malformed input")

func Encode(s string) string          // "aaab" → "3:a1:b"; defined for valid UTF-8 only
func Decode(s string) (string, error) // reverses Encode; ErrMalformed otherwise; never panics
```

## Format

| Text | Encoded |
|---|---|
| `""` | `""` |
| `"aaab"` | `"3:a1:b"` |
| `"héé"` | `"1:h2:é"` — counts **runes**, not bytes |
| `"xxxxxxxxxxxx"` (12) | `"12:x"` |
| `"::"` | `"2::"` — any rune, even `:` or a digit, may follow the colon |

`Decode` returns `ErrMalformed` for: a missing count or colon (`"a"`, `":a"`), a count with no rune after it (`"3:"`), a count below 1 or above 1,000,000 (`"0:a"`), a non-number (`"x:a"`), trailing junk (`"2:a3"`), or invalid UTF-8 after the colon (`"1:\xff"`).

## Rules
- Write all of these in `rle_test.go` — the checker requires each kind:
  - at least one **`Fuzz…`** test checking the property `Decode(Encode(s)) == s`. `Encode` is only defined for valid UTF-8, so `t.Skip` when `!utf8.ValidString(s)`.
  - a fuzz test that feeds `Decode` **arbitrary** input: it must never panic, and only ever return `ErrMalformed` as an error.
  - a **`Test…`** with a table of malformed inputs (check `errors.Is(err, ErrMalformed)`).
  - an **`Example…`** with an `// Output:` comment.
  - a **`Benchmark…`** using `for b.Loop()`.
- The bug hunt runs fuzz tests on their **seed corpus only** (the `f.Add` values), so your seeds must reach every tricky case. Think: empty, long runs, case, multi-byte runes, `:` and digits as text, broken encodings.
- All of your tests must pass on the correct code; the bug hunt then runs them once per bug and lists survivors as `Bug "…" = not caught`.

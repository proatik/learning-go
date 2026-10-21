# Caesar Cipher

Implement:

```go
func Encrypt(s string, shift int) string
func Decrypt(s string, shift int) string
```

`Encrypt` moves every ASCII letter `shift` places along the alphabet, **wrapping around** (`z` + 1 = `a`). `Decrypt` undoes it.

| Call | Result |
|---|---|
| `Encrypt("abc", 1)` | `"bcd"` |
| `Encrypt("xyz", 3)` | `"abc"` |
| `Encrypt("Hello, World!", 13)` | `"Uryyb, Jbeyq!"` |
| `Encrypt("abc", -1)` | `"zab"` |
| `Encrypt("abc", 27)` | `"bcd"` |
| `Encrypt("Café 2027", 1)` | `"Dbgé 2027"` |
| `Decrypt("bcd", 1)` | `"abc"` |

## Rules
- Only `a`–`z` and `A`–`Z` change; case is preserved.
- Everything else (spaces, digits, punctuation, `é`, `世`) is copied unchanged.
- `shift` may be any `int`: negative, zero or larger than 26.
- `Decrypt(Encrypt(s, k), k) == s` for every `s` and `k`.

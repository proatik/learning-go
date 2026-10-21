# Strings, Runes & Formatting

## 1. What a string is

A `string` is an **immutable, read-only sequence of bytes**. Go source is UTF-8, so string literals hold UTF-8 text — but a string may contain *any* bytes.

```go
s := "héllo"
fmt.Println(len(s)) // 6 — len counts BYTES: 'é' takes 2 bytes in UTF-8
fmt.Println(s[0])   // 104 — indexing gives a byte (uint8), not a character
fmt.Println(s[1])   // 195 — the first byte of 'é'
// s[0] = 'H'       // compile error: strings are immutable
```

| Term | Type | Meaning |
|---|---|---|
| byte | `byte` (= `uint8`) | one 8-bit unit of the string |
| rune | `rune` (= `int32`) | one Unicode **code point**, e.g. `'é'` = U+00E9 |
| character | — | what a human sees; may be several runes (e.g. `e` + combining accent) |

UTF-8 uses **1–4 bytes per rune**: ASCII = 1, `é` = 2, `世` = 3, most emoji = 4.

```go
import "unicode/utf8"

s := "Go, 世界"
len(s)                    // 10 bytes
utf8.RuneCountInString(s) // 6 runes
utf8.RuneLen('世')         // 3
utf8.ValidString(s)       // true
```

## 2. Iterating: bytes vs runes

```go
s := "aé世"

for i := 0; i < len(s); i++ { // byte by byte
	fmt.Printf("%d:%x ", i, s[i]) // 0:61 1:c3 2:a9 3:e4 4:b8 5:96
}

for i, r := range s { // rune by rune — range DECODES UTF-8
	fmt.Printf("%d:%c ", i, r) // 0:a 1:é 3:世
}
```

- `for i, r := range s` yields the **byte index** where each rune starts and the **rune** itself. Note the jumps: 0, 1, 3.
- An invalid UTF-8 byte yields `utf8.RuneError` (U+FFFD, `�`) and advances by 1 byte.
- Need just the runes? `for _, r := range s`.

## 3. Conversions

```go
b := []byte("héllo") // [104 195 169 108 108 111] — a copy of the bytes
r := []rune("héllo") // [104 233 108 108 111]     — one element per rune
string(b)            // "héllo"
string(r)            // "héllo"
string('世')          // "世" — a rune converts to its UTF-8 encoding
len(r)               // 5 — same as utf8.RuneCountInString
```

> `[]byte` and `[]rune` are **slices** — growable lists indexed with `x[i]`, with `len(x)`. Unlike a string you can modify them (`r[0] = 'H'`). Slices get their own module next.

Converting to `[]rune` is the simple way to work **per character position**: reverse, take the n-th rune, swap runes.

## 4. Substrings (slicing)

`s[a:b]` is the substring from **byte** `a` up to (not including) byte `b`. No copy is made.

```go
s := "Hello, 世界"
s[0:5]  // "Hello"
s[7:]   // "世界"     (omitted end = len(s))
s[:5]   // "Hello"    (omitted start = 0)
s[7:8]  // "\xe4"     — cuts 世 in half: invalid UTF-8!
```

Offsets returned by `strings.Index` and by `range` are byte offsets, so they are always safe to slice with.

## 5. Literals

```go
"tab\there\n"      // interpreted: escapes \n \t \\ \" \u00e9 \x41 work
`C:\new\folder`    // raw: backslashes are literal, no escapes at all
`line 1
line 2`            // raw strings may span lines (carriage returns are dropped)
'A'                // rune literal (int32 value 65), NOT a string
```

Raw strings are ideal for regular expressions, Windows paths, JSON and multi-line text. They can't contain a backtick.

## 6. Comparing

```go
"go" == "go"                    // true — byte-for-byte equality
"apple" < "banana"              // true — lexicographic by bytes
"Zebra" < "apple"               // true! 'Z' (90) < 'a' (97)
strings.EqualFold("Go", "GO")   // true — case-insensitive, Unicode-aware
strings.Compare(a, b)           // -1, 0, +1 (rarely needed; use == and <)
```

## 7. Building strings

Strings are immutable, so `s += x` allocates a **new** string and copies everything each time. In a loop that's O(n²).

```go
// ❌ slow for many pieces
s := ""
for i := range 1000 {
	s += strconv.Itoa(i)
}

// ✅ strings.Builder appends into a growing buffer
var b strings.Builder // zero value is ready to use
for i := range 1000 {
	b.WriteString(strconv.Itoa(i))
}
b.WriteByte('!')  // append one byte
b.WriteRune('世') // append one rune (UTF-8 encoded)
s := b.String()   // the result
b.Len()           // bytes written so far
```

`+` is perfectly fine for a few pieces (`"Hello, " + name + "!"`). Use a `Builder` in loops. Don't copy a `Builder` after writing to it.

## 8. The `strings` package

All functions return **new** strings (strings are immutable). Indexes are byte offsets; "not found" is `-1`.

| Function | Example | Result |
|---|---|---|
| `Contains(s, sub)` | `Contains("seafood", "foo")` | `true` |
| `ContainsRune(s, r)` | `ContainsRune("aé", 'é')` | `true` |
| `HasPrefix(s, p)` / `HasSuffix(s, p)` | `HasSuffix("main.go", ".go")` | `true` |
| `Index(s, sub)` | `Index("chicken", "ken")` | `4` |
| `LastIndex(s, sub)` | `LastIndex("go gopher", "go")` | `3` |
| `Count(s, sub)` | `Count("cheese", "e")` | `3` |
| `ToUpper` / `ToLower` | `ToUpper("héllo")` | `"HÉLLO"` |
| `TrimSpace(s)` | `TrimSpace("  hi \n")` | `"hi"` |
| `Trim(s, cutset)` | `Trim("!!hi!?", "!?")` | `"hi"` |
| `TrimPrefix` / `TrimSuffix` | `TrimSuffix("a.go", ".go")` | `"a"` |
| `Repeat(s, n)` | `Repeat("ab", 3)` | `"ababab"` |
| `Replace(s, old, new, n)` | `Replace("aaa", "a", "b", 2)` | `"bba"` |
| `ReplaceAll(s, old, new)` | `ReplaceAll("aaa", "a", "b")` | `"bbb"` |
| `Fields(s)` | `Fields("  a b\tc ")` | `["a" "b" "c"]` |
| `Split(s, sep)` | `Split("a,b,,c", ",")` | `["a" "b" "" "c"]` |
| `Join(elems, sep)` | `Join([]string{"a", "b"}, "-")` | `"a-b"` |
| `Cut(s, sep)` | `Cut("k=v", "=")` | `"k", "v", true` |
| `CutPrefix(s, p)` | `CutPrefix("v1.2", "v")` | `"1.2", true` |
| `CutSuffix(s, p)` | `CutSuffix("a.go", ".go")` | `"a", true` |
| `CutLast(s, sep)` *(Go 1.27)* | `CutLast("a.tar.gz", ".")` | `"a.tar", "gz", true` |
| `EqualFold(a, b)` | `EqualFold("Go", "GO")` | `true` |

`Cut` is the idiomatic way to split **once** around a separator:

```go
key, value, ok := strings.Cut("name=Gopher", "=")
if !ok {
	return errors.New("missing '='")
}
// key == "name", value == "Gopher"

before, after, found := strings.CutLast("archive.tar.gz", ".") // "archive.tar", "gz", true
```

Gotchas:
- `Trim(s, cutset)` removes any **characters** in `cutset` from both ends — it is not "remove this prefix". Use `TrimPrefix` for that.
- `Split("", ",")` returns `[""]` (one empty element), but `Fields("")` returns an empty list.
- `strings.Title` is deprecated (it mishandles Unicode word boundaries) — build title case yourself.

## 9. `unicode`: classifying runes

| Function | `true` for |
|---|---|
| `unicode.IsLetter(r)` | `'a'`, `'É'`, `'世'` |
| `unicode.IsDigit(r)` | `'0'`–`'9'` (and other scripts' decimal digits) |
| `unicode.IsSpace(r)` | `' '`, `'\t'`, `'\n'`, `'\r'`, U+00A0 … |
| `unicode.IsUpper(r)` / `IsLower(r)` | `'A'` / `'a'` |
| `unicode.IsPunct(r)` | `'!'`, `','`, `'.'` |
| `unicode.ToUpper(r)` / `ToLower(r)` | returns the converted rune: `ToUpper('é')` = `'É'` |

Arithmetic on runes works because they are integers: `'a' + 2 == 'c'`, `r - 'a'` is a letter's position in the alphabet, `r - '0'` is a digit's value.

## 10. `strconv`: strings ↔ numbers

| Function | Example | Result |
|---|---|---|
| `Itoa(i)` | `Itoa(42)` | `"42"` |
| `Atoi(s)` | `Atoi("42")` | `42, nil` |
| `Atoi(s)` (bad input) | `Atoi("4x")` | `0, strconv.Atoi: parsing "4x": invalid syntax` |
| `ParseInt(s, base, bits)` | `ParseInt("ff", 16, 64)` | `255, nil` |
| `ParseInt(s, 0, 64)` | `ParseInt("0b101", 0, 64)` | `5, nil` — base 0 reads the `0b`/`0o`/`0x` prefix |
| `FormatInt(i, base)` | `FormatInt(255, 2)` | `"11111111"` |
| `FormatInt(i, base)` | `FormatInt(-255, 16)` | `"-ff"` |
| `ParseBool(s)` | `ParseBool("true")` | `true, nil` (accepts `1 t T TRUE true True 0 f F FALSE false False`) |
| `ParseFloat(s, 64)` | `ParseFloat("3.14", 64)` | `3.14, nil` |
| `Quote(s)` | `Quote("hi\n")` | `"hi\n"` — quotes and escapes included (same as `%q`) |

`FormatInt` takes an `int64` and bases 2–36: `strconv.FormatInt(int64(n), 2)`. Always check the error from `Atoi`/`Parse*`.

## 11. `fmt` verbs

`Printf` prints; `Sprintf` returns the string; `Errorf` builds an `error`. All take the same format language.

| Verb | Meaning | Example | Output |
|---|---|---|---|
| `%v` | default format | `Sprintf("%v", 3.5)` | `3.5` |
| `%+v` | `%v` plus field names (structs, later) | | |
| `%T` | type of the value | `Sprintf("%T", 'a')` | `int32` |
| `%d` | integer, base 10 | `Sprintf("%d", 42)` | `42` |
| `%5d` | width 5, right-aligned | `Sprintf("[%5d]", 42)` | `[   42]` |
| `%-5d` | width 5, left-aligned | `Sprintf("[%-5d]", 42)` | `[42   ]` |
| `%05d` | pad with zeros | `Sprintf("%05d", 42)` | `00042` |
| `%+d` | always show sign | `Sprintf("%+d", 42)` | `+42` |
| `%b` | binary | `Sprintf("%b", 5)` | `101` |
| `%o` / `%O` | octal / with `0o` | `Sprintf("%o", 8)` | `10` |
| `%x` / `%X` | hex lower / upper | `Sprintf("%x %X", 255, 255)` | `ff FF` |
| `%#x` | hex with `0x` | `Sprintf("%#x", 255)` | `0xff` |
| `%x` on a string | hex of each byte | `Sprintf("% x", "é")` | `c3 a9` |
| `%c` | the character of a rune | `Sprintf("%c", 'é')` | `é` |
| `%U` | Unicode notation | `Sprintf("%U", 'é')` | `U+00E9` |
| `%q` | quoted string / rune | `Sprintf("%q %q", "hi", 'x')` | `"hi" 'x'` |
| `%s` | string | `Sprintf("%s", "go")` | `go` |
| `%10s` / `%-10s` | width 10, right / left | `Sprintf("[%-6s]", "go")` | `[go    ]` |
| `%.3s` | at most 3 runes | `Sprintf("%.3s", "héllo")` | `hél` |
| `%f` | float, 6 decimals | `Sprintf("%f", 3.14159)` | `3.141590` |
| `%.2f` | 2 decimals (rounded) | `Sprintf("%.2f", 3.14159)` | `3.14` |
| `%8.3f` | width 8, 3 decimals | `Sprintf("[%8.3f]", 3.14159)` | `[   3.142]` |
| `%e` | scientific | `Sprintf("%e", 1234.5678)` | `1.234568e+03` |
| `%g` | compact float | `Sprintf("%g", 1234.5678)` | `1234.5678` |
| `%t` | bool | `Sprintf("%t", true)` | `true` |
| `%p` | pointer address (pointers come later) | | `0xc000012345` |
| `%%` | a literal `%` | `Sprintf("%d%%", 50)` | `50%` |

Width and precision for strings count **runes**, not bytes, so `%-6s` lines up `"Café"` correctly.

**Explicit argument indexes** — `%[n]` picks the n-th argument (1-based) and following verbs continue from there:

```go
fmt.Printf("%[2]s %[1]s\n", "world", "hello") // hello world
fmt.Printf("%[1]d %[1]x %[1]b\n", 10)          // 10 a 1010
```

**Errors** — `Errorf` uses the same verbs; `%q` is great for showing bad input:

```go
return fmt.Errorf("invalid line %q: missing '='", line) // invalid line "abc": missing '='
```

`go vet` checks every `Printf`-style call: wrong verb for the type, missing or extra arguments.

## 12. Common mistakes

| Mistake | Why it's wrong | Fix |
|---|---|---|
| `s[i]` to get the i-th character | gives the i-th **byte** | `[]rune(s)[i]`, or `range s` |
| `len(s)` as "number of characters" | counts bytes | `utf8.RuneCountInString(s)` |
| `s[:3]` to take 3 characters | may cut a rune in half | convert to `[]rune` first, or slice at an offset from `range` / `Index` |
| `string(65)` to get `"65"` | gives `"A"` (rune 65); `go vet` warns, so `go test` fails | `strconv.Itoa(65)` |
| `s += x` inside a loop | copies the whole string every time | `strings.Builder` |
| `strings.Trim(s, "abc")` to remove prefix `"abc"` | removes any of the chars a, b, c | `strings.TrimPrefix` |
| Comparing case-insensitively with `ToLower(a) == ToLower(b)` | allocates; misses some Unicode folds | `strings.EqualFold(a, b)` |
| Ignoring the `error` from `strconv.Atoi` | `"4x"` silently becomes `0` | `n, err := strconv.Atoi(s); if err != nil { ... }` |

## Key takeaways

- A string is immutable bytes (usually UTF-8). `len` and `s[i]` work on **bytes**.
- `for i, r := range s` decodes **runes**; `[]rune(s)` gives per-character positions.
- Slicing `s[a:b]` uses byte offsets — only slice at rune boundaries.
- Build strings in loops with `strings.Builder`; split once with `strings.Cut`.
- `strconv` converts numbers ↔ text (check the errors); `unicode` classifies runes.
- Learn the `fmt` verbs and width/precision — `go vet` catches mismatches.

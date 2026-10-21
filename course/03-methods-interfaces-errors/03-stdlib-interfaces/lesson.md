# Interfaces in the Standard Library

The standard library is built around a handful of **tiny interfaces** (one or two methods each). Implement one method and your type plugs into `fmt`, `io`, `bufio`, `sort` and every package that accepts them.

| Interface | Package | Method(s) | You get |
|---|---|---|---|
| `Stringer` | `fmt` | `String() string` | Custom output for `Println`, `%v`, `%s` |
| `error` | built-in | `Error() string` | The error-handling protocol (next lesson) |
| `Reader` | `io` | `Read(p []byte) (n int, err error)` | Any byte source: strings, buffers, files, network |
| `Writer` | `io` | `Write(p []byte) (n int, err error)` | Any byte sink: buffers, stdout, files, hashes |
| `Interface` | `sort` | `Len`, `Less`, `Swap` | `sort.Sort` (legacy) |

## 1. `fmt.Stringer`

```go
type Stringer interface {
	String() string
}
```

`fmt` checks whether an operand implements `Stringer` and, for `%v`, `%s`, `Print` and `Println`, calls `String()` instead of printing the raw fields.

```go
type Point struct{ X, Y int }

func (p Point) String() string { return fmt.Sprintf("(%d, %d)", p.X, p.Y) }

p := Point{1, 2}
fmt.Println(p, &p)                             // (1, 2) (1, 2)
fmt.Printf("%v %s %+v\n", p, p, []Point{p})    // (1, 2) (1, 2) [(1, 2)]
fmt.Printf("%d\n", p)                          // {1 2}  ← %d is not a string verb: String() is ignored
```

- Works inside slices, maps and **exported** struct fields too: `fmt` calls `String()` on each element (it can't call methods through unexported fields).
- A value-receiver `String` is in the method set of **both** `T` and `*T`, so both print nicely.

### The pointer-receiver gotcha

```go
type Temp float64

func (t *Temp) String() string { return fmt.Sprintf("%.1f°C", float64(*t)) }

t := Temp(21.5)
fmt.Println(t, &t) // 21.5 21.5°C
```

- `String` is only in the method set of `*Temp`, so a plain `Temp` value **does not** implement `Stringer` and prints raw.
- Fix: use a **value receiver** for `String` (it only reads), unless the type must never be copied.

### The infinite-recursion trap

```go
type Weekday int

func (d Weekday) String() string {
	return fmt.Sprintf("%v", d) // ❌ %v calls d.String() → calls Sprintf → … stack overflow
}
```

- `go vet` reports: `fmt.Sprintf format %v with arg d causes recursive (…Weekday).String method call`.
- Fix: format the **underlying value** with a non-string verb or a conversion:

```go
func (d Weekday) String() string {
	return fmt.Sprintf("Weekday(%d)", int(d)) // ✅ int has no String method
}
```

### Related formatting interfaces (mention only)

- `fmt.Formatter` — `Format(f fmt.State, verb rune)`: full control over every verb and flag. Rarely needed.
- `fmt.GoStringer` — `GoString() string`: used by `%#v`.
- `encoding.TextMarshaler` (`MarshalText() ([]byte, error)`) and `json.Marshaler` (`MarshalJSON() ([]byte, error)`) let a type control how it is encoded as text / JSON.

## 2. The `error` interface (shape only)

```go
type error interface {
	Error() string
}
```

- `error` is a predeclared interface, exactly like `Stringer` but for failures.
- A `nil` error means success. Functions return it **last**: `(n int, err error)`.
- Creating, wrapping and inspecting errors is the whole next lesson. Here you only need: `if err != nil { ... }` and the special value `io.EOF`.

## 3. `io.Reader`: the byte source contract

```go
type Reader interface {
	Read(p []byte) (n int, err error)
}
```

`Read` fills **up to** `len(p)` bytes of `p` and reports how many it wrote.

| Rule | Meaning for you |
|---|---|
| `0 <= n <= len(p)` | Only `p[:n]` is valid data. Never use all of `p`. |
| May return `n < len(p)` with `err == nil` | A **short read** is normal. Call `Read` again. |
| May return `n > 0` **and** a non-nil `err` | **Process the `n` bytes first**, then look at `err`. |
| End of input is `err == io.EOF` | `io.EOF` is not a failure; it means "no more data". |
| At the end, either `(n, io.EOF)` or `(n, nil)` then `(0, io.EOF)` | Your loop must handle both. |
| `(0, nil)` is discouraged but legal | It means "nothing happened", **not** EOF. |
| Implementations must not retain `p` | Don't keep a reference to the caller's buffer. |

The canonical read loop:

```go
r := strings.NewReader("hello, reader")
buf := make([]byte, 5)
for {
	n, err := r.Read(buf)
	fmt.Printf("n=%d err=%v %q\n", n, err, buf[:n]) // 1. use buf[:n] first
	if err == io.EOF {                              // 2. then check the error
		break
	}
	if err != nil {
		// a real failure: handle it (next lesson)
		break
	}
}
```

```text
n=5 err=<nil> "hello"
n=5 err=<nil> ", rea"
n=3 err=<nil> "der"
n=0 err=EOF ""
```

- You will usually let a helper (`io.ReadAll`, `io.Copy`, `bufio.Scanner`) run this loop for you.
- When you **implement** `Read`, you must follow the same contract.

## 4. `io.Writer`: the byte sink contract

```go
type Writer interface {
	Write(p []byte) (n int, err error)
}
```

| Rule | Meaning |
|---|---|
| Writes all of `p` or fails | If `n < len(p)`, `err` **must** be non-nil. |
| Must not modify `p` | Not even temporarily. |
| Must not retain `p` | Copy the bytes if you need them later. |

Because a successful `Write` writes everything, callers rarely loop: they just check `err`.

## 5. Concrete readers and writers

| Type | Reader? | Writer? | Notes |
|---|---|---|---|
| `*strings.Reader` (`strings.NewReader(s)`) | ✅ | ❌ | Read a string as a stream |
| `*bytes.Buffer` (zero value ready) | ✅ | ✅ | Write, then read back; `.String()` shows unread bytes |
| `*strings.Builder` (zero value ready) | ❌ | ✅ | Most efficient way to build a string; don't copy a non-zero one |
| `*bytes.Reader` (`bytes.NewReader(b)`) | ✅ | ❌ | Like `strings.Reader` for `[]byte` |
| `os.Stdout`, `os.Stderr` | — | ✅ | The terminal (files come in Module 6) |
| `io.Discard` | — | ✅ | Throws everything away (like `/dev/null`) |

⚠️ The methods have **pointer receivers**: pass `&buf`, not `buf`.

```go
var b bytes.Buffer
b.WriteString("hi ")
fmt.Fprintf(&b, "%d+%d=%d", 2, 2, 4) // &b is an io.Writer
fmt.Println(b.String())             // hi 2+2=4

var sb strings.Builder
fmt.Fprintf(&sb, "x=%d", 7)
fmt.Println(sb.String(), sb.Len()) // x=7 3
```

## 6. `fmt.Fprint*`: format into any Writer

`Printf` is just `Fprintf(os.Stdout, ...)`. Write functions against `io.Writer` and the caller decides where output goes.

```go
func WriteGreeting(w io.Writer, name string) error {
	_, err := fmt.Fprintf(w, "Hello, %s!\n", name)
	return err
}

WriteGreeting(os.Stdout, "Ann")    // to the terminal
var b bytes.Buffer
WriteGreeting(&b, "Bob")           // to memory: perfect for tests
```

- `fmt.Fprint`, `fmt.Fprintln`, `fmt.Fprintf` all take a `Writer` first and return `(n int, err error)`.
- `io.WriteString(w, s)` writes a string without formatting (uses `w.WriteString` if `w` has one).

## 7. `io` helpers

```go
data, err := io.ReadAll(strings.NewReader("all of it")) // "all of it" <nil>

n, err := io.Copy(os.Stdout, strings.NewReader("copied\n")) // prints copied; n=7 err=<nil>
```

| Helper | Signature | Does |
|---|---|---|
| `io.ReadAll` | `(r Reader) ([]byte, error)` | Read until EOF; success is `err == nil` (not `io.EOF`) |
| `io.Copy` | `(dst Writer, src Reader) (int64, error)` | Stream src → dst until EOF; success is `err == nil` |
| `io.CopyN` | `(dst Writer, src Reader, n int64) (int64, error)` | Copy exactly n bytes (error if fewer) |
| `io.WriteString` | `(w Writer, s string) (int, error)` | Write a string |
| `io.MultiWriter` | `(ws ...Writer) Writer` | Every write goes to **all** writers (like Unix `tee`) |
| `io.MultiReader` | `(rs ...Reader) Reader` | Read rs one after another, as one stream |
| `io.LimitReader` | `(r Reader, n int64) Reader` | Reports EOF after n bytes |
| `io.TeeReader` | `(r Reader, w Writer) Reader` | Everything read from r is also written to w |

```go
var a, b bytes.Buffer
io.WriteString(io.MultiWriter(&a, &b), "twice") // a and b both hold "twice"

mr := io.MultiReader(strings.NewReader("ab"), strings.NewReader("cd"))
all, _ := io.ReadAll(mr) // "abcd"

lr := io.LimitReader(strings.NewReader("0123456789"), 4)
all, _ = io.ReadAll(lr) // "0123"

var log bytes.Buffer
tr := io.TeeReader(strings.NewReader("seen"), &log)
all, _ = io.ReadAll(tr) // all == "seen", log holds "seen" too
```

## 8. `bufio`: scanning and buffered writing

### `bufio.Scanner`: read lines or words

```go
sc := bufio.NewScanner(strings.NewReader("one two\nthree\n\nfour"))
for sc.Scan() {             // false at EOF or on error
	fmt.Printf("%q\n", sc.Text()) // the line, without "\n"
}
if err := sc.Err(); err != nil { // nil at a clean EOF
	// handle the error
}
```

```text
"one two"
"three"
""
"four"
```

- Lines: the default. `\n` (and a `\r` before it) is stripped; a final line without `\n` still counts; empty lines are returned as `""`.
- Words: call `sc.Split(bufio.ScanWords)` **before** the first `Scan`. Words are split on any Unicode space; never empty.
- Other splitters: `bufio.ScanRunes`, `bufio.ScanBytes`.
- `sc.Text()` returns a string; `sc.Bytes()` returns a slice that the **next `Scan` overwrites**.
- Lines longer than `bufio.MaxScanTokenSize` (64 KiB) stop the scan with an error unless you call `sc.Buffer`.

### `bufio.Writer`: batch small writes — then **Flush**

```go
var out bytes.Buffer
bw := bufio.NewWriter(&out)
fmt.Fprint(bw, "buffered")
fmt.Printf("%q\n", out.String()) // ""  ← still in bw's memory buffer
bw.Flush()
fmt.Printf("%q\n", out.String()) // "buffered"
```

- `bufio.Writer` collects writes in memory (4096 bytes by default) and passes them on in large chunks.
- Forgetting `Flush` silently **loses** the tail of your output. Always `return bw.Flush()` (it reports any write error).

## 9. Decorating: wrap a Reader or Writer

Because callers only see the interface, you can wrap one Reader/Writer in another that adds behaviour. This is how `LimitReader`, `TeeReader` and `bufio` work.

```go
// UpperWriter upper-cases ASCII letters before passing bytes on.
type UpperWriter struct {
	W io.Writer
}

func (u UpperWriter) Write(p []byte) (int, error) {
	q := make([]byte, len(p)) // a copy: never modify p
	for i, c := range p {
		if 'a' <= c && c <= 'z' {
			c -= 'a' - 'A'
		}
		q[i] = c
	}
	if _, err := u.W.Write(q); err != nil {
		return 0, err
	}
	return len(p), nil // success: report what the caller gave us
}

fmt.Fprintln(UpperWriter{W: os.Stdout}, "shout") // SHOUT
```

A **reader** decorator reads from the inner reader into the caller's buffer, transforms `p[:n]`, and returns the inner `n, err` unchanged:

```go
type rot13Reader struct{ r io.Reader }

func (rr rot13Reader) Read(p []byte) (int, error) {
	n, err := rr.r.Read(p)
	for i := range p[:n] { // only the n valid bytes
		p[i] = rot13(p[i])
	}
	return n, err // pass both through: n>0 with err is legal
}
```

- Decorators stack: `bufio.NewScanner(io.LimitReader(UpperReader{r}, 1024))`.
- A writer decorator should return `len(p)` on success (the caller's view), even if it wrote more or fewer bytes underneath.

## 10. `sort.Interface` (legacy) vs `slices.SortFunc`

```go
type Person struct {
	Name string
	Age  int
}

type ByAge []Person

func (a ByAge) Len() int           { return len(a) }
func (a ByAge) Less(i, j int) bool { return a[i].Age < a[j].Age }
func (a ByAge) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }

people := []Person{{"Ann", 31}, {"Bob", 25}, {"Cy", 40}}
sort.Sort(ByAge(people))               // [{Bob 25} {Ann 31} {Cy 40}]
sort.Sort(sort.Reverse(ByAge(people))) // [{Cy 40} {Ann 31} {Bob 25}]
```

- `ByAge(people)` is a **conversion**, not a copy of the elements: it shares the backing array, so `people` gets sorted.
- `sort.Sort` is not stable; `sort.Stable` keeps equal elements in their original order.
- `sort.Reverse(data)` wraps an `Interface` and flips `Less` — another decorator.
- **New code:** prefer `slices.SortFunc(people, func(a, b Person) int { return cmp.Compare(a.Age, b.Age) })`. The `sort` docs say it is "more ergonomic and runs faster". You will still meet `sort.Interface` in older code.
- `cmp.Ordered` is a **type constraint** (not a method interface) listing types that support `<`; it powers `slices.Sort`, `max`, `cmp.Compare`. Generics come in a later module.

## 11. Design lessons from the standard library

- **Accept interfaces, return concrete types.** `func Count(r io.Reader)` works with strings, buffers, files and network connections.
- **Keep interfaces small.** One method is easy to implement and easy to fake in tests.
- **Test with in-memory implementations:** `strings.NewReader` as input, `bytes.Buffer` / `strings.Builder` as output.

## Common mistakes

| Mistake | Fix |
|---|---|
| Using `buf` instead of `buf[:n]` after `Read` | Only the first `n` bytes are valid |
| Assuming one `Read` fills the buffer | Loop until `io.EOF`, or use `io.ReadAll` / `io.ReadFull` |
| Checking `err` before using `n` | Process `p[:n]` first: `Read` can return data **and** `io.EOF` |
| Treating `io.EOF` from `ReadAll` / `Copy` as possible | They return `nil` on success; only a real failure is non-nil |
| Forgetting `bw.Flush()` | End the function with `return bw.Flush()` |
| `String()` on `*T` but printing a `T` value | Use a value receiver for `String` |
| `fmt.Sprintf("%v", t)` inside `t.String()` | Format the underlying value: `int(t)`, `%d`, or fields |
| Passing `buf` (a `bytes.Buffer` value) where an `io.Writer` is expected | Pass `&buf`: the methods have pointer receivers |
| A writer decorator returning the inner `n` after transforming | Return `len(p)` on success, as the caller sees it |
| Modifying `p` inside `Write` | Transform a copy (`make([]byte, len(p))`) |
| `sc.Split(...)` after the first `Scan` | Set the split function before scanning (it panics otherwise) |
| Writing new code with `sort.Interface` | `slices.SortFunc` + `cmp.Compare` |

## Key takeaways

- One-method interfaces (`Stringer`, `error`, `Reader`, `Writer`) are what make the standard library composable.
- `String() string` controls `%v`/`%s`/`Println`; prefer a value receiver and never `%v` yourself inside it.
- `Read` may return fewer bytes than asked, and may return data together with `io.EOF`: **use `p[:n]` first, then check `err`**.
- `Write` writes everything or returns an error.
- `strings.Reader`, `bytes.Buffer`, `strings.Builder` are the in-memory readers/writers for tests; `fmt.Fprintf` writes to any of them.
- `io.ReadAll`, `io.Copy`, `io.MultiWriter`, `io.MultiReader`, `io.LimitReader`, `io.TeeReader` compose streams; `bufio.Scanner` splits lines/words; `bufio.Writer` needs `Flush`.
- Wrap readers/writers to add behaviour (decorators); callers can't tell the difference.
- `sort.Interface` is legacy; use `slices.SortFunc`.

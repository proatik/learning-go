# Structs

## 1. Declaring a struct type

A **struct** is a sequence of named **fields**, each with a type. You almost always give it a name with `type`.

```go
type Point struct {
	X int
	Y int
}

type Employee struct {
	Name      string
	Dept      string
	Salary    float64
	Skills    []string          // any type may be a field: slices, maps, other structs…
	Manager   string
	StartYear int
}

type Point3 struct {
	X, Y, Z float64 // consecutive fields of one type can be grouped
}
```

- Field names must be unique within the struct; the blank field `_` is allowed (padding) but rare.
- Field order is part of the type: `struct{ X, Y int }` and `struct{ Y, X int }` are different types.
- Declare types at package level so all functions can use them (a `type` inside a function is allowed, but only visible there).

## 2. Composite literals: keyed vs positional

```go
p1 := Point{X: 1, Y: 2} // keyed: field: value
p2 := Point{Y: 2}       // keyed: omitted fields get their zero value → {0 2}
p3 := Point{1, 2}       // positional: every field, in declaration order
p4 := Point{}           // all zero → {0 0}
```

| | Keyed `Point{X: 1, Y: 2}` | Positional `Point{1, 2}` |
|---|---|---|
| Omitting fields | ✅ omitted = zero value | ❌ `too few values in struct literal` |
| Order | any | must match the declaration |
| Survives adding / reordering fields | ✅ | ❌ breaks or, worse, silently swaps same-typed fields |
| Mixing styles | ❌ `mixture of field:value and value elements` | |

**Prefer keyed literals.** `go vet` warns about positional (unkeyed) literals of struct types **from other packages**. Positional is fine for tiny, stable types like `Point{1, 2}` in the same package.

Multi-line literals need a **trailing comma** after the last element (the lexer inserts a `;` at line end otherwise):

```go
e := Employee{
	Name:   "Ann",
	Dept:   "eng",
	Salary: 95000,
} // ← the comma after 95000 is required
```

## 3. Zero value and field access

The zero value of a struct has **every field set to its own zero value**, so a struct is usable without a constructor.

```go
type Config struct {
	Host string
	Port int
}

var c Config            // Config{Host: "", Port: 0}
fmt.Printf("%+v\n", c)  // {Host: Port:0}

p := Point{1, 2}
p.X = 10                // read and write fields with .
fmt.Println(p.X + p.Y)  // 12
```

Design types so their zero value is useful (Effective Go): an empty `Name` string, a `0` count and a `nil` slice you can `append` to are all ready to use.

Constructors are just ordinary functions, conventionally named `NewT`, used when fields need validation or non-zero defaults:

```go
func NewConfig(host string) Config {
	return Config{Host: host, Port: 8080}
}
```

## 4. Exported vs unexported fields

The capitalisation rule from L1 applies to **each field**.

```go
package store

type Item struct {
	Name  string // exported: other packages can read/write item.Name
	Price int    // exported
	cost  int    // unexported: only code in package store sees it
}
```

- Other packages can use `store.Item{Name: "pen"}` but **cannot** set or read `cost`, and positional literals of `Item` are impossible outside `store` (they'd need every field).
- `encoding/json`, `fmt` with `%v`, etc. all *see* that `cost` exists, but `json` only encodes **exported** fields (§14).

## 5. Structs are values

Assigning a struct, passing it to a function or returning it **copies every field**.

```go
a := Point{1, 2}
b := a        // b is an independent copy
b.X = 100
fmt.Println(a, b) // {1 2} {100 2}

func moveRight(p Point) Point { // p is a copy of the argument
	p.X++
	return p // return the new value; the caller decides what to do with it
}

a = moveRight(a) // the "copy, modify, return" style used throughout this lesson
```

⚠️ The copy is **shallow**: a slice or map field is a small header that refers to shared data, so both copies see changes to the *elements*.

```go
type Player struct {
	Name  string
	Score int
	Tags  []string
}

p := Player{Name: "ann", Tags: []string{"x"}}
q := p
q.Score = 5          // only q changes
q.Tags[0] = "CHANGED" // shared backing array: p sees this too
fmt.Printf("%+v\n%+v\n", p, q)
// {Name:ann Score:0 Tags:[CHANGED]}
// {Name:ann Score:5 Tags:[CHANGED]}
```

To let a function change the caller's struct you need a **pointer** — next lesson. Preview: `p := &s; p.X = 1` changes `s` itself (no `->`, the `.` works through the pointer).

## 6. Comparison and structs as map keys

Structs are **comparable** with `==` / `!=` when **all their fields are comparable**. Two struct values are equal if all corresponding fields are equal.

```go
Point{1, 2} == Point{X: 1, Y: 2} // true
p == Point{}                      // the idiomatic "is it the zero value?" test

type Bag struct{ Items []string }
// Bag{} == Bag{} → compile error: struct containing []string cannot be compared
```

| Field types | `==` | Map key |
|---|---|---|
| numbers, strings, bools, arrays of comparable, comparable structs, pointers | ✅ | ✅ |
| any slice, map or function field | ❌ compile error | ❌ |

A comparable struct is the natural **composite map key**, much cleaner than building `"x,y"` strings:

```go
visited := map[Point]bool{}
visited[Point{1, 2}] = true
fmt.Println(visited[Point{X: 1, Y: 2}], visited[Point{2, 1}]) // true false

type edge struct{ From, To string }
weights := map[edge]int{{"a", "b"}: 3} // element type elided in map literals (§8)
```

For structs that contain slices, compare in tests with `reflect.DeepEqual` or field by field (`slices.Equal` for the slice).

## 7. Anonymous structs

A struct type can be used **without a name**, directly in a literal. Ideal for one-off data.

```go
cfg := struct {
	Host string
	Port int
}{"localhost", 8080}
fmt.Println(cfg.Host) // localhost
```

The most common use, which you see in every test file of this course: **table-driven tests**.

```go
tests := []struct {
	name string
	a, b Point
	want Point
}{
	{"origin", Point{}, Point{}, Point{}},
	{"simple", Point{1, 2}, Point{3, 4}, Point{4, 6}},
}
for _, tt := range tests {
	if got := add(tt.a, tt.b); got != tt.want {
		t.Errorf("%s: add(%+v, %+v) = %+v, want %+v", tt.name, tt.a, tt.b, got, tt.want)
	}
}
```

An anonymous struct value is assignable to a named type with **identical fields**: `var pt Point = struct{ X, Y int }{3, 4}` compiles.

## 8. Slices of structs

```go
staff := []Employee{
	{Name: "Ann", Dept: "eng", Salary: 95000}, // element type Employee elided
	{Name: "Bob", Dept: "ops", Salary: 70000},
}
staff = append(staff, Employee{Name: "Cy", Dept: "eng"})
```

Inside slice, array and map literals, the element type may be **elided**: `{Name: "Ann"}` instead of `Employee{Name: "Ann"}`.

⚠️ The `range` value variable is a **copy** of the element:

```go
pts := []Point{{1, 1}, {2, 2}}
for _, pt := range pts {
	pt.X = 0 // modifies the copy: pts unchanged → [{1 1} {2 2}]
}
for i := range pts {
	pts[i].X = 0 // modifies the element itself → [{0 1} {0 2}]
}
```

Finding an element: `i := slices.IndexFunc(staff, func(e Employee) bool { return e.Name == "Bob" })` (-1 if none).

Functions that "update" a slice of structs and promise **value semantics** must not touch the caller's elements: build the result in a new slice (`slices.Clone(s)` first, or `append` into a fresh slice).

## 9. Maps of structs: copy, modify, store

A map element is **not addressable**, so you can't assign to one of its fields:

```go
counts := map[string]Point{"a": {1, 2}}
counts["a"].X = 10 // compile error: cannot assign to struct field counts["a"].X in map
```

Read a copy, modify it, store it back:

```go
p := counts["a"] // a missing key gives the zero struct, so this also works for new keys
p.X = 10
counts["a"] = p
fmt.Println(counts) // map[a:{10 2}]
```

Typical aggregation:

```go
type DeptStats struct {
	Count int
	Total float64
}

stats := map[string]DeptStats{}
for _, e := range staff {
	s := stats[e.Dept]
	s.Count++
	s.Total += e.Salary
	stats[e.Dept] = s
}
```

(A map of *pointers* to structs, `map[string]*DeptStats`, allows `stats[k].Count++` — next lesson.)

## 10. Sorting slices of structs

`slices.SortFunc` (L2) with a comparison function returning negative / 0 / positive. Combine several keys with `cmp.Or`, which returns its first non-zero argument:

```go
type Emp struct {
	Name string
	Dept string
	Age  int
}

emps := []Emp{{"Cy", "ops", 30}, {"Al", "dev", 41}, {"Bo", "dev", 25}, {"Di", "ops", 30}}
slices.SortFunc(emps, func(a, b Emp) int {
	return cmp.Or(
		cmp.Compare(a.Dept, b.Dept), // 1st key: dept ascending
		cmp.Compare(b.Age, a.Age),   // 2nd key: age DEScending (b, a swapped)
		cmp.Compare(a.Name, b.Name), // 3rd key: name ascending (tie-breaker)
	)
})
fmt.Println(emps) // [{Al dev 41} {Bo dev 25} {Cy ops 30} {Di ops 30}]
```

- **Descending** = swap the arguments: `cmp.Compare(b.X, a.X)`.
- `SortFunc` is not stable; use `slices.SortStableFunc` when equal elements must keep their input order, or add a final tie-breaker key.
- Sorting reorders the caller's slice in place. To keep the input intact, sort `slices.Clone(s)`.

## 11. Nested structs

A field can itself be a struct; chain selectors to reach inside.

```go
type Rect struct {
	Min, Max Point
}

r := Rect{Min: Point{0, 0}, Max: Point{4, 3}}
width := r.Max.X - r.Min.X // 4
r.Max.Y = 10               // assign through a chain of fields: fine for a variable
```

Nested literals can't elide the inner type in struct literals (only in slice/map/array literals): `Rect{Min: {0, 0}}` doesn't compile; write `Point{0, 0}`.

## 12. Embedding: promoted fields

A field declared with **a type but no name** is an **embedded field**. Its field name is the type name, and its fields are **promoted**: reachable as if declared in the outer struct.

```go
type Address struct {
	Street, City string
}

type Customer struct {
	Name    string
	Address // embedded: field name is "Address"
}

var c Customer
c.City = "Oslo"         // promoted: shorthand for…
c.Address.City = "Oslo" // …the full path; both work
```

**Go 1.27:** a keyed literal may use promoted field names directly (a key may be any valid field selector of the struct type):

```go
c := Customer{Name: "Ann", City: "Oslo", Street: "Main 1"} // Go 1.27+
fmt.Printf("%+v\n", c) // {Name:Ann Address:{Street:Main 1 City:Oslo}}

c = Customer{Name: "Ann", Address: Address{City: "Oslo"}} // the long form still works (and is required before 1.27)
```

- You may not mix both for the same embedded field: `Customer{Address: Address{}, Street: "x"}` → `cannot specify promoted field Street and enclosing embedded field Address`.
- With `go 1.26` or older in `go.mod`, the short form fails with `... requires go1.27 or later`.

Name resolution rules:

| Situation | Result |
|---|---|
| Outer struct has its own `Name` and the embedded one has `Name` too | The **shallower** one wins (`c.Name` is the outer); reach the other with `c.Address.Name` |
| Two embedded types at the same depth both have `ID` | `d.ID` → compile error `ambiguous selector d.ID`; write `d.Base.ID` / `d.Audit.ID` |

Embedding is **composition, not inheritance**: a `Customer` is not an `Address`; you can't pass a `Customer` where an `Address` is expected (pass `c.Address`). Embedding also promotes *methods* — that's Module 3.

## 13. Printing: `%v`, `%+v`, `%#v`

```go
type Config struct {
	Host string
	Port int
}
cfg := Config{"localhost", 8080}
fmt.Printf("%v\n", cfg)  // {localhost 8080}
fmt.Printf("%+v\n", cfg) // {Host:localhost Port:8080}
fmt.Printf("%#v\n", cfg) // main.Config{Host:"localhost", Port:8080}
```

| Verb | Shows | Use for |
|---|---|---|
| `%v` | values only | compact output, `fmt.Println` |
| `%+v` | `Field:value` | **debugging and test messages** (this course's default) |
| `%#v` | Go syntax: type, field names, quoted strings | copy-pasteable values, spotting `""` vs `" "` |

## 14. Struct tags and JSON (preview)

A **tag** is a string literal after a field's type. The language ignores it; packages read it with reflection. `encoding/json` uses the `json` key:

```go
type Product struct {
	ID       int      `json:"id"`             // JSON key "id"
	Name     string   `json:"name"`
	Price    float64  `json:"price"`
	Tags     []string `json:"tags,omitempty"` // omit when empty (len 0 / nil)
	Internal string   `json:"-"`              // never encoded
	secret   string   // unexported: never encoded
}

b, err := json.Marshal(Product{ID: 1, Name: "Pen", Price: 1.5, Internal: "x"})
fmt.Println(string(b), err) // {"id":1,"name":"Pen","price":1.5} <nil>
```

| Tag | Effect |
|---|---|
| none | key = Go field name (`"Name"`) |
| `json:"name"` | key = `"name"` |
| `json:"name,omitempty"` | skip if `false`, `0`, `""`, nil pointer, or empty slice/map/array |
| `json:",omitempty"` | default key, with omitempty |
| `json:"name,omitzero"` | skip if the zero value (works for structs too, e.g. `time.Time`) |
| `json:"-"` | always skip |

- Keys appear in **field declaration order**. Fields of an **embedded** struct without a tag are flattened into the outer object.
- The tag syntax is `` `key:"value"` `` — no spaces around `:`, double quotes required. `go vet` checks malformed tags.
- Full JSON coverage (decoding, `json.Unmarshal`, errors) comes in a later module.

## 15. The empty struct `struct{}`

`struct{}` has no fields and occupies **zero bytes**. Its only value is `struct{}{}`.

```go
seen := map[string]struct{}{} // a set: only the keys matter
seen["go"] = struct{}{}
_, ok := seen["go"] // true
```

`map[string]bool` is also fine and easier to read (`if seen[k]`); `struct{}` signals "no value here" and saves memory in big sets. You will meet `chan struct{}` (a pure signal) in the concurrency module.

## Common mistakes

| Mistake | Fix |
|---|---|
| Positional literal `Employee{"Ann", "eng", 95000, nil, "", 0}` | Keyed: `Employee{Name: "Ann", Dept: "eng", Salary: 95000}` |
| Missing trailing comma in a multi-line literal | Every line ends with `,` including the last |
| Expecting a function to change the caller's struct | It gets a copy: return the new struct (`p = move(p)`), or use a pointer (next lesson) |
| `for _, s := range xs { s.X = 1 }` | Modify through the index: `xs[i].X = 1` |
| `m[k].Field = v` | `s := m[k]; s.Field = v; m[k] = s` |
| Assuming a struct copy deep-copies its slice/map fields | Copy them explicitly: `q.Tags = slices.Clone(p.Tags)` |
| `==` on a struct with a slice/map field | Compare fields, `slices.Equal`, or `reflect.DeepEqual` in tests |
| Building string keys `fmt.Sprint(x, ",", y)` | Use a struct key: `map[Point]int` |
| Lower-case field that should be encoded to JSON or used by another package | Capitalize it; set the external name with a tag |
| `` `json: "name"` `` or `` `json:name` `` | `` `json:"name"` `` (`go vet` flags it) |
| Descending sort written as `cmp.Compare(a.X, b.X)` | Swap: `cmp.Compare(b.X, a.X)` |
| Ambiguous promoted field from two embedded types | Use the full path `d.Base.ID` |

## Key takeaways

- `type T struct { Field Type; ... }`; capitalized fields are exported.
- Prefer **keyed** literals; omitted fields are zero. The zero struct is ready to use.
- Structs are **values**: assignment, arguments and results copy all fields (shallowly). Write "copy, modify, return" functions; pointers come next.
- Structs with only comparable fields work with `==` and as **map keys**.
- Anonymous structs power table-driven tests and one-off data.
- `range` gives copies; map elements aren't addressable: copy, modify, store.
- Sort with `slices.SortFunc` + `cmp.Or(cmp.Compare(...), ...)`; swap arguments for descending.
- Embedding promotes fields (composition, not inheritance); Go 1.27 lets literals use promoted keys.
- `%+v` for debugging; tags like `` `json:"name,omitempty"` `` control encoding; `struct{}` costs nothing.

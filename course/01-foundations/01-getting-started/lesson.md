# Getting Started: Toolchain & Program Structure

## 1. The toolchain

Go ships as **one binary, `go`**, that compiles, tests, formats, vets and manages dependencies.

```sh
go version          # go version go1.27.3 windows/amd64
go env GOPATH       # where downloaded modules live
```

Go releases a new minor version every 6 months (Feb & Aug). Each `go.mod` declares the minimum version it needs; with the default `GOTOOLCHAIN=auto`, an older local `go` **automatically downloads** the required toolchain. That is how this repo runs on Go 1.27 even if you installed something older.

## 2. Modules and packages

| Concept | What it is | Defined by |
|---|---|---|
| **Module** | A versioned collection of packages; the unit of dependency management | `go.mod` at the module root |
| **Package** | A directory of `.go` files compiled together | `package <name>` at the top of each file |
| **Import path** | Module path + directory, e.g. `learngo/course/...` | location on disk |

```sh
mkdir hello && cd hello
go mod init example.com/hello     # creates go.mod
```

```text
// go.mod
module example.com/hello

go 1.27.0
```

Rules:
- All `.go` files in one directory belong to **one** package (test files may also use `<name>_test`).
- **`package main` + `func main()`** = an executable. Any other package name = a library.
- Files and directories starting with `_` or `.` are ignored by the go tool (this course hides solutions in `_solution/`).

## 3. Anatomy of a program

```go
package main // 1. package clause — always first

import ( // 2. imports — grouped in parentheses
	"fmt"
	"math"
)

// 3. top-level declarations: const, var, type, func
const greeting = "Hello"

func main() { // entry point: no args, no return value
	fmt.Println(greeting, "Go!")
	fmt.Printf("√2 ≈ %.3f\n", math.Sqrt(2))
}
```

```sh
go run .            # compile + run the package in the current directory
go build            # produce an executable (hello.exe on Windows)
```

The compiler is strict on purpose:
- An **unused import** is a compile error.
- An **unused local variable** is a compile error.
- No semicolons needed — the lexer inserts them at line ends. Consequence: the opening `{` **must** be on the same line (`func main() {`).

## 4. Exported vs unexported

Visibility is decided by the **first letter** of a name — no `public`/`private` keywords.

```go
package shapes

func Area(w, h float64) float64 { return w * h } // Exported: shapes.Area works from other packages
func perimeter(w, h float64) float64 { ... }     // unexported: only inside package shapes
```

You can only use exported names from imported packages: `math.Pi` ✅, `math.pi` ❌.

## 5. Everyday commands

| Command | Purpose |
|---|---|
| `go run .` | Build and run the current package |
| `go build` | Compile; produces a binary for `main` packages |
| `go test ./...` | Run tests in this and all sub-packages |
| `go fmt ./...` | Format code (there is one official style — never argue about formatting) |
| `go vet ./...` | Static checks for likely bugs (e.g. wrong `Printf` verbs) |
| `go doc fmt.Println` | Show documentation from the terminal |
| `go mod tidy` | Add missing / remove unused dependencies in `go.mod` |
| `go fix ./...` | Apply *modernizers* — rewrite code to use newer language features & APIs |

> VS Code with the official **Go extension** runs `gofmt` on save and shows `vet` / compile errors as you type.

## 6. Comments and doc comments

```go
// Greet returns a greeting for name.     ← doc comment: starts with the name, full sentence
func Greet(name string) string { ... }

/* block comments exist but line comments are idiomatic */
```

A comment directly above a top-level declaration is its **doc comment**; `go doc` and pkg.go.dev display it.

## 7. How the exercises are checked (reading tests)

Tests live in files ending in `_test.go`, in the same package. You don't write tests yet, but you should be able to read them — they are the spec.

```go
package hello

import "testing"

func TestHello(t *testing.T) {          // name starts with Test, takes *testing.T
	got := Hello()
	want := "Hello, Go!"
	if got != want {
		t.Errorf("Hello() = %q, want %q", got, want)  // marks the test as failed
	}
}
```

**Example tests** check what a function prints to stdout:

```go
func ExamplePrintBanner() {
	PrintBanner()
	// Output:
	// === Welcome ===
}
```

`go test` compiles the package with its tests, runs every `TestXxx` and `ExampleXxx`, and also runs a subset of **`go vet`** — a vet failure fails the test run.

```sh
go test -v ./course/01-foundations/01-getting-started/exercises/01-first-program/
```

## Key takeaways

- One `go` command does everything; `go.mod` pins the Go version.
- A directory = a package; `package main` + `func main()` = a program.
- Capitalized = exported. Unused imports/variables don't compile.
- `gofmt` style is not optional; `go vet` catches bugs early.

# Control Flow: if, for, switch

Go has exactly **three** control-flow keywords for branching and looping: `if`, `for` and `switch`. There is no `while`, no `do`, no ternary `?:`.

## 1. `if` / `else if` / `else`

```go
if x > 0 {
	fmt.Println("positive")
} else if x < 0 {
	fmt.Println("negative")
} else {
	fmt.Println("zero")
}
```

- **No parentheses** around the condition (`gofmt` removes them).
- **Braces are mandatory**, even for one statement. `{` must be on the same line as `if`; `else` must be on the same line as the closing `}`.
- The condition must be a `bool`. `if n { ... }` with an `int` does not compile — write `if n != 0`.

Idiomatic Go **returns early** instead of nesting:

```go
func Describe(age int) string {
	if age < 0 {
		return "invalid"
	}
	if age < 18 {
		return "minor"
	}
	return "adult" // no else needed: the earlier branches returned
}
```

## 2. `if` with a short statement

An `if` may start with a short statement, separated by `;`. Variables declared there are scoped to the `if` **and all its `else` branches**.

```go
if n, err := strconv.Atoi(s); err != nil {
	fmt.Println("bad number:", err)
} else {
	fmt.Println("double:", n*2) // n is visible here too
}
// n and err do NOT exist here
```

```go
if r := x % 2; r == 0 {
	fmt.Println("even")
}
```

Use it to keep short-lived variables out of the enclosing scope.

## 3. `for` — the only loop

| Form | Syntax | Like |
|---|---|---|
| Three-part | `for i := 0; i < n; i++ { }` | C `for` |
| Condition only | `for n > 1 { }` | `while` |
| Infinite | `for { }` | `while (true)` — leave with `break` or `return` |
| Range over int (Go 1.22+) | `for i := range n { }` | `i` = 0, 1, …, n-1 |
| Range, no variable | `for range n { }` | repeat n times |

```go
for i := 0; i < 3; i++ { // init; condition; post
	fmt.Print(i, " ") // 0 1 2
}

n := 27
steps := 0
for n != 1 { // "while" loop
	if n%2 == 0 {
		n /= 2
	} else {
		n = 3*n + 1
	}
	steps++
}

for i := range 5 { // 0 1 2 3 4 — preferred when counting from 0
	fmt.Print(i)
}

for range 3 {
	fmt.Println("hi") // printed 3 times
}
```

- `i++` and `i--` are **statements**, not expressions: `x = i++` does not compile.
- `for i := range n` with `n <= 0` runs zero times.
- To count from 1 to n inclusive: `for i := 1; i <= n; i++` (or `for i := range n` and use `i+1`).

### Loop variables are per-iteration (Go 1.22+)

Each iteration gets a **fresh copy** of variables declared by `for` (`i` above). Before Go 1.22 the variable was shared across iterations, which caused famous bugs when closures or goroutines captured it. With `go 1.22` or later in `go.mod`, that bug is gone. You will see why this matters in the Functions and Concurrency lessons.

## 4. `break` and `continue`

| Statement | Effect |
|---|---|
| `break` | exit the innermost `for` (or `switch`) immediately |
| `continue` | skip the rest of this iteration; go to the next one |

```go
for i := range 10 {
	if i%2 == 0 {
		continue // skip even numbers
	}
	if i > 7 {
		break // stop entirely
	}
	fmt.Print(i, " ") // 1 3 5 7
}
```

### Labeled `break` / `continue`

A plain `break` only leaves the **innermost** loop. To control an outer loop, put a **label** before it:

```go
outer:
	for i := 1; i <= 9; i++ {
		for j := 1; j <= 9; j++ {
			if i*j == 42 {
				fmt.Println(i, j) // 6 7
				break outer        // leaves BOTH loops
			}
		}
	}
```

```go
next:
	for n := 2; n <= 20; n++ {
		for d := 2; d*d <= n; d++ {
			if n%d == 0 {
				continue next // n is not prime: move on to the next n
			}
		}
		fmt.Print(n, " ") // 2 3 5 7 11 13 17 19
	}
```

- A label is an identifier followed by `:`; `gofmt` out-dents it one level.
- An **unused label is a compile error**, like an unused variable.

## 5. `switch`

### Value switch

```go
switch day {
case "Sat", "Sun": // several values in one case
	fmt.Println("weekend")
case "Mon", "Tue", "Wed", "Thu", "Fri":
	fmt.Println("weekday")
default: // optional; runs when nothing matched
	fmt.Println("unknown")
}
```

- Cases are tested **top to bottom**; the first match runs and the switch **ends**. There is **no implicit fall-through** and no `break` needed.
- Case values don't have to be constants — any expression of a comparable type works.
- Duplicate constant cases (`case 1, 1:`) are a compile error.

### Tagless switch — a cleaner if-else chain

`switch {` with no value is `switch true {`: the first case whose condition is true wins.

```go
switch {
case score >= 90:
	grade = "A"
case score >= 80: // only reached if score < 90
	grade = "B"
default:
	grade = "F"
}
```

### Short statement in a switch

```go
switch r := n % 3; r {
case 0:
	fmt.Println("divisible by 3")
default:
	fmt.Println("remainder", r)
}
```

Like `if`, `r` is scoped to the switch.

### `fallthrough` (rare)

`fallthrough` must be the **last statement** of a case; it jumps into the **next case's body without testing its condition**.

```go
switch n {
case 0:
	fmt.Println("zero")
	fallthrough
case 1:
	fmt.Println("zero or one") // runs for 0 AND 1
}
```

Prefer `case 0, 1:` — real code almost never needs `fallthrough`.

### `break` inside a `switch`

`break` inside a `switch` exits the **switch**, not a surrounding `for`. Use a label to break out of the loop:

```go
loop:
	for {
		switch cmd := next(); cmd {
		case "quit":
			break loop // a plain break would only leave the switch
		}
	}
```

## 6. `min` and `max` builtins (Go 1.21+)

```go
min(3, 1, 2)       // 1
max(2.5, 1)        // 2.5
max(lo, min(x, hi)) // clamp x into [lo, hi]
```

They accept **one or more** arguments of the same ordered type (integers, floats, strings) — no import needed.

## Common mistakes

| Mistake | Fix |
|---|---|
| `if (x > 0) {` | Works, but `gofmt` drops the parens — write `if x > 0 {` |
| `if x > 0 return x` | Braces are always required |
| `}` newline `else {` | `else` must be on the same line: `} else {` |
| `if n {` where `n` is an `int` | Conditions must be `bool`: `if n != 0 {` |
| Using a short-statement variable after the `if` | It's out of scope — declare it before the `if` if you need it later |
| Writing `break` at the end of each `case` | Unnecessary — Go cases never fall through |
| Expecting `break` in a `switch` to exit the loop | It exits the switch; use a labeled `break` |
| `for i := 0; i < n; i++` vs `<= n` | Decide whether `n` is included; off-by-one errors are the classic loop bug |
| `for i := range n` expecting 1…n | It yields 0…n-1 |
| `while n > 0 {` | No `while` keyword — use `for n > 0 {` |
| Forgetting to change the condition variable in a `while`-style loop | Infinite loop — make sure every iteration makes progress |
| Calling a variable `min` or `max` | It shadows the builtin in that scope; pick another name |

## Key takeaways

- `if`: no parens, mandatory braces, `bool` conditions; prefer early `return` over deep nesting.
- `if v := f(); cond {}` / `switch v := f(); v {}` scope `v` to the statement.
- `for` is the only loop: three-part, condition-only, infinite, and `for i := range n`.
- `break`/`continue` act on the innermost loop; labels reach outer loops.
- `switch` cases don't fall through; list several values with commas; tagless `switch {}` replaces if-else chains.
- Loop variables are fresh per iteration since Go 1.22; `min`/`max` are builtins since Go 1.21.

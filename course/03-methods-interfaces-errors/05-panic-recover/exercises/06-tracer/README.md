# Call Tracer

Write the `defer t.Enter("name")()` tracing idiom: one line at the top of a function logs both its entry **and** its exit, with indentation showing nesting.

```go
type Tracer struct {
	W     io.Writer // where trace lines go; nil disables tracing
	depth int       // current nesting level
}

func (t *Tracer) Enter(name string) func()
```

Usage:

```go
func fact(t *Tracer, n int) int {
	defer t.Enter(fmt.Sprintf("fact(%d)", n))() // Enter runs now, the returned func at exit
	if n <= 1 {
		return 1
	}
	return n * fact(t, n-1)
}
```

`fact(&Tracer{W: os.Stdout}, 3)` prints:

```text
> fact(3)
| > fact(2)
| | > fact(1)
| | < fact(1)
| < fact(2)
< fact(3)
```

## Rules
- `Enter(name)` immediately writes `<indent>> <name>\n`, then increases the depth by one.
- The returned function decreases the depth by one, then writes `<indent>< <name>\n`.
- `<indent>` is `"| "` repeated `depth` times (the depth **before** entering / **after** leaving), so an enter line and its exit line have the same indent.
- Because it is deferred, the exit line must also be printed when the function **panics**, and the depth must be back to its old value afterwards.
- If `W` is `nil`, print nothing, but still return a **non-nil** function (a `defer` of a nil func would panic) and keep the depth unchanged.

# Let `go vet` Find the Bugs

`report.go` **compiles**, but it has bugs that `go vet` detects. `go test` runs vet first, so the tests won't even start until you fix them.

Required behaviour:

```go
Summary("apples", 3) // "apples: 3 items"
Percent(1, 4)        // "25.0%"
```

## Steps
1. Press **Run tests**.
2. Read each vet message: it names the file, the line and the problem.
3. Fix the format strings. Don't change the function signatures.

# Validation Error Type

Write a custom error type that tells the caller **which field** failed and **why**.

```go
type User struct {
	Name  string
	Email string
	Age   int
}

type ValidationError struct {
	Field  string
	Reason string
}

func (e *ValidationError) Error() string // "<field>: <reason>"

func Validate(u User) error
```

## Examples

| Call | Result |
|---|---|
| `(&ValidationError{"age", "must be at least 13"}).Error()` | `"age: must be at least 13"` |
| `Validate(User{"Ann", "ann@go.dev", 30})` | `nil` |
| `Validate(User{"", "ann@go.dev", 30})` | `&ValidationError{Field: "name", Reason: "must not be empty"}` |
| `Validate(User{"Ann", "ann.go.dev", 30})` | `&ValidationError{Field: "email", Reason: "must look like user@host"}` |
| `Validate(User{"Ann", "ann@go.dev", 9})` | `&ValidationError{Field: "age", Reason: "must be at least 13"}` |

## Rules
- Check the fields in this order and return the **first** failure as a `*ValidationError` (pointer):

| Field | Fails when | Reason |
|---|---|---|
| `name` | empty after `strings.TrimSpace` | `must not be empty` |
| `name` | `Name` has more than 20 **characters** (runes, not bytes) | `must be at most 20 characters` |
| `email` | not exactly one `@`, or nothing before it, or nothing after it | `must look like user@host` |
| `age` | less than 13 | `must be at least 13` |
| `age` | more than 150 | `must be at most 150` |

- On success, return a literal `nil`. The result type is `error`, not `*ValidationError`.
- Callers will use `errors.As` / `errors.AsType[*ValidationError]` to read `Field`.

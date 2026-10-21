# Errors Across Layers

Build a tiny 3-layer app: **repository → service → handler**. Lower layers add context by wrapping. The handler turns errors into HTTP-style status codes **once**, using `errors.Is` / `errors.AsType`.

```go
var ErrNotFound = errors.New("not found")

type InputError struct {
	Field string
	Msg   string
}

func (e *InputError) Error() string // "invalid <field>: <msg>"

type User struct {
	ID   int
	Name string
}

type Repo struct {
	Users map[int]User
	Fail  error // if non-nil, every Find fails with it (e.g. a database outage)
}

func (r *Repo) Find(id int) (User, error)

type Service struct {
	Repo *Repo
}

func (s *Service) GetUser(id int) (User, error)

func StatusOf(err error) int
func Handle(s *Service, id int) (status int, body string)
```

## Examples

With `Users = {7: {7, "Ann"}}`:

| Call | Result |
|---|---|
| `repo.Find(7)` | `User{7, "Ann"}, nil` |
| `repo.Find(9)` | error `find user 9: not found` |
| `svc.GetUser(9)` | error `get user: find user 9: not found` |
| `svc.GetUser(0)` | error `invalid id: must be positive` (a `*InputError`) |
| `Handle(svc, 7)` | `200, "Ann"` |
| `Handle(svc, 9)` | `404, "user not found"` |
| `Handle(svc, -3)` | `400, "invalid id: must be positive"` |
| `Handle(svc, 7)` with `Fail = errors.New("db down")` | `500, "internal error"` |

## Rules
- `Find`:
  - If `r.Fail != nil`, return `fmt.Errorf("find user %d: %w", id, r.Fail)`.
  - If the id is missing, return `fmt.Errorf("find user %d: %w", id, ErrNotFound)`.
- `GetUser`:
  - If `id <= 0`, return `&InputError{Field: "id", Msg: "must be positive"}` without calling the repo.
  - Otherwise wrap any repo error as `get user: <err>`.
  - Return a zero `User` with every error.
- `StatusOf`:
  - `nil` → `200`.
  - anything matching `ErrNotFound` → `404`.
  - anything containing a `*InputError` → `400`.
  - everything else → `500`.
  - Wrapped errors must be matched too.
- `Handle` calls `GetUser` and returns `StatusOf(err)` with this body:

| Status | Body |
|---|---|
| 200 | the user's `Name` |
| 404 | `user not found` |
| 400 | the `*InputError`'s message |
| 500 | `internal error` (never leak internal details) |

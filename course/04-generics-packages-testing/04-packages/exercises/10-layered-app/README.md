# Layered App

Build a small to-do back end in three layers, the same shape as the REST API in Module 7:

```text
todo (app.go)          wiring: picks the implementations (dependency injection)
  └─ service/          business rules; declares the Store interface it needs
       └─ store/       persistence: Task type + in-memory implementation
```

```go
// package store  (store/store.go)
var ErrNotFound = errors.New("store: not found")

type Task struct {
	ID      int
	Title   string
	Done    bool
	Created time.Time
}

type Memory struct{ /* unexported fields */ }

func NewMemory() *Memory
func (m *Memory) Create(t store.Task) (store.Task, error)
func (m *Memory) Get(id int) (store.Task, error)
func (m *Memory) Update(t store.Task) error
func (m *Memory) List() ([]store.Task, error)

// package service  (service/service.go)
const MaxTitle = 80

var (
	ErrEmptyTitle   = errors.New("service: title is empty")
	ErrTitleTooLong = errors.New("service: title is too long")
	ErrDuplicate    = errors.New("service: an open task has this title")
	ErrAlreadyDone  = errors.New("service: task is already done")
)

type Store interface { // declared by the consumer: exactly what Service uses
	Create(t store.Task) (store.Task, error)
	Get(id int) (store.Task, error)
	Update(t store.Task) error
	List() ([]store.Task, error)
}

func New(st Store, now func() time.Time) *Service
func (s *Service) Add(title string) (store.Task, error)
func (s *Service) Complete(id int) error
func (s *Service) Pending() ([]store.Task, error)
func (s *Service) Stats() (done, pending int, err error)

// package todo  (app.go)
func New(now func() time.Time) *service.Service
```

## Examples

| Call (on `todo.New(clock)`) | Result |
|---|---|
| `Add("  Write tests ")` | `Task{ID: 1, Title: "Write tests", Created: clock()}`, `nil` |
| `Add("write TESTS")` | error `ErrDuplicate` (an open task has the same title, ignoring case) |
| `Add("")` / `Add("   ")` | error `ErrEmptyTitle` |
| `Complete(1)` then `Add("Write tests")` | allowed again: the old one is done |
| `Complete(1)` twice | second → `ErrAlreadyDone` |
| `Complete(42)` | error wrapping `store.ErrNotFound` |
| `Stats()` | `done, pending, nil` |

## Rules
- **store.Memory**:
  - `Create` ignores `t.ID`, assigns IDs `1, 2, 3, …`, and returns the stored task.
  - `Get` and `Update` of a missing ID → `fmt.Errorf("%w: task %d", ErrNotFound, id)`.
  - `Update` replaces the whole task.
  - `List` returns all tasks sorted by ID: an empty or nil slice when there are none.
- **service.Add**:
  1. trims spaces (`strings.TrimSpace`);
  2. empty → `ErrEmptyTitle`;
  3. more than `MaxTitle` **runes** → `ErrTitleTooLong`;
  4. an existing **not-done** task with the same title (`strings.EqualFold`) → `ErrDuplicate`;
  5. otherwise creates `store.Task{Title: title, Created: now()}`.
- **Complete**: `Get` the task. If it's already `Done` → `ErrAlreadyDone`; otherwise set `Done` and `Update`.
- **Pending** returns the not-done tasks in ID order. **Stats** counts done and pending tasks.
- **Store errors** (from `Create`, `Get`, `Update`, `List`) are returned **wrapped** (`fmt.Errorf("…: %w", err)`), so `errors.Is` still finds them. A test uses a fake store that fails.
- `service.New` with a nil `now` uses `time.Now`.
- `todo.New` wires `service.New(store.NewMemory(), now)`. It's the only package that knows the concrete store.

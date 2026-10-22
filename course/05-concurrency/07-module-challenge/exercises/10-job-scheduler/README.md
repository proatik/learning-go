# Job Scheduler

Run a set of jobs that depend on each other (a **DAG**, directed acyclic graph) as concurrently as the dependencies allow: a job starts the moment **all** of its dependencies have succeeded. When a job fails, everything that depends on it is skipped, while unrelated jobs carry on.

```go
type Job struct {
	Name string
	Deps []string
	Run  func(ctx context.Context) error
}

type Status int

const (
	Skipped Status = iota // never ran
	Done                  // ran and returned nil
	Failed                // ran and returned an error
)

func (s Status) String() string // "skipped", "done", "failed" (given)

var (
	ErrDuplicate  = errors.New("dag: duplicate job")
	ErrUnknownDep = errors.New("dag: unknown dependency")
	ErrCycle      = errors.New("dag: dependency cycle")
)

func Run(ctx context.Context, jobs []Job) (map[string]Status, error)
```

## Examples

| Jobs (each takes 1s) | Result |
|---|---|
| `a`; `b` and `c` need `a`; `d` needs `b` and `c` | all `done` after **3s** (`b` and `c` run together), `nil` |
| 5 jobs with no dependencies | all `done` after **1s** |
| `a`; `b` needs `a` and fails with `boom`; `c` needs `b`; `d` needs `a` | `a`, `d` done; `b` failed; `c` skipped; error `job "b": boom` |
| `a` needs `b`, `b` needs `a` | nothing runs: `nil` map, error wrapping `ErrCycle` |
| `a` needs `zz` (not in the list) | nothing runs: error wrapping `ErrUnknownDep` |
| two jobs named `a` | nothing runs: error wrapping `ErrDuplicate` |
| no jobs | `map[]`, `nil` |

## Rules
- **Validate first, run nothing if invalid.** Return a `nil` map and an error wrapping (`%w`) `ErrDuplicate`, `ErrUnknownDep` or `ErrCycle` (a job depending on itself is a cycle). Name the offending job in the message, e.g. `dag: unknown dependency: "a" needs "zz"`.
- A job runs, on its own goroutine, as soon as all of its `Deps` are `Done`. Jobs with no dependencies start at once. Every `Run` gets `ctx`.
- If any dependency is `Failed` or `Skipped`, the job is `Skipped` (never run). If `ctx` is done before a job could start, it is `Skipped` too.
- The returned map has one `Status` per job.
- The error is the **first** job failure in time, as `fmt.Errorf("job %q: %w", name, err)`. If no job failed but some were skipped because `ctx` ended, return `ctx.Err()`. Otherwise `nil`.
- `Run` returns after every job has finished or been skipped, with no goroutines left behind.

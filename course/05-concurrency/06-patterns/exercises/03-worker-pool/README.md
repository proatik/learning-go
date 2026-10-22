# Worker Pool

Start a fixed number of workers that share one job channel (fan-out) and send their results to one output channel (fan-in).

```go
type Result struct {
	Job   int
	Value int
}

func Run(n int, jobs <-chan int, work func(int) int) <-chan Result
```

```text
             ┌─▶ worker 1 ─┐
jobs ────────┼─▶ worker 2 ─┼──▶ results   (closed after the last worker exits)
             └─▶ worker n ─┘
```

## Examples

| Code | Result |
|---|---|
| `Run(4, jobs 1..20, func(j int) int { return j * 10 })` | 20 results `{1 10}`, `{2 20}`, … in any order, then closed |
| `Run(3, 10 jobs that sleep 1s, …)` | at most 3 jobs run at once; done after 4s (3+3+3+1) |
| `Run(10, 4 jobs that sleep 1s, …)` | all 4 run at once; done after 1s |
| `Run(0, 3 jobs that sleep 1s, …)` | treated as 1 worker; done after 3s |
| `Run(2, closed empty channel, …)` | closed straight away, no results |

## Rules
- Start exactly `max(n, 1)` workers. Each loops `for j := range jobs`, so it exits when `jobs` is closed and drained.
- Every job produces exactly one `Result{Job: j, Value: work(j)}`. Result order is unspecified.
- While `jobs` is still open, workers keep waiting and the results channel stays open.
- Close the results channel once, after **all** workers have finished. Never close `jobs`: the caller owns it.
- `Run` returns at once.

# Idle Timeout

Collect events until the stream goes **quiet** for `d`: one `time.Timer` that you `Reset` on every event. This is how a server drops an idle connection.

```go
func UntilIdle(events <-chan string, d time.Duration) []string
```

## Examples

All with `d = 250ms`.

| Events arrive at | `events` closed at | Returns | Returns after |
|---|---|---|---|
| `a`@100ms, `b`@300ms, `c`@600ms | 2s | `[a b]` (300→600 is a 300ms gap) | 550ms |
| `a`@200ms, `b`@400ms, `c`@600ms | 5s | `[a b c]` (every gap < 250ms) | 850ms |
| `a`@100ms, `b`@300ms, `c`@500ms, `d`@700ms | 800ms | `[a b c d]` (closed before going idle) | 800ms |
| none | 1s | `nil` | 250ms |
| none | 0s (already closed) | `nil` | 0s |

## Rules
- The idle clock starts when `UntilIdle` is called and restarts at **every** event.
- Return the events received so far (in order) as soon as `d` passes with no event, or as soon as `events` is closed, whichever comes first. Return `nil` if there were none.
- Use a single `time.NewTimer(d)` and call `Reset(d)` on each event. (Since Go 1.23 `Reset` needs no draining.) Stop the timer before returning.
- The tests never send an event at the same instant the timer would fire.

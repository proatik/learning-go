# A fake clock

Write `TestSession` with a **fake clock**, and catch every broken version of `Session`.

```go
// Session expires a fixed TTL after it was started.
type Session struct { /* start, ttl, now */ }

// New starts a session that expires ttl from now(). A nil now means time.Now.
func New(ttl time.Duration, now func() time.Time) *Session

func (s *Session) Expired() bool             // true once now() >= start+ttl
func (s *Session) Remaining() time.Duration  // start+ttl - now(), never below 0
func (s *Session) Touch()                    // restart the TTL from now()
```

`Session` is already written (in the read-only checker). You write the test in `session_test.go`.

## Examples

`s := New(time.Hour, clock)` with the fake clock at 12:00:

| Fake time | `s.Expired()` | `s.Remaining()` |
|---|---|---|
| 12:00 | `false` | `1h0m0s` |
| 12:59:59.999999999 | `false` | `1ns` |
| 13:00 (exactly the deadline) | `true` | `0s` |
| 14:00 | `true` | `0s` |

After `s.Touch()` at 12:40, the new deadline is 13:40.

## Rules

- **Never sleep and never use the real clock.** Pass a func that returns a variable you control:
  `now := time.Date(...)`, `clock := func() time.Time { return now }`, then move time with `now = now.Add(...)`.
- Test just before, exactly at, and after the deadline.
- Test `Touch`: the TTL restarts from the fake time of the call, not from the old deadline.
- **The checker** runs your test on the correct `Session` (it must pass), then once per planted bug (it must fail every time).

# Token Ring

Ping-pong, generalised: `n` player goroutines sit in a **ring**, each connected to the next by a channel. A token (the log of turns so far) travels around the ring; each player adds its name and passes it on. After `turns` turns, `Ring` collects the token and shuts every player down cleanly.

```go
func Ring(names []string, turns int) []string
```

```text
          links[0]            links[1]            links[2]
main ──▶ player 0 ──────────▶ player 1 ──────────▶ player 2 ──┐
            ▲                                                 │
            └──────────────────── links[0] ◀──────────────────┘
```

## Examples

| Call | Result |
|---|---|
| `Ring([ping pong], 4)` | `[ping pong ping pong]` |
| `Ring([a b c], 7)` | `[a b c a b c a]` |
| `Ring([solo], 3)` | `[solo solo solo]` (a player passes to itself) |
| `Ring([a b], 0)` | `[]` |
| `Ring([a b], -2)` | `[]` (same as `0`) |
| `Ring([], 5)` | `[]` |

## Rules
- Start **one goroutine per player**. Player `i` receives the token on `links[i]` and sends it on `links[(i+1) % n]`. Players only talk through these channels.
- The token is a `[]string`. It starts **empty** at player `0`. A player that receives it either:
  - finds `len(token) == turns`: the game is over, so it hands the token back to `Ring`, or
  - appends its own name and passes the token to the next player.
- `Ring` returns the finished token. `turns <= 0` returns an empty log; no names returns an empty log without starting anything.
- ⚠️ With **one** player, `links[0]` is both its input and its output. An unbuffered send to yourself blocks forever (nobody else can receive), so it must still work.
- **No leaks:** when `Ring` returns, every player goroutine has exited. A player blocked forever on a receive is a leak (the tests count goroutines).

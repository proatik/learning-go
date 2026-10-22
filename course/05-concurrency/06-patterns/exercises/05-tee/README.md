# Tee Channel

Split one channel into two: like the Unix `tee` command, every value from `in` goes to **both** outputs.

```go
func Tee[T any](ctx context.Context, in <-chan T) (<-chan T, <-chan T)
```

```text
            ┌──▶ out1   1, 2, 3
in 1,2,3 ───┤
            └──▶ out2   1, 2, 3
```

## Examples

| Code | Result |
|---|---|
| `Tee(ctx, ch(1, 2, 3))` | both outputs yield `1 2 3`, then both close |
| `Tee(ctx, ch[int]())` | both outputs close straight away |
| reader takes from `out2` first, then `out1`, for every value | works: no deadlock |
| `in` never closes, then `cancel()` | both outputs close, the goroutine exits |

(`ch(…)` is a closed channel holding those values.)

## Rules
- Both outputs are **unbuffered** and receive every value, in `in`'s order.
- A value goes to both outputs **before** the next value is read from `in`. Deliver it in whichever order the readers are ready: a `select` over both sends, setting a local copy of a channel to `nil` once it has had its copy (a nil channel's case never fires).
- Close both outputs when `in` is closed **or** `ctx` is done. Every blocking operation (the receive from `in` and every send) must also watch `ctx.Done()`, so a cancelled Tee never leaks its goroutine.
- Tee owns its outputs and closes them; it never closes `in`.

# Stack on a Slice Pointer

Implement a LIFO stack stored in a plain `[]int`. Functions that change the stack's **length** take a `*[]int`; read-only ones take the slice itself.

```go
func Push(s *[]int, v int)
func PushAll(s *[]int, vs ...int)
func Pop(s *[]int) (int, bool)
func Peek(s []int) (int, bool)
```

The top of the stack is the **last** element of the slice.

| Function | Behavior |
|---|---|
| `Push(&s, v)` | appends `v` |
| `PushAll(&s, 1, 2, 3)` | pushes in order, so `3` ends on top |
| `Pop(&s)` | removes and returns the top, `true`; on an empty stack returns `0, false` and leaves `s` alone |
| `Peek(s)` | returns the top without removing it; `0, false` when empty |

```go
var s []int       // nil slice = empty stack
Push(&s, 1)
Push(&s, 2)       // s == [1 2]
v, ok := Pop(&s)  // 2, true; s == [1]
v, ok = Peek(s)   // 1, true; s unchanged
Pop(&s)           // 1, true; s == []
v, ok = Pop(&s)   // 0, false
```

## Rules
- `s` (the pointer) is never nil; `*s` may be a nil slice.
- With `s *[]int`, write `*s = append(*s, v)`, `len(*s)`, `(*s)[i]`: indexing doesn't auto-dereference slice pointers.
- `Pop` shrinks with a re-slice: `*s = (*s)[:n-1]`.

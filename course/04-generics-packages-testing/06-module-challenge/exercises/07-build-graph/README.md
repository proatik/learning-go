# Build Planner on a Generic Graph

Write a reusable generic `graph` package (BFS iterator, topological sort, a generic error type), then use it from the root `build` package to plan build order and find which targets a change affects.

```go
// package graph — learngo/course/04-generics-packages-testing/06-module-challenge/exercises/07-build-graph/graph
var ErrCycle = errors.New("graph: cycle")

type CycleError[N cmp.Ordered] struct {
	Nodes []N // sorted
}

func (e *CycleError[N]) Error() string
func (e *CycleError[N]) Is(target error) bool

type Graph[N cmp.Ordered] struct {
	adj map[N][]N // node → successors in insertion order; every node is a key
}

func New[N cmp.Ordered]() *Graph[N]
func (g *Graph[N]) AddNode(n N)
func (g *Graph[N]) AddEdge(from, to N)
func (g *Graph[N]) Nodes() iter.Seq[N]
func (g *Graph[N]) Neighbors(n N) iter.Seq[N]
func (g *Graph[N]) BFS(start N) iter.Seq2[N, int]
func (g *Graph[N]) TopoSort() ([]N, error)

// package build — the exercise root
var ErrUnknownDep = errors.New("build: unknown dependency")

type Target struct {
	Name string
	Deps []string
}

func Plan(targets []Target) ([]string, error)
func Affected(targets []Target, changed string) ([]string, error)
```

## Examples

| Code | Result |
|---|---|
| edges `a→b a→c b→d c→d`: `BFS("a")` | `a 0`, `b 1`, `c 1`, `d 2` |
| edges `3→1 2→1 1→0 10→0`: `TopoSort()` | `[2 3 1 10 0]` — smallest ready node first |
| edges `a→b b→c c→a c→d e→a`: `TopoSort()` | `nil`, `"graph: cycle among [a b c d]"` |
| `Plan([{app [lib]} {lib [core]} {core []}])` | `[core lib app]` |
| `Plan([{app [lib]}])` | `build: unknown dependency: "app" needs "lib"` |
| `Plan([{a [b]} {b [a]}])` | `"build: graph: cycle among [a b]"` |
| `Affected(targets, "lib")` | `lib` and everything depending on it, sorted |

## Rules
- **graph** knows nothing about builds; **build** imports **graph** (never the other way round).
- **AddEdge** adds missing nodes and ignores a repeated edge. **Nodes** yields sorted nodes; **Neighbors** yields successors in insertion order (none for an unknown node).
- **BFS** yields `(node, distance)` breadth first, each reachable node once, successors in insertion order; nothing for an unknown `start`. It must terminate on cycles and honour `break`.
- **TopoSort** uses Kahn's algorithm: repeatedly take the **smallest** node whose remaining in-degree is 0. Empty graph → empty slice, `nil`.
- On a cycle, **TopoSort** returns `nil` and `&CycleError{Nodes: ...}` holding every node never ordered, sorted. `Error()` is `fmt.Sprintf("graph: cycle among %v", e.Nodes)`; `Is` returns `target == ErrCycle`.
- **build** adds every target as a node and an edge **dep → target** for each dependency. A dependency that is not a target → `fmt.Errorf("%w: %q needs %q", ErrUnknownDep, target, dep)`, checking targets and deps in input order.
- **Plan** wraps a graph error as `fmt.Errorf("build: %w", err)` so `errors.As` still finds `*graph.CycleError[string]`.
- **Affected** returns the nodes `BFS` reaches from `changed`, sorted; an unknown `changed` gives an empty, non-nil slice.

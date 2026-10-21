# Graph Shortest Path

A **directed graph** is stored as an adjacency list: `g["a"]` lists the nodes `a` has an edge **to**. Build graphs and find shortest paths with breadth-first search (BFS).

```go
var (
	ErrUnknownNode = errors.New("graph: unknown node")
	ErrNoPath      = errors.New("graph: no path")
)

func AddEdge(g map[string][]string, from, to string)
func Distances(g map[string][]string, from string) map[string]int
func ShortestPath(g map[string][]string, from, to string) ([]string, error)
```

## Examples

```go
g := map[string][]string{}
AddEdge(g, "home", "park")
AddEdge(g, "home", "shop")
AddEdge(g, "park", "school")
AddEdge(g, "shop", "school")
AddEdge(g, "school", "library")

ShortestPath(g, "home", "library") // [home park school library], nil
ShortestPath(g, "library", "home") // nil, ErrNoPath
Distances(g, "home")               // map[home:0 park:1 shop:1 school:2 library:3]
```

## Rules
- A node **exists** if it is a key of `g` **or** appears in any adjacency list. A neighbour that is not a key simply has no outgoing edges.
- **AddEdge** appends `to` to `g[from]` unless it is already there (no duplicate edges), and makes sure `to` is also a key (`g[to]` exists, possibly with an empty list). Self-loops are allowed. The caller's map is modified — maps are reference-like, no pointer needed.
- **Distances** returns the number of edges on the shortest path from `from` to every node reachable from it, including `from` itself at `0`. If `from` does not exist, return `nil`.
- **ShortestPath** returns the nodes of a shortest path, `from` and `to` included. `from == to` (existing) gives `[from]`.
  - `from` or `to` does not exist → `nil, ErrUnknownNode`.
  - `to` is not reachable → `nil, ErrNoPath`.
  - **Ties:** BFS visits each node's neighbours in **slice order**, and a node's predecessor is fixed the **first time** it is discovered. Return the path this produces (in the example, `park` is listed before `shop`, so the path goes through `park`).
- Return the sentinel errors themselves (or wrap them with `%w`); tests use `errors.Is`.
- Graphs may contain cycles. Do not modify `g` in `Distances` or `ShortestPath`.

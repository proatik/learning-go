# sort.Interface, Old and New

Implement `sort.Interface` the classic way, then do the same job with the modern `slices.SortFunc`.

```go
type Person struct {
	Name string
	Age  int
}

type ByAge []Person  // implements sort.Interface: youngest first
type ByName []Person // implements sort.Interface: A → Z

func (a ByAge) Len() int
func (a ByAge) Less(i, j int) bool
func (a ByAge) Swap(i, j int)
// …and the same three methods for ByName

func SortOldestFirst(ps []Person)   // sort package: oldest first, ties keep their order
func SortAgeThenName(ps []Person)   // slices.SortFunc: age ascending, then name A → Z
```

## Examples

Starting from `ps := []Person{{"Cy", 40}, {"Ann", 31}, {"Bob", 25}, {"Al", 31}}`:

| Call | `ps` afterwards |
|---|---|
| `sort.Sort(ByAge(ps))` | `[{Bob 25} {Ann 31} {Al 31} {Cy 40}]` (the two 31s in either order) |
| `sort.Sort(ByName(ps))` | `[{Al 31} {Ann 31} {Bob 25} {Cy 40}]` |
| `SortOldestFirst(ps)` | `[{Cy 40} {Ann 31} {Al 31} {Bob 25}]` |
| `SortAgeThenName(ps)` | `[{Bob 25} {Al 31} {Ann 31} {Cy 40}]` |

## Rules
- `ByAge.Less` compares only `Age`; `ByName.Less` compares only `Name` (plain `<` on strings).
- `SortOldestFirst` must use the `sort` package with your `ByAge` (hint: `sort.Reverse`, and a **stable** sort so equal ages keep their original relative order).
- `SortAgeThenName` must use `slices.SortFunc` with `cmp.Compare` (no `sort` package).
- All functions sort **in place**; nil and empty slices are fine.

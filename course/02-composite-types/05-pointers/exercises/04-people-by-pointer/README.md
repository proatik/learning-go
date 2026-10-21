# People by Pointer

Update structs in place: through a `*Person`, and inside a `[]Person` by index. Return pointers **into** the slice so callers can edit the element they found.

```go
type Person struct {
	Name string
	Age  int
}

func Birthday(p *Person)
func BirthdayAll(people []Person)
func Find(people []Person, name string) *Person
func Oldest(people []Person) *Person
```

| Function | Behavior |
|---|---|
| `Birthday(p)` | `p.Age` grows by 1; does nothing if `p` is nil |
| `BirthdayAll(people)` | every person in the slice gets one year older |
| `Find(people, name)` | pointer to the **first** element with that `Name`, or `nil` |
| `Oldest(people)` | pointer to the element with the highest `Age` (first one on ties), or `nil` if empty |

```go
people := []Person{{"Ana", 30}, {"Ben", 41}, {"Cy", 41}}
p := Find(people, "Ana")
p.Age = 99          // people[0].Age is now 99: p points INTO the slice
o := Oldest(people) // &people[0] (99)
```

## Rules
- The returned pointers must point to the slice's own elements (`&people[i]`), not to copies.
- `for _, p := range people` gives copies: use the index.
- Reuse `Birthday` in `BirthdayAll`.

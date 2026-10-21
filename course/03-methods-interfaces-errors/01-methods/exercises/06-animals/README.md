# Animals and Embedding

Methods promoted through embedding, shadowed by the outer type, and extended by calling the embedded version.

```go
type Animal struct {
	Name  string
	Sound string
}

func (a Animal) Speak() string     // "Rex says woof"
func (a Animal) Describe() string  // "Rex is an animal"
func (a *Animal) Rename(name string)

type Dog struct {
	Animal
	Breed string
}

func NewDog(name, breed string) Dog // Sound is "woof"
func (d Dog) Describe() string      // shadows Animal.Describe
func (d Dog) Fetch() string         // "Rex fetches the ball"

type GuideDog struct {
	Dog
	Handler string
}

func (g GuideDog) Describe() string // extends Dog.Describe
```

## Examples

| Code | Result |
|---|---|
| `Animal{"Tom", "meow"}.Speak()` | `"Tom says meow"` |
| `Animal{Name: "Tom"}.Describe()` | `"Tom is an animal"` |
| `NewDog("Rex", "beagle").Describe()` | `"Rex is a beagle"` |
| `NewDog("Rex", "").Describe()` | `"Rex is a dog"` |
| `NewDog("Rex", "beagle").Speak()` | `"Rex says woof"` (promoted) |
| `NewDog("Rex", "beagle").Animal.Describe()` | `"Rex is an animal"` |
| `GuideDog{NewDog("Max", "labrador"), "Sam"}.Describe()` | `"Max is a labrador guiding Sam"` |
| `g.Rename("Buddy")` on a `GuideDog` | `g.Name` becomes `"Buddy"` |

## Rules
- `Speak`: `"<Name> says <Sound>"`; `Animal.Describe`: `"<Name> is an animal"`.
- `Dog.Describe`: `"<Name> is a <Breed>"`, or `"<Name> is a dog"` when `Breed` is empty.
- `GuideDog.Describe` must **reuse** `Dog.Describe` and append `" guiding <Handler>"`; with an empty `Handler`, it returns `Dog.Describe` unchanged.
- `Rename` changes `Name` in place; it is promoted to `Dog` and `GuideDog` automatically.

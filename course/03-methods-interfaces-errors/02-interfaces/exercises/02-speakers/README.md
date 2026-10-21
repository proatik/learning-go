# Speakers

Make three unrelated types satisfy one `Speaker` interface, then write code that handles them all the same way.

```go
type Speaker interface {
	Name() string
	Sound() string
}

type Dog struct{ Nick string }
type Cat struct {
	Nick   string
	Grumpy bool
}
type Robot struct {
	Model   string
	Version int
}

func Greet(s Speaker) string
func Chorus(speakers []Speaker) string
func CountSound(speakers []Speaker, sound string) int
```

## Examples

| Call | Result |
|---|---|
| `Dog{"Rex"}.Name()` / `.Sound()` | `"Rex"` / `"Woof"` |
| `Cat{"Tom", false}.Sound()` | `"Meow"` |
| `Cat{"Tom", true}.Sound()` | `"Hiss"` |
| `Robot{"RX", 2}.Name()` / `.Sound()` | `"RX-v2"` / `"Beep boop"` |
| `Greet(Dog{"Rex"})` | `"Rex says Woof!"` |
| `Chorus([]Speaker{Dog{"Rex"}, Robot{"RX", 2}})` | `"Rex says Woof! RX-v2 says Beep boop!"` |
| `Chorus(nil)` | `""` |
| `CountSound([]Speaker{Dog{"A"}, Cat{"B", false}, Dog{"C"}}, "Woof")` | `2` |

## Rules
- All methods use **value receivers**.
- `Robot.Name()` is `Model + "-v" + Version`.
- `Greet` returns `"<Name> says <Sound>!"`.
- `Chorus` joins the greetings with a single space, in order; empty or nil input → `""`.
- `CountSound` counts speakers whose `Sound()` equals `sound` exactly (case-sensitive).

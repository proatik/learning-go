# Nested Maps: A Gradebook

A gradebook maps **student → (subject → score)**:

```go
type Gradebook map[string]map[string]int

func AddGrade(g Gradebook, student, subject string, score int)
func Average(g Gradebook, student string) (avg float64, ok bool)
func SubjectAverages(g Gradebook) map[string]float64
func PrintReport(g Gradebook)
```

## 1. `AddGrade`

Record (or overwrite) `student`'s score in `subject`. `g` is non-nil, but the **inner map for a new student doesn't exist yet**: writing `g[student][subject]` straight away panics.

```go
g := Gradebook{}
AddGrade(g, "alice", "math", 90)
AddGrade(g, "alice", "physics", 85)
AddGrade(g, "bob", "math", 70)
AddGrade(g, "bob", "math", 75) // overwrite
// g == Gradebook{"alice": {"math": 90, "physics": 85}, "bob": {"math": 75}}
```

## 2. `Average`

The mean of one student's scores. `ok` is `false` (and `avg` is `0`) if the student is unknown **or** has no grades.

| Call (with `g` above) | Result |
|---|---|
| `Average(g, "alice")` | `87.5, true` |
| `Average(g, "bob")` | `75, true` |
| `Average(g, "zoe")` | `0, false` |

## 3. `SubjectAverages`

For every subject anyone has a score in, the mean over the students who have one. Never nil.

`SubjectAverages(g)` → `map[math:82.5 physics:85]`

## 4. `PrintReport`

One line per student, **students sorted**, and within a line **subjects sorted**, then the average with two decimals. A student with an empty inner map prints `(no grades)`. An empty gradebook prints nothing.

```
alice: math=90 physics=85 (avg 87.50)
bob: math=75 (avg 75.00)
carl: (no grades)
```

## Rules
- Use comma-ok to detect a missing inner map, create it, store it, then write.
- Reading through a missing student is safe (`g["zoe"]["math"]` is `0`); writing is not.
- Every output must be deterministic: sort keys with `slices.Sorted(maps.Keys(...))`.
- Build each report line with `strings.Builder` or `fmt.Sprintf` + `strings.Join`, whichever you prefer.

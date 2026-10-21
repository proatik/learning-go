# Student Ranking

Sort a slice of structs by **two keys** and compute competition ranks.

```go
type Student struct {
	Name  string
	Grade int
}

// Ranked is a Student plus its rank; Student is embedded, so r.Name and r.Grade work.
type Ranked struct {
	Student
	Rank int
}

func SortStudents(students []Student)              // in place
func TopN(students []Student, n int) []Student     // new slice, input untouched
func Ranks(students []Student) []Ranked            // new slice, input untouched
```

**Order:** grade **descending**, then name **ascending** (byte order: `"Bob" < "ann"`).

| Input | Result |
|---|---|
| `SortStudents` on `[{cy 70} {bob 90} {amy 90}]` | becomes `[{amy 90} {bob 90} {cy 70}]` |
| `TopN([{cy 70} {bob 90} {amy 90}], 2)` | `[{amy 90} {bob 90}]` |
| `TopN(..., 10)` | all 3, sorted |
| `TopN(..., 0)` | empty |
| `Ranks([{cy 70} {bob 90} {amy 90} {dan 60}])` | `[{amy 90} 1, {bob 90} 1, {cy 70} 3, {dan 60} 4]` |

## Rules
- Use `slices.SortFunc` with `cmp.Or` and `cmp.Compare`.
- `TopN` and `Ranks` must **not** reorder or change the caller's slice. `TopN` with `n <= 0` returns an empty (or nil) slice.
- Ranks use **standard competition ranking** ("1224"): equal grades share a rank; the next rank skips accordingly. A rank equals 1 + the number of students with a strictly higher grade.

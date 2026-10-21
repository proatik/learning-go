# Employee Reports

Work with a **slice of structs**: filter it, aggregate it into a map, and find the maximum.

```go
type Employee struct {
	Name   string
	Dept   string
	Salary float64
}

func FilterByDept(emps []Employee, dept string) []Employee
func TotalByDept(emps []Employee) map[string]float64
func HighestPaid(emps []Employee) (Employee, bool)
```

Given

```go
emps := []Employee{
	{Name: "Ann", Dept: "eng", Salary: 120},
	{Name: "Bob", Dept: "ops", Salary: 80},
	{Name: "Cy", Dept: "eng", Salary: 100},
}
```

| Call | Result |
|---|---|
| `FilterByDept(emps, "eng")` | `[{Ann eng 120} {Cy eng 100}]` |
| `FilterByDept(emps, "hr")` | empty slice |
| `TotalByDept(emps)` | `map[eng:220 ops:80]` |
| `HighestPaid(emps)` | `{Ann eng 120}, true` |
| `HighestPaid(nil)` | `Employee{}, false` |

## Rules
- `FilterByDept` keeps the input order and must not modify `emps`. Department names match exactly (case-sensitive).
- `TotalByDept` always returns a **non-nil** map (empty for no employees).
- `HighestPaid`: on a tie, the **first** employee with the highest salary wins.

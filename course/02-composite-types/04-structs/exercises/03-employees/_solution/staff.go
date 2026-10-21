package staff

// Employee is a member of staff.
type Employee struct {
	Name   string
	Dept   string
	Salary float64
}

// FilterByDept returns the employees of dept, in input order.
func FilterByDept(emps []Employee, dept string) []Employee {
	var out []Employee // appending to a fresh slice never touches emps
	for _, e := range emps {
		if e.Dept == dept {
			out = append(out, e)
		}
	}
	return out
}

// TotalByDept returns the total salary per department. The map is never nil.
func TotalByDept(emps []Employee) map[string]float64 {
	totals := make(map[string]float64)
	for _, e := range emps {
		totals[e.Dept] += e.Salary // a missing key reads as 0
	}
	return totals
}

// HighestPaid returns the first employee with the highest salary and true,
// or the zero Employee and false if emps is empty.
func HighestPaid(emps []Employee) (Employee, bool) {
	if len(emps) == 0 {
		return Employee{}, false
	}
	best := emps[0]
	for _, e := range emps[1:] {
		if e.Salary > best.Salary { // strictly greater: ties keep the first
			best = e
		}
	}
	return best, true
}

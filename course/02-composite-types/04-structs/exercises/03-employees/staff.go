package staff

// Employee is a member of staff.
type Employee struct {
	Name   string
	Dept   string
	Salary float64
}

// FilterByDept returns the employees of dept, in input order.
func FilterByDept(emps []Employee, dept string) []Employee {
	// TODO
	return nil
}

// TotalByDept returns the total salary per department. The map is never nil.
func TotalByDept(emps []Employee) map[string]float64 {
	// TODO
	return nil
}

// HighestPaid returns the first employee with the highest salary and true,
// or the zero Employee and false if emps is empty.
func HighestPaid(emps []Employee) (Employee, bool) {
	// TODO
	return Employee{}, false
}

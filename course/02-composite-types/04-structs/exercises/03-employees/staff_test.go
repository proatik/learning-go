package staff

import (
	"maps"
	"slices"
	"testing"
)

func sample() []Employee {
	return []Employee{
		{Name: "Ann", Dept: "eng", Salary: 120},
		{Name: "Bob", Dept: "ops", Salary: 80},
		{Name: "Cy", Dept: "eng", Salary: 100},
		{Name: "Di", Dept: "sales", Salary: 120},
		{Name: "Ed", Dept: "ops", Salary: 95.5},
		{Name: "Flo", Dept: "Eng", Salary: 70},
	}
}

func TestFilterByDept(t *testing.T) {
	tests := []struct {
		dept string
		want []Employee
	}{
		{"eng", []Employee{{Name: "Ann", Dept: "eng", Salary: 120}, {Name: "Cy", Dept: "eng", Salary: 100}}},
		{"ops", []Employee{{Name: "Bob", Dept: "ops", Salary: 80}, {Name: "Ed", Dept: "ops", Salary: 95.5}}},
		{"Eng", []Employee{{Name: "Flo", Dept: "Eng", Salary: 70}}},
		{"hr", nil},
		{"", nil},
	}
	for _, tt := range tests {
		emps := sample()
		got := FilterByDept(emps, tt.dept)
		if !slices.Equal(got, tt.want) {
			t.Errorf("FilterByDept(sample, %q) = %+v, want %+v", tt.dept, got, tt.want)
		}
		if !slices.Equal(emps, sample()) {
			t.Errorf("FilterByDept(sample, %q) modified its input: %+v", tt.dept, emps)
		}
	}
	if got := FilterByDept(nil, "eng"); len(got) != 0 {
		t.Errorf("FilterByDept(nil, \"eng\") = %+v, want empty", got)
	}
}

func TestTotalByDept(t *testing.T) {
	got := TotalByDept(sample())
	want := map[string]float64{"eng": 220, "ops": 175.5, "sales": 120, "Eng": 70}
	if !maps.Equal(got, want) {
		t.Errorf("TotalByDept(sample) = %v, want %v", got, want)
	}

	one := []Employee{{Name: "Zed", Dept: "x", Salary: 0}}
	got = TotalByDept(one)
	if v, ok := got["x"]; !ok || v != 0 {
		t.Errorf("TotalByDept(%+v) = %v, want map[x:0] (a department with 0 total still appears)", one, got)
	}

	for _, in := range [][]Employee{nil, {}} {
		got := TotalByDept(in)
		if got == nil {
			t.Errorf("TotalByDept(%#v) = nil, want a non-nil empty map", in)
		} else if len(got) != 0 {
			t.Errorf("TotalByDept(%#v) = %v, want an empty map", in, got)
		}
	}
}

func TestHighestPaid(t *testing.T) {
	tests := []struct {
		emps   []Employee
		want   Employee
		wantOK bool
	}{
		{sample(), Employee{Name: "Ann", Dept: "eng", Salary: 120}, true},
		{[]Employee{{Name: "Bob", Dept: "ops", Salary: 80}}, Employee{Name: "Bob", Dept: "ops", Salary: 80}, true},
		{[]Employee{{Name: "A", Salary: 1}, {Name: "B", Salary: 3}, {Name: "C", Salary: 2}}, Employee{Name: "B", Salary: 3}, true},
		{[]Employee{{Name: "A", Salary: 5}, {Name: "B", Salary: 9}, {Name: "C", Salary: 9}}, Employee{Name: "B", Salary: 9}, true},
		{[]Employee{{Name: "Neg", Salary: -5}, {Name: "Neg2", Salary: -1}}, Employee{Name: "Neg2", Salary: -1}, true},
		{[]Employee{}, Employee{}, false},
		{nil, Employee{}, false},
	}
	for _, tt := range tests {
		got, ok := HighestPaid(tt.emps)
		if got != tt.want || ok != tt.wantOK {
			t.Errorf("HighestPaid(%+v) = %+v, %v; want %+v, %v", tt.emps, got, ok, tt.want, tt.wantOK)
		}
	}
}

package gradebook

import (
	"maps"
	"reflect"
	"testing"
)

func TestAddGrade(t *testing.T) {
	g := Gradebook{}
	AddGrade(g, "alice", "math", 90) // must not panic: the inner map is created
	AddGrade(g, "alice", "physics", 85)
	AddGrade(g, "bob", "math", 70)
	AddGrade(g, "bob", "math", 75)
	AddGrade(g, "carl", "art", 0)

	want := Gradebook{
		"alice": {"math": 90, "physics": 85},
		"bob":   {"math": 75},
		"carl":  {"art": 0},
	}
	if !reflect.DeepEqual(g, want) {
		t.Errorf("after AddGrade calls, g = %v, want %v", g, want)
	}
}

func TestAddGradeKeepsStudentsSeparate(t *testing.T) {
	g := Gradebook{}
	AddGrade(g, "alice", "math", 90)
	AddGrade(g, "bob", "art", 60)
	if len(g["alice"]) != 1 || len(g["bob"]) != 1 {
		t.Errorf("g = %v; each student needs their own inner map", g)
	}
}

func TestAddGradeExistingInnerMap(t *testing.T) {
	inner := map[string]int{"math": 50}
	g := Gradebook{"dana": inner}
	AddGrade(g, "dana", "chem", 77)
	want := map[string]int{"math": 50, "chem": 77}
	if !maps.Equal(g["dana"], want) {
		t.Errorf("g[\"dana\"] = %v, want %v (keep existing scores)", g["dana"], want)
	}
	if !maps.Equal(inner, want) {
		t.Errorf("existing inner map = %v, want %v: write into it, don't replace it", inner, want)
	}
}

var sample = Gradebook{
	"alice": {"math": 90, "physics": 85},
	"bob":   {"math": 75},
	"carl":  {},
	"dana":  {"math": 60, "art": 100, "chem": 71, "physics": 89},
}

func TestAverage(t *testing.T) {
	tests := []struct {
		student string
		want    float64
		wantOK  bool
	}{
		{"alice", 87.5, true},
		{"bob", 75, true},
		{"dana", 80, true},
		{"carl", 0, false}, // known student, no grades
		{"zoe", 0, false},  // unknown student
		{"", 0, false},
	}
	for _, tt := range tests {
		got, ok := Average(sample, tt.student)
		if got != tt.want || ok != tt.wantOK {
			t.Errorf("Average(sample, %q) = (%v, %t), want (%v, %t)", tt.student, got, ok, tt.want, tt.wantOK)
		}
	}
	if got, ok := Average(nil, "alice"); got != 0 || ok {
		t.Errorf("Average(nil, \"alice\") = (%v, %t), want (0, false)", got, ok)
	}
}

func TestSubjectAverages(t *testing.T) {
	tests := []struct {
		name string
		g    Gradebook
		want map[string]float64
	}{
		{"sample", sample, map[string]float64{"math": 75, "physics": 87, "art": 100, "chem": 71}},
		{
			"two students",
			Gradebook{"alice": {"math": 90, "physics": 85}, "bob": {"math": 75}},
			map[string]float64{"math": 82.5, "physics": 85},
		},
		{"one empty student", Gradebook{"x": {}}, map[string]float64{}},
		{"empty", Gradebook{}, map[string]float64{}},
		{"nil", nil, map[string]float64{}},
	}
	for _, tt := range tests {
		got := SubjectAverages(tt.g)
		if got == nil {
			t.Errorf("SubjectAverages(%s) = nil, want an initialized map", tt.name)
			continue
		}
		if !maps.Equal(got, tt.want) {
			t.Errorf("SubjectAverages(%s) = %v, want %v", tt.name, got, tt.want)
		}
	}
}

func ExamplePrintReport() {
	PrintReport(Gradebook{
		"bob":   {"math": 75},
		"carl":  {},
		"alice": {"physics": 85, "math": 90},
	})
	// Output:
	// alice: math=90 physics=85 (avg 87.50)
	// bob: math=75 (avg 75.00)
	// carl: (no grades)
}

func ExamplePrintReport_sorted() {
	g := Gradebook{}
	AddGrade(g, "zed", "physics", 70)
	AddGrade(g, "zed", "art", 81)
	AddGrade(g, "zed", "math", 66)
	AddGrade(g, "amy", "chem", 99)
	AddGrade(g, "amy", "biology", 98)
	PrintReport(g)
	// Output:
	// amy: biology=98 chem=99 (avg 98.50)
	// zed: art=81 math=66 physics=70 (avg 72.33)
}

func ExamplePrintReport_empty() {
	PrintReport(Gradebook{})
	PrintReport(nil)
	// Output:
}

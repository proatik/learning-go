package ranking

import (
	"slices"
	"testing"
)

func class() []Student {
	return []Student{
		{Name: "cy", Grade: 70},
		{Name: "bob", Grade: 90},
		{Name: "amy", Grade: 90},
		{Name: "dan", Grade: 60},
		{Name: "eve", Grade: 70},
		{Name: "Zed", Grade: 70},
	}
}

func TestSortStudents(t *testing.T) {
	tests := []struct {
		in, want []Student
	}{
		{
			class(),
			[]Student{{"amy", 90}, {"bob", 90}, {"Zed", 70}, {"cy", 70}, {"eve", 70}, {"dan", 60}},
		},
		{
			[]Student{{"cy", 70}, {"bob", 90}, {"amy", 90}},
			[]Student{{"amy", 90}, {"bob", 90}, {"cy", 70}},
		},
		{
			[]Student{{"b", 1}, {"a", 1}, {"c", 1}},
			[]Student{{"a", 1}, {"b", 1}, {"c", 1}},
		},
		{
			[]Student{{"low", -5}, {"high", 100}, {"mid", 0}},
			[]Student{{"high", 100}, {"mid", 0}, {"low", -5}},
		},
		{[]Student{{"solo", 50}}, []Student{{"solo", 50}}},
		{nil, nil},
	}
	for _, tt := range tests {
		got := slices.Clone(tt.in)
		SortStudents(got)
		if !slices.Equal(got, tt.want) {
			t.Errorf("SortStudents(%+v) gave %+v, want %+v", tt.in, got, tt.want)
		}
	}
}

func TestTopN(t *testing.T) {
	tests := []struct {
		n    int
		want []Student
	}{
		{1, []Student{{"amy", 90}}},
		{2, []Student{{"amy", 90}, {"bob", 90}}},
		{4, []Student{{"amy", 90}, {"bob", 90}, {"Zed", 70}, {"cy", 70}}},
		{6, []Student{{"amy", 90}, {"bob", 90}, {"Zed", 70}, {"cy", 70}, {"eve", 70}, {"dan", 60}}},
		{100, []Student{{"amy", 90}, {"bob", 90}, {"Zed", 70}, {"cy", 70}, {"eve", 70}, {"dan", 60}}},
		{0, nil},
		{-1, nil},
	}
	for _, tt := range tests {
		in := class()
		got := TopN(in, tt.n)
		if !slices.Equal(got, tt.want) {
			t.Errorf("TopN(class, %d) = %+v, want %+v", tt.n, got, tt.want)
		}
		if !slices.Equal(in, class()) {
			t.Errorf("TopN(class, %d) modified its input: %+v", tt.n, in)
		}
	}
	if got := TopN(nil, 3); len(got) != 0 {
		t.Errorf("TopN(nil, 3) = %+v, want empty", got)
	}
}

func TestRanks(t *testing.T) {
	tests := []struct {
		in   []Student
		want []Ranked
	}{
		{
			class(),
			[]Ranked{
				{Student{"amy", 90}, 1},
				{Student{"bob", 90}, 1},
				{Student{"Zed", 70}, 3},
				{Student{"cy", 70}, 3},
				{Student{"eve", 70}, 3},
				{Student{"dan", 60}, 6},
			},
		},
		{
			[]Student{{"cy", 70}, {"bob", 90}, {"amy", 90}, {"dan", 60}},
			[]Ranked{{Student{"amy", 90}, 1}, {Student{"bob", 90}, 1}, {Student{"cy", 70}, 3}, {Student{"dan", 60}, 4}},
		},
		{
			[]Student{{"a", 3}, {"b", 2}, {"c", 1}},
			[]Ranked{{Student{"a", 3}, 1}, {Student{"b", 2}, 2}, {Student{"c", 1}, 3}},
		},
		{
			[]Student{{"x", 5}, {"y", 5}, {"z", 5}},
			[]Ranked{{Student{"x", 5}, 1}, {Student{"y", 5}, 1}, {Student{"z", 5}, 1}},
		},
		{[]Student{{"solo", 0}}, []Ranked{{Student{"solo", 0}, 1}}},
		{nil, nil},
	}
	for _, tt := range tests {
		in := slices.Clone(tt.in)
		got := Ranks(in)
		if !slices.Equal(got, tt.want) {
			t.Errorf("Ranks(%+v) = %+v, want %+v", tt.in, got, tt.want)
		}
		if !slices.Equal(in, tt.in) {
			t.Errorf("Ranks(%+v) modified its input: %+v", tt.in, in)
		}
	}
}

func TestRankedPromotedFields(t *testing.T) {
	got := Ranks([]Student{{"amy", 90}})
	if len(got) != 1 {
		t.Fatalf("Ranks([{amy 90}]) returned %d results, want 1", len(got))
	}
	if got[0].Name != "amy" || got[0].Grade != 90 || got[0].Rank != 1 {
		t.Errorf("Ranks([{amy 90}])[0] = %+v, want Name amy, Grade 90, Rank 1", got[0])
	}
}

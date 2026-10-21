package ranking

import (
	"slices"
	"testing"
)

var rankTests = []struct {
	names         []string
	scores        []int
	wantRank      []string
	wantStandings []string
}{
	{
		[]string{"cat", "amy", "bob", "dan"},
		[]int{80, 95, 95, 70},
		[]string{"amy", "bob", "cat", "dan"},
		[]string{"1. amy (95)", "1. bob (95)", "3. cat (80)", "4. dan (70)"},
	},
	{
		[]string{"x", "y"},
		[]int{1, 2},
		[]string{"y", "x"},
		[]string{"1. y (2)", "2. x (1)"},
	},
	{
		[]string{"c", "b", "a"},
		[]int{5, 5, 5},
		[]string{"a", "b", "c"},
		[]string{"1. a (5)", "1. b (5)", "1. c (5)"},
	},
	{
		[]string{"eve", "dan", "cy", "bo", "al"},
		[]int{60, 90, 60, 90, 75},
		[]string{"bo", "dan", "al", "cy", "eve"},
		[]string{"1. bo (90)", "1. dan (90)", "3. al (75)", "4. cy (60)", "4. eve (60)"},
	},
	{
		[]string{"zed", "Zed", "ann"},
		[]int{10, 10, 0},
		[]string{"Zed", "zed", "ann"}, // 'Z' sorts before 'z'
		[]string{"1. Zed (10)", "1. zed (10)", "3. ann (0)"},
	},
	{
		[]string{"neg", "pos", "zero"},
		[]int{-5, 5, 0},
		[]string{"pos", "zero", "neg"},
		[]string{"1. pos (5)", "2. zero (0)", "3. neg (-5)"},
	},
	{
		[]string{"solo"},
		[]int{42},
		[]string{"solo"},
		[]string{"1. solo (42)"},
	},
	{nil, nil, []string{}, []string{}},
}

func TestRank(t *testing.T) {
	for _, tt := range rankTests {
		names, scores := slices.Clone(tt.names), slices.Clone(tt.scores)
		got, err := Rank(names, scores)
		if err != nil {
			t.Errorf("Rank(%q, %v) returned error %v, want nil", tt.names, tt.scores, err)
			continue
		}
		if !slices.Equal(got, tt.wantRank) {
			t.Errorf("Rank(%q, %v) = %q, want %q", tt.names, tt.scores, got, tt.wantRank)
		}
		if !slices.Equal(names, tt.names) || !slices.Equal(scores, tt.scores) {
			t.Errorf("Rank modified its inputs: %q %v became %q %v", tt.names, tt.scores, names, scores)
		}
	}
}

func TestStandings(t *testing.T) {
	for _, tt := range rankTests {
		names, scores := slices.Clone(tt.names), slices.Clone(tt.scores)
		got, err := Standings(names, scores)
		if err != nil {
			t.Errorf("Standings(%q, %v) returned error %v, want nil", tt.names, tt.scores, err)
			continue
		}
		if !slices.Equal(got, tt.wantStandings) {
			t.Errorf("Standings(%q, %v) = %q, want %q", tt.names, tt.scores, got, tt.wantStandings)
		}
		if !slices.Equal(names, tt.names) || !slices.Equal(scores, tt.scores) {
			t.Errorf("Standings modified its inputs: %q %v became %q %v", tt.names, tt.scores, names, scores)
		}
	}
}

func TestLengthMismatch(t *testing.T) {
	tests := []struct {
		names  []string
		scores []int
	}{
		{[]string{"x"}, []int{1, 2}},
		{[]string{"x", "y"}, []int{1}},
		{nil, []int{1}},
		{[]string{"x"}, nil},
	}
	for _, tt := range tests {
		if got, err := Rank(tt.names, tt.scores); err == nil || got != nil {
			t.Errorf("Rank(%q, %v) = %q, %v; want nil, <error>", tt.names, tt.scores, got, err)
		}
		if got, err := Standings(tt.names, tt.scores); err == nil || got != nil {
			t.Errorf("Standings(%q, %v) = %q, %v; want nil, <error>", tt.names, tt.scores, got, err)
		}
	}
}

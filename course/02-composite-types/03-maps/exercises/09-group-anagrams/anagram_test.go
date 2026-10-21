package anagram

import (
	"reflect"
	"slices"
	"testing"
)

func TestKey(t *testing.T) {
	tests := []struct {
		in, want string
	}{
		{"listen", "eilnst"},
		{"Silent", "eilnst"},
		{"ENLIST", "eilnst"},
		{"día", "adí"},
		{"cba", "abc"},
		{"aAb", "aab"},
		{"x", "x"},
		{"", ""},
	}
	for _, tt := range tests {
		if got := Key(tt.in); got != tt.want {
			t.Errorf("Key(%q) = %q, want %q", tt.in, got, tt.want)
		}
	}
}

func TestGroupAnagrams(t *testing.T) {
	tests := []struct {
		in   []string
		want [][]string
	}{
		{
			[]string{"eat", "tea", "tan", "ate", "nat", "bat"},
			[][]string{{"ate", "eat", "tea"}, {"bat"}, {"nat", "tan"}},
		},
		{
			[]string{"Listen", "silent", "enlist", "google"},
			[][]string{{"Listen", "enlist", "silent"}, {"google"}},
		},
		{
			[]string{"zz", "yy", "xx"},
			[][]string{{"xx"}, {"yy"}, {"zz"}},
		},
		{
			[]string{"stop", "pots", "tops", "opts", "post", "spot", "go"},
			[][]string{{"go"}, {"opts", "post", "pots", "spot", "stop", "tops"}},
		},
		{
			[]string{"día", "aíd", "dia"},
			[][]string{{"aíd", "día"}, {"dia"}},
		},
		{[]string{"a", "a"}, [][]string{{"a", "a"}}},
		{[]string{"", "b", ""}, [][]string{{"", ""}, {"b"}}},
		{[]string{"solo"}, [][]string{{"solo"}}},
	}
	for _, tt := range tests {
		in := slices.Clone(tt.in)
		// Run several times: depending on map order would only fail some of the time.
		for range 5 {
			got := GroupAnagrams(in)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GroupAnagrams(%q) = %q, want %q", tt.in, got, tt.want)
				break
			}
		}
		if !slices.Equal(in, tt.in) {
			t.Errorf("GroupAnagrams modified its input: %q became %q", tt.in, in)
		}
	}
}

func TestGroupAnagramsEmpty(t *testing.T) {
	for _, in := range [][]string{nil, {}} {
		if got := GroupAnagrams(in); len(got) != 0 {
			t.Errorf("GroupAnagrams(%q) = %q, want an empty result", in, got)
		}
	}
}

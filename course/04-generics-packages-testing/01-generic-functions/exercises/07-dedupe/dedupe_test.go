package dedupe_test

import (
	"slices"
	"strings"
	"testing"

	dd "learngo/course/04-generics-packages-testing/01-generic-functions/exercises/07-dedupe"
)

// User has a slice field, so User itself is NOT comparable.
type User struct {
	Name  string
	Email string
	Tags  []string
}

func names(us []User) []string {
	out := make([]string, 0, len(us))
	for _, u := range us {
		out = append(out, u.Name)
	}
	return out
}

func TestDedupeInts(t *testing.T) {
	tests := []struct {
		in, want []int
	}{
		{[]int{3, 1, 3, 2, 1}, []int{3, 1, 2}},
		{[]int{1, 2, 3}, []int{1, 2, 3}},
		{[]int{7, 7, 7}, []int{7}},
		{[]int{0, -1, 0, -1}, []int{0, -1}},
		{[]int{5}, []int{5}},
		{[]int{}, []int{}},
		{nil, []int{}},
	}
	for _, tt := range tests {
		if got := dd.Dedupe(tt.in); !slices.Equal(got, tt.want) {
			t.Errorf("Dedupe(%v) = %v, want %v", tt.in, got, tt.want)
		}
	}
}

func TestDedupeStrings(t *testing.T) {
	tests := []struct {
		in, want []string
	}{
		{[]string{"b", "a", "b"}, []string{"b", "a"}},
		{[]string{"", "x", ""}, []string{"", "x"}},
		{[]string{"Go", "go"}, []string{"Go", "go"}},
		{[]string{"日本", "日本", "中国"}, []string{"日本", "中国"}},
	}
	for _, tt := range tests {
		if got := dd.Dedupe(tt.in); !slices.Equal(got, tt.want) {
			t.Errorf("Dedupe(%q) = %q, want %q", tt.in, got, tt.want)
		}
	}
}

func TestDedupeDoesNotModify(t *testing.T) {
	in := []int{3, 1, 3, 2, 1}
	dd.Dedupe(in)
	if want := []int{3, 1, 3, 2, 1}; !slices.Equal(in, want) {
		t.Errorf("Dedupe(%v) modified the input: %v, want %v", want, in, want)
	}
}

func TestDedupeBy(t *testing.T) {
	users := []User{
		{"ann", "ann@x.io", []string{"admin"}},
		{"bob", "bob@x.io", nil},
		{"ann2", "ann@x.io", nil},
		{"cy", "cy@x.io", []string{"dev"}},
		{"bob2", "bob@x.io", nil},
	}
	got := dd.DedupeBy(users, func(u User) string { return u.Email })
	if want := []string{"ann", "bob", "cy"}; !slices.Equal(names(got), want) {
		t.Errorf("DedupeBy(users, email) = %q, want %q", names(got), want)
	}

	words := []string{"Go", "go", "GO", "rust", "Rust"}
	if got, want := dd.DedupeBy(words, strings.ToLower), []string{"Go", "rust"}; !slices.Equal(got, want) {
		t.Errorf("DedupeBy(%q, strings.ToLower) = %q, want %q", words, got, want)
	}

	byLen := dd.DedupeBy([]string{"a", "bb", "c", "dd", "eee"}, func(s string) int { return len(s) })
	if want := []string{"a", "bb", "eee"}; !slices.Equal(byLen, want) {
		t.Errorf("DedupeBy(%q, len) = %q, want %q", []string{"a", "bb", "c", "dd", "eee"}, byLen, want)
	}

	parity := dd.DedupeBy([]int{4, 7, 2, 9, 1}, func(n int) bool { return n%2 == 0 })
	if want := []int{4, 7}; !slices.Equal(parity, want) {
		t.Errorf("DedupeBy(%v, isEven) = %v, want %v", []int{4, 7, 2, 9, 1}, parity, want)
	}

	if got := dd.DedupeBy([]User(nil), func(u User) string { return u.Email }); len(got) != 0 {
		t.Errorf("DedupeBy(nil, email) = %v, want []", got)
	}
}

func TestDedupeByCallsKeyOnce(t *testing.T) {
	in := []int{1, 1, 2, 3, 3}
	calls := 0
	dd.DedupeBy(in, func(n int) int { calls++; return n })
	if calls != len(in) {
		t.Errorf("DedupeBy(%v, key) called key %d times, want %d", in, calls, len(in))
	}
}

func TestIntersect(t *testing.T) {
	tests := []struct {
		a, b, want []int
	}{
		{[]int{1, 2, 3, 2}, []int{2, 3, 4}, []int{2, 3}},
		{[]int{3, 2, 1}, []int{1, 2, 3}, []int{3, 2, 1}},
		{[]int{1, 1, 1}, []int{1}, []int{1}},
		{[]int{1, 2}, []int{3, 4}, []int{}},
		{[]int{1}, nil, []int{}},
		{nil, []int{1}, []int{}},
		{nil, nil, []int{}},
	}
	for _, tt := range tests {
		if got := dd.Intersect(tt.a, tt.b); !slices.Equal(got, tt.want) {
			t.Errorf("Intersect(%v, %v) = %v, want %v", tt.a, tt.b, got, tt.want)
		}
	}

	a, b := []string{"go", "zig", "c", "go"}, []string{"c", "go", "rust"}
	if got, want := dd.Intersect(a, b), []string{"go", "c"}; !slices.Equal(got, want) {
		t.Errorf("Intersect(%q, %q) = %q, want %q", a, b, got, want)
	}
}

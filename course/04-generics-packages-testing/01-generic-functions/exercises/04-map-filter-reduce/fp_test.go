package fp_test

import (
	"slices"
	"strconv"
	"strings"
	"testing"

	fp "learngo/course/04-generics-packages-testing/01-generic-functions/exercises/04-map-filter-reduce"
)

type User struct {
	Name   string
	Age    int
	Active bool
}

func isEven(n int) bool { return n%2 == 0 }

func TestMap(t *testing.T) {
	ints := []int{1, 2, 3}
	if got, want := fp.Map(ints, strconv.Itoa), []string{"1", "2", "3"}; !slices.Equal(got, want) {
		t.Errorf("Map(%v, strconv.Itoa) = %q, want %q", ints, got, want)
	}

	words := []string{"go", "gopher", "", "日本"}
	lengths := fp.Map(words, func(s string) int { return len(s) })
	if want := []int{2, 6, 0, 6}; !slices.Equal(lengths, want) {
		t.Errorf("Map(%q, len) = %v, want %v", words, lengths, want)
	}

	if got, want := fp.Map(words, strings.ToUpper), []string{"GO", "GOPHER", "", "日本"}; !slices.Equal(got, want) {
		t.Errorf("Map(%q, strings.ToUpper) = %q, want %q", words, got, want)
	}

	users := []User{{"ann", 30, true}, {"bob", 17, false}}
	names := fp.Map(users, func(u User) string { return u.Name })
	if want := []string{"ann", "bob"}; !slices.Equal(names, want) {
		t.Errorf("Map(users, name) = %q, want %q", names, want)
	}

	if got := fp.Map([]int{}, strconv.Itoa); len(got) != 0 {
		t.Errorf("Map([], strconv.Itoa) = %q, want []", got)
	}
	if got := fp.Map(nil, strconv.Itoa); len(got) != 0 {
		t.Errorf("Map(nil, strconv.Itoa) = %q, want []", got)
	}
}

func TestFilter(t *testing.T) {
	tests := []struct {
		in, want []int
	}{
		{[]int{1, 2, 3, 4}, []int{2, 4}},
		{[]int{4, 3, 2, 1}, []int{4, 2}},
		{[]int{2, 2, 2}, []int{2, 2, 2}},
		{[]int{-2, -1, 0}, []int{-2, 0}},
		{[]int{1, 3, 5}, []int{}},
		{[]int{}, []int{}},
		{nil, []int{}},
	}
	for _, tt := range tests {
		if got := fp.Filter(tt.in, isEven); !slices.Equal(got, tt.want) {
			t.Errorf("Filter(%v, isEven) = %v, want %v", tt.in, got, tt.want)
		}
	}

	users := []User{{"ann", 30, true}, {"bob", 17, false}, {"cy", 45, true}}
	active := fp.Filter(users, func(u User) bool { return u.Active })
	if want := []User{{"ann", 30, true}, {"cy", 45, true}}; !slices.Equal(active, want) {
		t.Errorf("Filter(users, active) = %v, want %v", active, want)
	}
}

func TestFilterDoesNotModify(t *testing.T) {
	in := []int{1, 2, 3, 4, 5, 6}
	got := fp.Filter(in, isEven)
	if want := []int{1, 2, 3, 4, 5, 6}; !slices.Equal(in, want) {
		t.Errorf("Filter(%v, isEven) modified the input: %v, want %v", want, in, want)
	}
	if len(got) > 0 && &got[0] == &in[0] {
		t.Errorf("Filter(%v, isEven) shares memory with the input, want a new slice", in)
	}
}

func TestReduce(t *testing.T) {
	add := func(a, b int) int { return a + b }
	tests := []struct {
		in   []int
		init int
		want int
	}{
		{[]int{1, 2, 3}, 0, 6},
		{[]int{1, 2, 3}, 10, 16},
		{[]int{-5, 5}, 0, 0},
		{[]int{}, 7, 7},
		{nil, 0, 0},
	}
	for _, tt := range tests {
		if got := fp.Reduce(tt.in, tt.init, add); got != tt.want {
			t.Errorf("Reduce(%v, %d, add) = %d, want %d", tt.in, tt.init, got, tt.want)
		}
	}

	words := []string{"a", "bb", "ccc"}
	total := fp.Reduce(words, 0, func(n int, s string) int { return n + len(s) })
	if total != 6 {
		t.Errorf("Reduce(%q, 0, sumLen) = %d, want %d", words, total, 6)
	}

	// Left to right: the order of f(acc, x) calls matters for non-commutative f.
	joined := fp.Reduce(words, ">", func(acc, s string) string { return acc + s })
	if joined != ">abbccc" {
		t.Errorf("Reduce(%q, %q, concat) = %q, want %q", words, ">", joined, ">abbccc")
	}

	digits := []int{1, 2, 3}
	n := fp.Reduce(digits, 0, func(acc, d int) int { return acc*10 + d })
	if n != 123 {
		t.Errorf("Reduce(%v, 0, acc*10+d) = %d, want %d", digits, n, 123)
	}
}

func TestFind(t *testing.T) {
	tests := []struct {
		in     []int
		want   int
		wantOK bool
	}{
		{[]int{1, 4, 6}, 4, true},
		{[]int{0, 1}, 0, true},
		{[]int{1, 3}, 0, false},
		{nil, 0, false},
	}
	for _, tt := range tests {
		got, ok := fp.Find(tt.in, isEven)
		if got != tt.want || ok != tt.wantOK {
			t.Errorf("Find(%v, isEven) = %d, %t, want %d, %t", tt.in, got, ok, tt.want, tt.wantOK)
		}
	}

	users := []User{{"ann", 30, true}, {"bob", 17, false}}
	u, ok := fp.Find(users, func(u User) bool { return u.Age < 18 })
	if !ok || u.Name != "bob" {
		t.Errorf("Find(users, minor) = %v, %t, want %v, %t", u, ok, users[1], true)
	}
	u, ok = fp.Find(users, func(u User) bool { return u.Age > 100 })
	if ok || u != (User{}) {
		t.Errorf("Find(users, age>100) = %v, %t, want %v, %t", u, ok, User{}, false)
	}
}

func TestCallsOncePerElement(t *testing.T) {
	in := []int{1, 2, 3, 4}

	calls := 0
	fp.Map(in, func(n int) int { calls++; return n })
	if calls != 4 {
		t.Errorf("Map(%v, f) called f %d times, want %d", in, calls, 4)
	}

	calls = 0
	fp.Filter(in, func(n int) bool { calls++; return true })
	if calls != 4 {
		t.Errorf("Filter(%v, f) called f %d times, want %d", in, calls, 4)
	}

	calls = 0
	fp.Find(in, func(n int) bool { calls++; return n == 2 })
	if calls != 2 {
		t.Errorf("Find(%v, ==2) called pred %d times, want %d (stop at the first match)", in, calls, 2)
	}
}

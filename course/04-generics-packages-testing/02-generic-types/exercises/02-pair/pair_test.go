package pair_test

import (
	"fmt"
	"slices"
	"strings"
	"testing"

	pair "learngo/course/04-generics-packages-testing/02-generic-types/exercises/02-pair"
)

// show formats a pair without calling its String method.
func show[A, B any](p pair.Pair[A, B]) string {
	return fmt.Sprintf("{%v %v}", p.First, p.Second)
}

// showAll formats a slice of pairs without calling String.
func showAll[A, B any](ps []pair.Pair[A, B]) string {
	parts := make([]string, len(ps))
	for i, p := range ps {
		parts[i] = show(p)
	}
	return "[" + strings.Join(parts, " ") + "]"
}

func TestNew(t *testing.T) {
	p := pair.New("age", 42)
	if want := (pair.Pair[string, int]{First: "age", Second: 42}); p != want {
		t.Errorf("New(%q, %d) = %s, want %s", "age", 42, show(p), show(want))
	}
	var np *int
	q := pair.New(np, []string(nil))
	if q.First != nil || q.Second != nil {
		t.Errorf("New(nil, nil) = %s, want {<nil> []}", show(q))
	}
}

func TestSwap(t *testing.T) {
	p := pair.Pair[string, int]{First: "age", Second: 42}
	got := p.Swap()
	if want := (pair.Pair[int, string]{First: 42, Second: "age"}); got != want {
		t.Errorf("%s.Swap() = %s, want %s", show(p), show(got), show(want))
	}
	if back := got.Swap(); back != p {
		t.Errorf("Swap().Swap() = %s, want %s", show(back), show(p))
	}
	same := pair.Pair[int, int]{First: 1, Second: 2}
	if got, want := same.Swap(), (pair.Pair[int, int]{First: 2, Second: 1}); got != want {
		t.Errorf("%s.Swap() = %s, want %s", show(same), show(got), show(want))
	}
}

func TestString(t *testing.T) {
	tests := []struct {
		name string
		got  string
		want string
	}{
		{`New("age", 42)`, pair.New("age", 42).String(), "(age, 42)"},
		{`New(1.5, true)`, pair.New(1.5, true).String(), "(1.5, true)"},
		{`New("", 0)`, pair.New("", 0).String(), "(, 0)"},
		{`New("héllo", "世界")`, pair.New("héllo", "世界").String(), "(héllo, 世界)"},
		{`New([]int{1, 2}, -3)`, pair.New([]int{1, 2}, -3).String(), "([1 2], -3)"},
	}
	for _, tt := range tests {
		if tt.got != tt.want {
			t.Errorf("%s.String() = %q, want %q", tt.name, tt.got, tt.want)
		}
	}
}

func TestZip(t *testing.T) {
	tests := []struct {
		as   []int
		bs   []string
		want []pair.Pair[int, string]
	}{
		{[]int{1, 2, 3}, []string{"a", "b", "c"}, []pair.Pair[int, string]{{1, "a"}, {2, "b"}, {3, "c"}}},
		{[]int{1, 2, 3}, []string{"a", "b"}, []pair.Pair[int, string]{{1, "a"}, {2, "b"}}},
		{[]int{1}, []string{"a", "b", "c"}, []pair.Pair[int, string]{{1, "a"}}},
		{[]int{1}, nil, nil},
		{nil, []string{"a"}, nil},
		{nil, nil, nil},
	}
	for _, tt := range tests {
		asBefore, bsBefore := slices.Clone(tt.as), slices.Clone(tt.bs)
		got := pair.Zip(tt.as, tt.bs)
		if len(got) != len(tt.want) || (len(got) > 0 && !slices.Equal(got, tt.want)) {
			t.Errorf("Zip(%v, %q) = %s, want %s", tt.as, tt.bs, showAll(got), showAll(tt.want))
		}
		if !slices.Equal(tt.as, asBefore) || !slices.Equal(tt.bs, bsBefore) {
			t.Errorf("Zip(%v, %q) modified its inputs", asBefore, bsBefore)
		}
	}
}

func TestUnzip(t *testing.T) {
	ps := []pair.Pair[int, string]{{1, "a"}, {2, "b"}, {-3, "ü"}}
	as, bs := pair.Unzip(ps)
	if want := []int{1, 2, -3}; !slices.Equal(as, want) {
		t.Errorf("Unzip(%s) firsts = %v, want %v", showAll(ps), as, want)
	}
	if want := []string{"a", "b", "ü"}; !slices.Equal(bs, want) {
		t.Errorf("Unzip(%s) seconds = %q, want %q", showAll(ps), bs, want)
	}
	as, bs = pair.Unzip[int, string](nil)
	if len(as) != 0 || len(bs) != 0 {
		t.Errorf("Unzip(nil) = (%v, %q), want ([], [])", as, bs)
	}
}

func TestZipUnzipRoundTrip(t *testing.T) {
	names := []string{"ann", "bob"}
	ages := []float64{31.5, 28}
	gotNames, gotAges := pair.Unzip(pair.Zip(names, ages))
	if !slices.Equal(gotNames, names) || !slices.Equal(gotAges, ages) {
		t.Errorf("Unzip(Zip(%q, %v)) = (%q, %v), want (%q, %v)", names, ages, gotNames, gotAges, names, ages)
	}
}

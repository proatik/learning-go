package coll_test

import (
	"slices"
	"strconv"
	"strings"
	"testing"

	coll "learngo/course/04-generics-packages-testing/02-generic-types/exercises/07-generic-methods"
)

func isEven(n int) bool { return n%2 == 0 }

func TestOfItemsLen(t *testing.T) {
	src := []int{3, 1, 2}
	l := coll.Of(src...)
	src[0] = 99
	if got, want := l.Items(), []int{3, 1, 2}; !slices.Equal(got, want) {
		t.Errorf("Of(3, 1, 2).Items() = %v, want %v (Of must copy its arguments)", got, want)
	}
	if got := l.Len(); got != 3 {
		t.Errorf("Of(3, 1, 2).Len() = %d, want 3", got)
	}
	items := l.Items()
	if len(items) == 3 {
		items[0] = -1
	}
	if got := l.Items(); len(got) == 0 || got[0] != 3 {
		t.Errorf("modifying the Items() result changed the list: Items() = %v, want [3 1 2]", got)
	}
	var zero coll.List[string]
	if got := zero.Len(); got != 0 {
		t.Errorf("zero List Len() = %d, want 0", got)
	}
}

func TestFilter(t *testing.T) {
	tests := []struct {
		in, want []int
	}{
		{[]int{1, 2, 3, 4}, []int{2, 4}},
		{[]int{1, 3}, []int{}},
		{[]int{}, []int{}},
		{[]int{0, -2, 7}, []int{0, -2}},
	}
	for _, tt := range tests {
		l := coll.Of(tt.in...)
		if got := l.Filter(isEven).Items(); !slices.Equal(got, tt.want) {
			t.Errorf("Of(%v).Filter(isEven) = %v, want %v", tt.in, got, tt.want)
		}
		if got := l.Items(); !slices.Equal(got, tt.in) {
			t.Errorf("Of(%v).Filter changed the receiver to %v", tt.in, got)
		}
	}
}

func TestMap(t *testing.T) {
	strs := coll.Of(1, -2, 30).Map(strconv.Itoa)
	if got, want := strs.Items(), []string{"1", "-2", "30"}; !slices.Equal(got, want) {
		t.Errorf("Of(1, -2, 30).Map(strconv.Itoa) = %q, want %q", got, want)
	}
	half := coll.Of(1, 2).Map[float64](func(n int) float64 { return float64(n) / 2 })
	if got, want := half.Items(), []float64{0.5, 1}; !slices.Equal(got, want) {
		t.Errorf("Of(1, 2).Map[float64](half) = %v, want %v", got, want)
	}
	lens := coll.Of("héllo", "", "世界").Map(func(s string) int { return len([]rune(s)) })
	if got, want := lens.Items(), []int{5, 0, 2}; !slices.Equal(got, want) {
		t.Errorf("Of(\"héllo\", \"\", \"世界\").Map(runeCount) = %v, want %v", got, want)
	}
	if got := coll.Of[int]().Map(strconv.Itoa).Len(); got != 0 {
		t.Errorf("Of().Map(strconv.Itoa).Len() = %d, want 0", got)
	}
}

func TestFlatMap(t *testing.T) {
	split := func(s string) []string { return strings.Split(s, "") }
	got := coll.Of("ab", "", "c").FlatMap(split).Items()
	if want := []string{"a", "b", "c"}; !slices.Equal(got, want) {
		t.Errorf("Of(\"ab\", \"\", \"c\").FlatMap(split) = %q, want %q", got, want)
	}
	repeat := func(n int) []int { return slices.Repeat([]int{n}, n) }
	got2 := coll.Of(1, 0, 3).FlatMap(repeat).Items()
	if want := []int{1, 3, 3, 3}; !slices.Equal(got2, want) {
		t.Errorf("Of(1, 0, 3).FlatMap(repeat) = %v, want %v", got2, want)
	}
}

func TestReduce(t *testing.T) {
	join := func(acc string, n int) string { return acc + strconv.Itoa(n) }
	if got := coll.Of(1, 2, 3).Reduce("", join); got != "123" {
		t.Errorf("Of(1, 2, 3).Reduce(\"\", join) = %q, want %q", got, "123")
	}
	if got := coll.Of(1, 2, 3).Reduce(">", join); got != ">123" {
		t.Errorf("Of(1, 2, 3).Reduce(\">\", join) = %q, want %q", got, ">123")
	}
	if got := coll.Of[int]().Reduce("init", join); got != "init" {
		t.Errorf("Of().Reduce(\"init\", join) = %q, want %q", got, "init")
	}
	sub := func(acc, n int) int { return acc - n }
	if got := coll.Of(1, 2, 3).Reduce(10, sub); got != 4 {
		t.Errorf("Of(1, 2, 3).Reduce(10, sub) = %d, want 4", got)
	}
}

func TestGroupBy(t *testing.T) {
	words := coll.Of("go", "c", "zig", "rust", "java", "d")
	groups := words.GroupBy(func(s string) int { return len(s) })
	want := map[int][]string{1: {"c", "d"}, 2: {"go"}, 3: {"zig"}, 4: {"rust", "java"}}
	if len(groups) != len(want) {
		t.Errorf("GroupBy(len) gave %d groups, want %d", len(groups), len(want))
	}
	for k, w := range want {
		if got := groups[k].Items(); !slices.Equal(got, w) {
			t.Errorf("GroupBy(len)[%d] = %q, want %q", k, got, w)
		}
	}
	parity := coll.Of(1, 2, 3, 4, 5).GroupBy(isEven)
	if got := parity[true].Items(); !slices.Equal(got, []int{2, 4}) {
		t.Errorf("GroupBy(isEven)[true] = %v, want [2 4]", got)
	}
	if got := parity[false].Items(); !slices.Equal(got, []int{1, 3, 5}) {
		t.Errorf("GroupBy(isEven)[false] = %v, want [1 3 5]", got)
	}
	if got := coll.Of[int]().GroupBy(isEven); len(got) != 0 {
		t.Errorf("Of().GroupBy(isEven) has %d groups, want 0", len(got))
	}
}

type cents int64

func TestSum(t *testing.T) {
	if got := coll.Sum(coll.Of(1.5, 2.5)); got != 4 {
		t.Errorf("Sum(Of(1.5, 2.5)) = %v, want 4", got)
	}
	if got := coll.Sum(coll.Of(1, -2, 10)); got != 9 {
		t.Errorf("Sum(Of(1, -2, 10)) = %d, want 9", got)
	}
	if got := coll.Sum(coll.Of[cents](199, 1)); got != 200 {
		t.Errorf("Sum(Of[cents](199, 1)) = %d, want 200", got)
	}
	if got := coll.Sum(coll.Of[int]()); got != 0 {
		t.Errorf("Sum(Of()) = %d, want 0", got)
	}
}

func TestChaining(t *testing.T) {
	got := coll.Of(1, 2, 3, 4).Filter(isEven).Map(strconv.Itoa).Items()
	if want := []string{"2", "4"}; !slices.Equal(got, want) {
		t.Errorf("Of(1, 2, 3, 4).Filter(isEven).Map(strconv.Itoa) = %q, want %q", got, want)
	}
	total := coll.Sum(coll.Of("a", "bb", "ccc").Map(func(s string) int { return len(s) }))
	if total != 6 {
		t.Errorf("Sum(Of(\"a\", \"bb\", \"ccc\").Map(len)) = %d, want 6", total)
	}
}

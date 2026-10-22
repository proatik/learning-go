package result

import (
	"errors"
	"fmt"
	"slices"
	"strconv"
	"testing"
)

var errBoom = errors.New("boom")

func TestOkFailOf(t *testing.T) {
	if v, err := Ok(42).Get(); v != 42 || err != nil {
		t.Errorf("Ok(42).Get() = %d, %v, want 42, <nil>", v, err)
	}
	if !Ok("x").IsOk() {
		t.Errorf("Ok(%q).IsOk() = false, want true", "x")
	}
	if v, err := Fail[int](errBoom).Get(); v != 0 || !errors.Is(err, errBoom) {
		t.Errorf("Fail[int](errBoom).Get() = %d, %v, want 0, %v", v, err, errBoom)
	}
	if Fail[string](errBoom).IsOk() {
		t.Errorf("Fail[string](errBoom).IsOk() = true, want false")
	}
	if v, err := Of(7, nil).Get(); v != 7 || err != nil {
		t.Errorf("Of(7, nil).Get() = %d, %v, want 7, <nil>", v, err)
	}
	if v, err := Of(7, errBoom).Get(); v != 0 || !errors.Is(err, errBoom) {
		t.Errorf("Of(7, errBoom).Get() = %d, %v, want 0, %v", v, err, errBoom)
	}
	if v, err := Of(strconv.Atoi("12")).Get(); v != 12 || err != nil {
		t.Errorf("Of(strconv.Atoi(%q)).Get() = %d, %v, want 12, <nil>", "12", v, err)
	}
}

func TestOrElse(t *testing.T) {
	if got := Ok(3).OrElse(-1); got != 3 {
		t.Errorf("Ok(3).OrElse(-1) = %d, want 3", got)
	}
	if got := Fail[int](errBoom).OrElse(-1); got != -1 {
		t.Errorf("Fail(errBoom).OrElse(-1) = %d, want -1", got)
	}
	if got := Ok(0).OrElse(5); got != 0 {
		t.Errorf("Ok(0).OrElse(5) = %d, want 0", got)
	}
}

func TestContext(t *testing.T) {
	r := Fail[int](errBoom).Context("load config").Context("start")
	_, err := r.Get()
	if want := "start: load config: boom"; err == nil || err.Error() != want {
		t.Errorf("Fail(boom).Context(%q).Context(%q) error = %v, want %q", "load config", "start", err, want)
	}
	if !errors.Is(err, errBoom) {
		t.Errorf("errors.Is(%v, errBoom) = false, want true (wrap with %%w)", err)
	}
	if v, err := Ok(1).Context("never").Get(); v != 1 || err != nil {
		t.Errorf("Ok(1).Context(%q).Get() = %d, %v, want 1, <nil>", "never", v, err)
	}
}

func TestMap(t *testing.T) {
	s := Ok(21).Map(func(n int) int { return n * 2 }).Map(strconv.Itoa)
	if v, err := s.Get(); v != "42" || err != nil {
		t.Errorf("Ok(21).Map(double).Map(Itoa).Get() = %q, %v, want %q, <nil>", v, err, "42")
	}
	called := false
	f := Fail[int](errBoom).Map(func(n int) string { called = true; return "x" })
	if v, err := f.Get(); v != "" || !errors.Is(err, errBoom) {
		t.Errorf("Fail(boom).Map(f).Get() = %q, %v, want %q, %v", v, err, "", errBoom)
	}
	if called {
		t.Errorf("Fail(boom).Map(f) called f, want f skipped")
	}
}

func TestThen(t *testing.T) {
	parse := func(s string) (int, error) { return strconv.Atoi(s) }
	positive := func(n int) (int, error) {
		if n <= 0 {
			return 0, fmt.Errorf("%d is not positive", n)
		}
		return n, nil
	}
	tests := []struct {
		in      string
		want    int
		wantErr bool
	}{
		{"5", 5, false},
		{"-3", 0, true},
		{"0", 0, true},
		{"abc", 0, true},
	}
	for _, tt := range tests {
		got, err := Ok(tt.in).Then(parse).Then(positive).Get()
		if got != tt.want || (err != nil) != tt.wantErr {
			t.Errorf("Ok(%q).Then(parse).Then(positive).Get() = %d, %v, want %d, error %v", tt.in, got, err, tt.want, tt.wantErr)
		}
	}

	called := false
	r := Fail[string](errBoom).Then(func(s string) (int, error) { called = true; return 1, nil })
	if _, err := r.Get(); !errors.Is(err, errBoom) || called {
		t.Errorf("Fail(boom).Then(f) error = %v, f called = %v; want %v, false", err, called, errBoom)
	}
}

func TestTry(t *testing.T) {
	if v, err := Try(func() int { return 9 }).Get(); v != 9 || err != nil {
		t.Errorf("Try(returns 9).Get() = %d, %v, want 9, <nil>", v, err)
	}
	r := Try(func() int {
		var s []int
		return s[3]
	})
	_, err := r.Get()
	if !errors.Is(err, ErrPanic) {
		t.Errorf("Try(index out of range) error = %v, want %v", err, ErrPanic)
	}
	r2 := Try(func() string { panic("disk on fire") })
	if _, err := r2.Get(); err == nil || err.Error() != "result: panic: disk on fire" {
		t.Errorf("Try(panic(%q)) error = %v, want %q", "disk on fire", err, "result: panic: disk on fire")
	}
}

func TestCollect(t *testing.T) {
	got, err := Collect(slices.Values([]Result[int]{Ok(1), Ok(2), Ok(3)}))
	if err != nil || !slices.Equal(got, []int{1, 2, 3}) {
		t.Errorf("Collect(Ok 1, 2, 3) = %v, %v, want [1 2 3], <nil>", got, err)
	}

	empty, err := Collect(slices.Values([]Result[int]{}))
	if err != nil || empty == nil || len(empty) != 0 {
		t.Errorf("Collect(empty) = %#v, %v, want []int{}, <nil>", empty, err)
	}

	pulled := 0
	seq := func(yield func(Result[int]) bool) {
		for _, r := range []Result[int]{Ok(1), Ok(2), Fail[int](errBoom), Ok(4), Ok(5)} {
			pulled++
			if !yield(r) {
				return
			}
		}
	}
	got, err = Collect(seq)
	if got != nil {
		t.Errorf("Collect(with a failure) values = %v, want nil", got)
	}
	if want := "item 2: boom"; err == nil || err.Error() != want || !errors.Is(err, errBoom) {
		t.Errorf("Collect(Ok, Ok, Fail(boom), ...) error = %v, want %q wrapping errBoom", err, want)
	}
	if pulled != 3 {
		t.Errorf("Collect pulled %d items, want 3 (stop at the first failure)", pulled)
	}
}

func Example() {
	parse := func(s string) (int, error) { return strconv.Atoi(s) }
	for _, in := range []string{"21", "x"} {
		r := Ok(in).Then(parse).Map(func(n int) int { return n * 2 }).Context("input " + in)
		v, err := r.Get()
		fmt.Println(v, err)
	}
	// Output:
	// 42 <nil>
	// 0 input x: strconv.Atoi: parsing "x": invalid syntax
}

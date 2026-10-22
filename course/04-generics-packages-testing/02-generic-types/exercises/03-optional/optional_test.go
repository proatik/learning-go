package optional

import "testing"

// state formats o from its fields, without calling its methods.
func state[T any](o Optional[T]) string {
	if !o.ok {
		return "None"
	}
	return "Some"
}

func TestSome(t *testing.T) {
	if v, ok := Some(5).Get(); v != 5 || !ok {
		t.Errorf("Some(5).Get() = (%d, %v), want (5, true)", v, ok)
	}
	if v, ok := Some(0).Get(); v != 0 || !ok {
		t.Errorf("Some(0).Get() = (%d, %v), want (0, true)", v, ok)
	}
	if !Some("").IsSome() {
		t.Errorf("Some(%q).IsSome() = false, want true", "")
	}
	var np *int
	if v, ok := Some(np).Get(); v != nil || !ok {
		t.Errorf("Some(nil).Get() = (%v, %v), want (<nil>, true)", v, ok)
	}
}

func TestNoneAndZeroValue(t *testing.T) {
	if v, ok := None[int]().Get(); v != 0 || ok {
		t.Errorf("None[int]().Get() = (%d, %v), want (0, false)", v, ok)
	}
	if None[string]().IsSome() {
		t.Errorf("None[string]().IsSome() = true, want false")
	}
	if got := state(None[float64]()); got != "None" {
		t.Errorf("None[float64]() = %s, want None", got)
	}
	var o Optional[string]
	if o.IsSome() {
		t.Errorf("zero Optional[string] IsSome() = true, want false")
	}
	if v, ok := o.Get(); v != "" || ok {
		t.Errorf("zero Optional[string] Get() = (%q, %v), want (\"\", false)", v, ok)
	}
}

func TestOrElse(t *testing.T) {
	tests := []struct {
		name string
		got  string
		want string
	}{
		{`Some("ann").OrElse("guest")`, Some("ann").OrElse("guest"), "ann"},
		{`Some("").OrElse("guest")`, Some("").OrElse("guest"), ""},
		{`None[string]().OrElse("guest")`, None[string]().OrElse("guest"), "guest"},
		{`Optional[string]{}.OrElse("世界")`, Optional[string]{}.OrElse("世界"), "世界"},
	}
	for _, tt := range tests {
		if tt.got != tt.want {
			t.Errorf("%s = %q, want %q", tt.name, tt.got, tt.want)
		}
	}
	if got := Some(0).OrElse(-1); got != 0 {
		t.Errorf("Some(0).OrElse(-1) = %d, want 0", got)
	}
}

func TestString(t *testing.T) {
	tests := []struct {
		name string
		got  string
		want string
	}{
		{"Some(42)", Some(42).String(), "Some(42)"},
		{"Some(-1.5)", Some(-1.5).String(), "Some(-1.5)"},
		{`Some("héllo")`, Some("héllo").String(), "Some(héllo)"},
		{"Some([]int{1, 2})", Some([]int{1, 2}).String(), "Some([1 2])"},
		{"None[int]()", None[int]().String(), "None"},
		{"Optional[bool]{}", Optional[bool]{}.String(), "None"},
	}
	for _, tt := range tests {
		if tt.got != tt.want {
			t.Errorf("%s.String() = %q, want %q", tt.name, tt.got, tt.want)
		}
	}
}

func TestFirst(t *testing.T) {
	if v, ok := First([]string{"a", "b"}).Get(); v != "a" || !ok {
		t.Errorf("First([a b]).Get() = (%q, %v), want (\"a\", true)", v, ok)
	}
	if v, ok := First([]int{0}).Get(); v != 0 || !ok {
		t.Errorf("First([0]).Get() = (%d, %v), want (0, true)", v, ok)
	}
	if got := state(First([]int{})); got != "None" {
		t.Errorf("First([]) = %s, want None", got)
	}
	if got := state(First[int](nil)); got != "None" {
		t.Errorf("First(nil) = %s, want None", got)
	}
}

func TestLookup(t *testing.T) {
	m := map[string]int{"x": 0, "y": 7}
	tests := []struct {
		key  string
		want int
		ok   bool
	}{
		{"x", 0, true},
		{"y", 7, true},
		{"z", 0, false},
		{"", 0, false},
	}
	for _, tt := range tests {
		v, ok := Lookup(m, tt.key).Get()
		if v != tt.want || ok != tt.ok {
			t.Errorf("Lookup(%v, %q).Get() = (%d, %v), want (%d, %v)", m, tt.key, v, ok, tt.want, tt.ok)
		}
	}
	var nilMap map[int]string
	if got := state(Lookup(nilMap, 1)); got != "None" {
		t.Errorf("Lookup(nil map, 1) = %s, want None", got)
	}
}

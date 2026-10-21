package must_test

import (
	"errors"
	"io"
	"strconv"
	"testing"

	must "learngo/course/03-methods-interfaces-errors/05-panic-recover/exercises/03-must-atoi"
)

// catch runs f and returns the value it panicked with (nil if it did not).
func catch(f func()) (r any) {
	defer func() { r = recover() }()
	f()
	return nil
}

func TestMust(t *testing.T) {
	for _, n := range []int{42, 0, -7} {
		var got int
		if r := catch(func() { got = must.Must(n, nil) }); r != nil {
			t.Errorf("Must(%d, nil) panicked: %v", n, r)
			continue
		}
		if got != n {
			t.Errorf("Must(%d, nil) = %d, want %d", n, got, n)
		}
	}
}

func TestMustPanicsWithSameError(t *testing.T) {
	custom := errors.New("custom")
	for _, err := range []error{io.EOF, custom} {
		r := catch(func() { must.Must(1, err) })
		if r == nil {
			t.Errorf("Must(1, %v) did not panic, want panic(%v)", err, err)
			continue
		}
		if got, ok := r.(error); !ok || got != err {
			t.Errorf("Must(1, %v) panicked with %#v, want the same error value %v", err, r, err)
		}
	}
}

func TestMustAtoi(t *testing.T) {
	tests := []struct {
		s    string
		want int
	}{
		{"0", 0},
		{"8080", 8080},
		{"-17", -17},
		{"+5", 5},
		{"007", 7},
		{"9223372036854775807", 9223372036854775807},
	}
	for _, tt := range tests {
		var got int
		if r := catch(func() { got = must.MustAtoi(tt.s) }); r != nil {
			t.Errorf("MustAtoi(%q) panicked: %v", tt.s, r)
			continue
		}
		if got != tt.want {
			t.Errorf("MustAtoi(%q) = %d, want %d", tt.s, got, tt.want)
		}
	}
}

func TestMustAtoiPanics(t *testing.T) {
	tests := []struct {
		s      string
		msg    string
		target error
	}{
		{"12a", `MustAtoi("12a"): strconv.Atoi: parsing "12a": invalid syntax`, strconv.ErrSyntax},
		{"", `MustAtoi(""): strconv.Atoi: parsing "": invalid syntax`, strconv.ErrSyntax},
		{" 1", `MustAtoi(" 1"): strconv.Atoi: parsing " 1": invalid syntax`, strconv.ErrSyntax},
		{"٣", `MustAtoi("٣"): strconv.Atoi: parsing "٣": invalid syntax`, strconv.ErrSyntax},
		{"99999999999999999999", `MustAtoi("99999999999999999999"): strconv.Atoi: parsing "99999999999999999999": value out of range`, strconv.ErrRange},
	}
	for _, tt := range tests {
		r := catch(func() { must.MustAtoi(tt.s) })
		if r == nil {
			t.Errorf("MustAtoi(%q) did not panic, want panic %q", tt.s, tt.msg)
			continue
		}
		err, ok := r.(error)
		if !ok {
			t.Errorf("MustAtoi(%q) panicked with %T %v, want an error value", tt.s, r, r)
			continue
		}
		if err.Error() != tt.msg {
			t.Errorf("MustAtoi(%q) panic = %q, want %q", tt.s, err.Error(), tt.msg)
		}
		if !errors.Is(err, tt.target) {
			t.Errorf("MustAtoi(%q) panic = %v, want it to wrap %v", tt.s, err, tt.target)
		}
	}
}

func TestMustPort(t *testing.T) {
	for _, tt := range []struct {
		s    string
		want int
	}{{"1", 1}, {"80", 80}, {"8080", 8080}, {"65535", 65535}} {
		var got int
		if r := catch(func() { got = must.MustPort(tt.s) }); r != nil {
			t.Errorf("MustPort(%q) panicked: %v", tt.s, r)
			continue
		}
		if got != tt.want {
			t.Errorf("MustPort(%q) = %d, want %d", tt.s, got, tt.want)
		}
	}
}

func TestMustPortPanics(t *testing.T) {
	tests := []struct {
		s      string
		msg    string
		target error
	}{
		{"0", `MustPort("0"): port out of range`, must.ErrPortRange},
		{"65536", `MustPort("65536"): port out of range`, must.ErrPortRange},
		{"-80", `MustPort("-80"): port out of range`, must.ErrPortRange},
		{"http", `MustAtoi("http"): strconv.Atoi: parsing "http": invalid syntax`, strconv.ErrSyntax},
	}
	for _, tt := range tests {
		r := catch(func() { must.MustPort(tt.s) })
		if r == nil {
			t.Errorf("MustPort(%q) did not panic, want panic %q", tt.s, tt.msg)
			continue
		}
		err, ok := r.(error)
		if !ok {
			t.Errorf("MustPort(%q) panicked with %T %v, want an error value", tt.s, r, r)
			continue
		}
		if err.Error() != tt.msg {
			t.Errorf("MustPort(%q) panic = %q, want %q", tt.s, err.Error(), tt.msg)
		}
		if !errors.Is(err, tt.target) {
			t.Errorf("MustPort(%q) panic = %v, want it to wrap %v", tt.s, err, tt.target)
		}
	}
}

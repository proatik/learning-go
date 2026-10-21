package ipaddr

import (
	"fmt"
	"testing"
)

// Compile-time checks: both types (and their pointers) are Stringers.
var (
	_ fmt.Stringer = IPAddr{}
	_ fmt.Stringer = &IPAddr{}
	_ fmt.Stringer = Color{}
	_ fmt.Stringer = &Color{}
)

func TestIPAddrString(t *testing.T) {
	tests := []struct {
		ip   IPAddr
		want string
	}{
		{IPAddr{127, 0, 0, 1}, "127.0.0.1"},
		{IPAddr{}, "0.0.0.0"},
		{IPAddr{8, 8, 8, 8}, "8.8.8.8"},
		{IPAddr{255, 255, 255, 255}, "255.255.255.255"},
		{IPAddr{192, 168, 10, 254}, "192.168.10.254"},
		{IPAddr{10, 0, 1, 9}, "10.0.1.9"},
	}
	for _, tt := range tests {
		if got := fmt.Sprint(tt.ip); got != tt.want {
			t.Errorf("fmt.Sprint(IPAddr%v) = %q, want %q", [4]byte(tt.ip), got, tt.want)
		}
		ip := tt.ip
		if got := fmt.Sprint(&ip); got != tt.want {
			t.Errorf("fmt.Sprint(&IPAddr%v) = %q, want %q", [4]byte(tt.ip), got, tt.want)
		}
	}
}

func TestIPAddrInSlice(t *testing.T) {
	in := []IPAddr{{8, 8, 8, 8}, {1, 1, 1, 1}}
	want := "[8.8.8.8 1.1.1.1]"
	if got := fmt.Sprintf("%v", in); got != want {
		t.Errorf("fmt.Sprintf(%%v, %s) = %q, want %q", "[]IPAddr{{8,8,8,8},{1,1,1,1}}", got, want)
	}
	m := map[string]IPAddr{"loopback": {127, 0, 0, 1}}
	want = "map[loopback:127.0.0.1]"
	if got := fmt.Sprint(m); got != want {
		t.Errorf("fmt.Sprint(%s) = %q, want %q", "map{loopback: {127,0,0,1}}", got, want)
	}
}

func TestColorString(t *testing.T) {
	tests := []struct {
		c    Color
		want string
	}{
		{Color{255, 128, 0}, "#ff8000"},
		{Color{}, "#000000"},
		{Color{255, 255, 255}, "#ffffff"},
		{Color{1, 2, 3}, "#010203"},
		{Color{R: 0xab, G: 0xcd, B: 0xef}, "#abcdef"},
		{Color{B: 16}, "#000010"},
	}
	for _, tt := range tests {
		if got := fmt.Sprint(tt.c); got != tt.want {
			t.Errorf("fmt.Sprint(Color{%d, %d, %d}) = %q, want %q", tt.c.R, tt.c.G, tt.c.B, got, tt.want)
		}
		c := tt.c
		if got := fmt.Sprintf("%s", &c); got != tt.want {
			t.Errorf("fmt.Sprintf(%%s, &Color{%d, %d, %d}) = %q, want %q", tt.c.R, tt.c.G, tt.c.B, got, tt.want)
		}
	}
}

package caesar

import "testing"

func TestEncrypt(t *testing.T) {
	tests := []struct {
		s     string
		shift int
		want  string
	}{
		{"abc", 1, "bcd"},
		{"xyz", 3, "abc"},
		{"XYZ", 3, "ABC"},
		{"Hello, World!", 13, "Uryyb, Jbeyq!"},
		{"abc", -1, "zab"},
		{"abc", -27, "zab"},
		{"abc", 27, "bcd"},
		{"abc", 26, "abc"},
		{"abc", 0, "abc"},
		{"abc", 52, "abc"},
		{"Café 2027", 1, "Dbgé 2027"},
		{"世界 go!", 2, "世界 iq!"},
		{"", 5, ""},
		{"The quick brown fox jumps over the lazy dog.", 3, "Wkh txlfn eurzq ira mxpsv ryhu wkh odcb grj."},
	}
	for _, tt := range tests {
		if got := Encrypt(tt.s, tt.shift); got != tt.want {
			t.Errorf("Encrypt(%q, %d) = %q, want %q", tt.s, tt.shift, got, tt.want)
		}
	}
}

func TestDecrypt(t *testing.T) {
	tests := []struct {
		s     string
		shift int
		want  string
	}{
		{"bcd", 1, "abc"},
		{"abc", 3, "xyz"},
		{"Uryyb, Jbeyq!", 13, "Hello, World!"},
		{"zab", -1, "abc"},
		{"Dbgé 2027", 1, "Café 2027"},
		{"", 3, ""},
	}
	for _, tt := range tests {
		if got := Decrypt(tt.s, tt.shift); got != tt.want {
			t.Errorf("Decrypt(%q, %d) = %q, want %q", tt.s, tt.shift, got, tt.want)
		}
	}
}

func TestRoundTrip(t *testing.T) {
	s := "Go 1.27: héllo, 世界! Zz"
	for shift := -60; shift <= 60; shift++ {
		if got := Decrypt(Encrypt(s, shift), shift); got != s {
			t.Errorf("Decrypt(Encrypt(%q, %d), %d) = %q, want %q", s, shift, shift, got, s)
		}
	}
}

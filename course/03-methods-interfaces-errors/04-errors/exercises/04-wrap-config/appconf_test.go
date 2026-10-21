package appconf

import (
	"errors"
	"strconv"
	"testing"
)

var files = map[string]string{
	"app.conf":    "localhost:8080",
	"db.conf":     "  db.local:5432\n",
	"edge.conf":   "h:1",
	"max.conf":    "h:65535",
	"empty.conf":  "",
	"noport.conf": "localhost",
	"nohost.conf": ":80",
	"blank.conf":  "h:",
	"alpha.conf":  "h:abc",
	"huge.conf":   "h:99999999999999999999",
	"zero.conf":   "h:0",
	"neg.conf":    "h:-1",
	"high.conf":   "h:70000",
	"ipv6.conf":   "::1:80",
}

func TestLoadConfigOK(t *testing.T) {
	tests := []struct {
		name string
		want Config
	}{
		{"app.conf", Config{"localhost", 8080}},
		{"db.conf", Config{"db.local", 5432}},
		{"edge.conf", Config{"h", 1}},
		{"max.conf", Config{"h", 65535}},
	}
	for _, tt := range tests {
		got, err := LoadConfig(files, tt.name)
		if err != nil {
			t.Errorf("LoadConfig(%q) error = %v, want <nil>", tt.name, err)
			continue
		}
		if got != tt.want {
			t.Errorf("LoadConfig(%q) = %+v, want %+v", tt.name, got, tt.want)
		}
	}
}

func TestLoadConfigSentinels(t *testing.T) {
	tests := []struct {
		name string
		want error
		msg  string
	}{
		{"missing.conf", ErrNotFound, `config "missing.conf": file not found`},
		{"noport.conf", ErrMissingPort, `config "noport.conf": missing port`},
		{"empty.conf", ErrMissingPort, `config "empty.conf": missing port`},
		{"nohost.conf", ErrMissingHost, `config "nohost.conf": missing host`},
		{"ipv6.conf", ErrMissingHost, `config "ipv6.conf": missing host`},
		{"zero.conf", ErrPortRange, `config "zero.conf": port 0: port out of range`},
		{"neg.conf", ErrPortRange, `config "neg.conf": port -1: port out of range`},
		{"high.conf", ErrPortRange, `config "high.conf": port 70000: port out of range`},
	}
	for _, tt := range tests {
		got, err := LoadConfig(files, tt.name)
		if !errors.Is(err, tt.want) {
			t.Errorf("LoadConfig(%q) error = %v, want %v", tt.name, err, tt.want)
			continue
		}
		if err.Error() != tt.msg {
			t.Errorf("LoadConfig(%q) error = %q, want %q", tt.name, err.Error(), tt.msg)
		}
		if got != (Config{}) {
			t.Errorf("LoadConfig(%q) = %+v, want the zero Config with an error", tt.name, got)
		}
	}
}

func TestLoadConfigNumError(t *testing.T) {
	tests := []struct {
		name  string
		cause error
		num   string
		msg   string
	}{
		{"alpha.conf", strconv.ErrSyntax, "abc", `config "alpha.conf": port "abc": strconv.Atoi: parsing "abc": invalid syntax`},
		{"blank.conf", strconv.ErrSyntax, "", `config "blank.conf": port "": strconv.Atoi: parsing "": invalid syntax`},
		{"huge.conf", strconv.ErrRange, "99999999999999999999", `config "huge.conf": port "99999999999999999999": strconv.Atoi: parsing "99999999999999999999": value out of range`},
	}
	for _, tt := range tests {
		_, err := LoadConfig(files, tt.name)
		ne, ok := errors.AsType[*strconv.NumError](err)
		if !ok {
			t.Errorf("LoadConfig(%q) error = %v, want a wrapped *strconv.NumError", tt.name, err)
			continue
		}
		if ne.Num != tt.num {
			t.Errorf("LoadConfig(%q) NumError.Num = %q, want %q", tt.name, ne.Num, tt.num)
		}
		if !errors.Is(err, tt.cause) {
			t.Errorf("LoadConfig(%q) error = %v, want %v", tt.name, err, tt.cause)
		}
		if errors.Is(err, ErrPortRange) {
			t.Errorf("LoadConfig(%q) error = %v, should not match ErrPortRange", tt.name, err)
		}
		if err.Error() != tt.msg {
			t.Errorf("LoadConfig(%q) error = %q, want %q", tt.name, err.Error(), tt.msg)
		}
	}
}

func TestParseHasNoFileName(t *testing.T) {
	_, err := parse("localhost")
	if err != ErrMissingPort {
		t.Errorf("parse(%q) error = %v, want ErrMissingPort itself (LoadConfig adds the context)", "localhost", err)
	}
}

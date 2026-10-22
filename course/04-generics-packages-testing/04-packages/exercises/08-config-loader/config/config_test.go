package config_test

import (
	"errors"
	"strings"
	"testing"
	"time"

	"learngo/course/04-generics-packages-testing/04-packages/exercises/08-config-loader/config"
)

func TestDefault(t *testing.T) {
	want := config.Config{Host: "localhost", Port: 8080, Debug: false, Timeout: 5 * time.Second}
	if got := config.Default(); got != want {
		t.Errorf("Default() = %+v, want %+v", got, want)
	}
}

func TestLoadOK(t *testing.T) {
	def := config.Config{Host: "localhost", Port: 8080, Timeout: 5 * time.Second}
	tests := []struct {
		in   map[string]string
		want config.Config
	}{
		{nil, def},
		{map[string]string{}, def},
		{map[string]string{"port": "9090", "debug": "true"}, config.Config{Host: "localhost", Port: 9090, Debug: true, Timeout: 5 * time.Second}},
		{map[string]string{"host": "example.com"}, config.Config{Host: "example.com", Port: 8080, Timeout: 5 * time.Second}},
		{map[string]string{"timeout": "1m30s", "debug": "0"}, config.Config{Host: "localhost", Port: 8080, Timeout: 90 * time.Second}},
		{map[string]string{"host": "h", "port": "1", "debug": "T", "timeout": "250ms"}, config.Config{Host: "h", Port: 1, Debug: true, Timeout: 250 * time.Millisecond}},
		{map[string]string{"port": "65535"}, config.Config{Host: "localhost", Port: 65535, Timeout: 5 * time.Second}},
	}
	for _, tt := range tests {
		got, err := config.Load(tt.in)
		if err != nil || got != tt.want {
			t.Errorf("Load(%v) = %+v, %v, want %+v, <nil>", tt.in, got, err, tt.want)
		}
	}
}

func TestLoadErrors(t *testing.T) {
	tests := []struct {
		in      map[string]string
		want    []error
		notWant []error
	}{
		{map[string]string{"colour": "red"}, []error{config.ErrUnknownKey}, []error{config.ErrBadValue, config.ErrInvalid}},
		{map[string]string{"Port": "80"}, []error{config.ErrUnknownKey}, nil},
		{map[string]string{"port": "abc"}, []error{config.ErrBadValue}, []error{config.ErrUnknownKey, config.ErrInvalid}},
		{map[string]string{"debug": "maybe"}, []error{config.ErrBadValue}, nil},
		{map[string]string{"timeout": "5"}, []error{config.ErrBadValue}, nil},
		{map[string]string{"port": "70000"}, []error{config.ErrInvalid}, []error{config.ErrBadValue}},
		{map[string]string{"port": "0"}, []error{config.ErrInvalid}, nil},
		{map[string]string{"host": ""}, []error{config.ErrInvalid}, nil},
		{map[string]string{"timeout": "-1s"}, []error{config.ErrInvalid}, nil},
		{map[string]string{"zzz": "1", "port": "x", "host": ""}, []error{config.ErrUnknownKey, config.ErrBadValue, config.ErrInvalid}, nil},
	}
	for _, tt := range tests {
		got, err := config.Load(tt.in)
		if got != (config.Config{}) {
			t.Errorf("Load(%v) config = %+v, want {}", tt.in, got)
		}
		for _, w := range tt.want {
			if !errors.Is(err, w) {
				t.Errorf("Load(%v) error = %v, want %v", tt.in, err, w)
			}
		}
		for _, nw := range tt.notWant {
			if errors.Is(err, nw) {
				t.Errorf("Load(%v) error = %v, want no %v", tt.in, err, nw)
			}
		}
	}
}

func TestLoadErrorMessage(t *testing.T) {
	in := map[string]string{"zzz": "1", "port": "x", "host": ""}
	want := strings.Join([]string{
		`config: bad value: port="x"`,
		`config: unknown key: "zzz"`,
		`config: invalid: host is empty`,
	}, "\n")
	_, err := config.Load(in)
	if err == nil || err.Error() != want {
		t.Errorf("Load(%v) error = %v, want %q", in, err, want)
	}
}

func TestValidate(t *testing.T) {
	valid := config.Default()
	if err := valid.Validate(); err != nil {
		t.Errorf("Default().Validate() = %v, want <nil>", err)
	}
	c := config.Config{Host: "h", Port: 1, Timeout: time.Nanosecond}
	if err := c.Validate(); err != nil {
		t.Errorf("%+v.Validate() = %v, want <nil>", c, err)
	}
	want := strings.Join([]string{
		"config: invalid: host is empty",
		"config: invalid: port 0 out of range",
		"config: invalid: timeout must be positive",
	}, "\n")
	err := config.Config{}.Validate()
	if !errors.Is(err, config.ErrInvalid) || err.Error() != want {
		t.Errorf("Config{}.Validate() = %v, want %q", err, want)
	}
	big := config.Config{Host: "h", Port: 65536, Timeout: time.Second}
	if err := big.Validate(); err == nil || err.Error() != "config: invalid: port 65536 out of range" {
		t.Errorf("%+v.Validate() = %v, want %q", big, err, "config: invalid: port 65536 out of range")
	}
}

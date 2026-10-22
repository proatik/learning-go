package daemon_test

import (
	"errors"
	"testing"

	daemon "learngo/course/04-generics-packages-testing/04-packages/exercises/08-config-loader"
	"learngo/course/04-generics-packages-testing/04-packages/exercises/08-config-loader/config"
)

func TestDescribe(t *testing.T) {
	tests := []struct {
		in   map[string]string
		want string
	}{
		{nil, "localhost:8080 debug=false timeout=5s"},
		{map[string]string{"host": "example.com"}, "example.com:8080 debug=false timeout=5s"},
		{map[string]string{"port": "443", "debug": "true", "timeout": "2m"}, "localhost:443 debug=true timeout=2m0s"},
	}
	for _, tt := range tests {
		got, err := daemon.Describe(tt.in)
		if got != tt.want || err != nil {
			t.Errorf("Describe(%v) = %q, %v, want %q, <nil>", tt.in, got, err, tt.want)
		}
	}
}

func TestDescribeError(t *testing.T) {
	in := map[string]string{"port": "nope"}
	got, err := daemon.Describe(in)
	if got != "" || !errors.Is(err, config.ErrBadValue) {
		t.Errorf("Describe(%v) = %q, %v, want %q, %v", in, got, err, "", config.ErrBadValue)
	}
}

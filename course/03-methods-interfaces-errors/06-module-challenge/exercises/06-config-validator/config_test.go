package config_test

import (
	"errors"
	"fmt"
	"slices"
	"testing"

	config "learngo/course/03-methods-interfaces-errors/06-module-challenge/exercises/06-config-validator"
)

func valid() config.Config {
	return config.Config{Name: "api", Port: 8080, Env: "prod", Replicas: 3, Tags: []string{"web", "public"}}
}

func TestFieldError(t *testing.T) {
	fe := &config.FieldError{Field: "port", Err: config.ErrOutOfRange}
	var err error = fe
	if got := err.Error(); got != "port: out of range" {
		t.Errorf("FieldError{port, ErrOutOfRange}.Error() = %q, want %q", got, "port: out of range")
	}
	if !errors.Is(err, config.ErrOutOfRange) {
		t.Errorf("errors.Is(%v, ErrOutOfRange) = false, want true (add an Unwrap method)", err)
	}
	if errors.Is(err, config.ErrRequired) {
		t.Errorf("errors.Is(%v, ErrRequired) = true, want false", err)
	}
	wrapped := fmt.Errorf("load: %w", err)
	if got, ok := errors.AsType[*config.FieldError](wrapped); !ok || got != fe {
		t.Errorf("errors.AsType[*FieldError](%v) = %v, %v, want the same *FieldError, true", wrapped, got, ok)
	}
}

func TestValidConfigs(t *testing.T) {
	tests := []config.Config{
		valid(),
		{Name: "x", Port: 1, Env: "dev", Replicas: 1},
		{Name: "x", Port: 65535, Env: "staging", Replicas: 10, Tags: []string{}},
		{Name: " worker ", Port: 9000, Env: "dev", Replicas: 1, Tags: []string{"a", "A"}}, // tags are case-sensitive
		{Name: "db", Port: 5432, Env: "prod", Replicas: 2},
	}
	for _, c := range tests {
		// Must be a true nil interface, not a typed nil.
		if err := c.Validate(); err != nil {
			t.Errorf("Validate(%+v) = %v, want nil", c, err)
		}
	}
}

func TestSingleProblems(t *testing.T) {
	tests := []struct {
		mod   func(*config.Config)
		field string
		want  error
	}{
		{func(c *config.Config) { c.Name = "" }, "name", config.ErrRequired},
		{func(c *config.Config) { c.Name = " \t " }, "name", config.ErrRequired},
		{func(c *config.Config) { c.Port = 0 }, "port", config.ErrOutOfRange},
		{func(c *config.Config) { c.Port = -80 }, "port", config.ErrOutOfRange},
		{func(c *config.Config) { c.Port = 65536 }, "port", config.ErrOutOfRange},
		{func(c *config.Config) { c.Env = "" }, "env", config.ErrRequired},
		{func(c *config.Config) { c.Env = "Prod" }, "env", config.ErrInvalid},
		{func(c *config.Config) { c.Env = "test" }, "env", config.ErrInvalid},
		{func(c *config.Config) { c.Replicas = 0 }, "replicas", config.ErrOutOfRange},
		{func(c *config.Config) { c.Replicas = 11 }, "replicas", config.ErrOutOfRange},
		{func(c *config.Config) { c.Replicas = 1 }, "replicas", config.ErrInvalid}, // prod needs >= 2
		{func(c *config.Config) { c.Tags = []string{"a", ""} }, "tags[1]", config.ErrRequired},
		{func(c *config.Config) { c.Tags = []string{"  "} }, "tags[0]", config.ErrRequired},
		{func(c *config.Config) { c.Tags = []string{"a", "b", "a"} }, "tags[2]", config.ErrInvalid},
	}
	for _, tt := range tests {
		c := valid()
		tt.mod(&c)
		err := c.Validate()
		if err == nil {
			t.Errorf("Validate(%+v) = nil, want a %s error", c, tt.field)
			continue
		}
		if !errors.Is(err, tt.want) {
			t.Errorf("Validate(%+v) error = %v, want %v", c, err, tt.want)
		}
		fe, ok := errors.AsType[*config.FieldError](err)
		if !ok {
			t.Errorf("Validate(%+v) error = %v (%T), want it to contain a *FieldError", c, err, err)
			continue
		}
		if fe.Field != tt.field || fe.Err != tt.want {
			t.Errorf("Validate(%+v) FieldError = {%s %v}, want {%s %v}", c, fe.Field, fe.Err, tt.field, tt.want)
		}
		if want := tt.field + ": " + tt.want.Error(); err.Error() != want {
			t.Errorf("Validate(%+v) error = %q, want %q", c, err.Error(), want)
		}
	}
}

func TestAllProblemsReported(t *testing.T) {
	c := config.Config{Port: 70000, Env: "qa", Replicas: 0, Tags: []string{"x", "", "x", "y", "y"}}
	err := c.Validate()
	if err == nil {
		t.Fatalf("Validate(%+v) = nil, want every problem joined", c)
	}
	want := "name: required\nport: out of range\nenv: invalid\nreplicas: out of range\ntags[1]: required\ntags[2]: invalid\ntags[4]: invalid"
	if err.Error() != want {
		t.Errorf("Validate(%+v) error =\n%s\nwant\n%s", c, err.Error(), want)
	}
	for _, sentinel := range []error{config.ErrRequired, config.ErrOutOfRange, config.ErrInvalid} {
		if !errors.Is(err, sentinel) {
			t.Errorf("errors.Is(err, %v) = false, want true (errors.Join searches every branch)", sentinel)
		}
	}
	// errors.As / AsType stop at the FIRST match in the tree.
	if fe, ok := errors.AsType[*config.FieldError](err); !ok || fe.Field != "name" {
		t.Errorf("errors.AsType[*FieldError](err) = %v, %v, want the name error, true", fe, ok)
	}
	wantFields := []string{"name", "port", "env", "replicas", "tags[1]", "tags[2]", "tags[4]"}
	if got := config.Fields(err); !slices.Equal(got, wantFields) {
		t.Errorf("Fields(Validate(...)) = %q, want %q", got, wantFields)
	}
}

func TestProdReplicasNotDoubleReported(t *testing.T) {
	c := valid()
	c.Replicas = 0 // out of range: report only that, not also "invalid"
	if got := config.Fields(c.Validate()); !slices.Equal(got, []string{"replicas"}) {
		t.Errorf("Fields(Validate(replicas=0, prod)) = %q, want [\"replicas\"]", got)
	}
	if err := c.Validate(); errors.Is(err, config.ErrInvalid) {
		t.Errorf("Validate(replicas=0, prod) error = %v, want only ErrOutOfRange", err)
	}
}

func TestFields(t *testing.T) {
	fe := func(f string) error { return &config.FieldError{Field: f, Err: config.ErrInvalid} }
	other := errors.New("disk full")
	tests := []struct {
		name string
		err  error
		want []string
	}{
		{"nil", nil, nil},
		{"plain error", other, nil},
		{"one field", fe("a"), []string{"a"}},
		{"wrapped", fmt.Errorf("load: %w", fe("a")), []string{"a"}},
		{"double wrapped", fmt.Errorf("x: %w", fmt.Errorf("y: %w", fe("deep"))), []string{"deep"}},
		{"joined", errors.Join(fe("a"), other, fe("b")), []string{"a", "b"}},
		{"nested joins", errors.Join(fe("a"), errors.Join(fe("b"), fmt.Errorf("w: %w", fe("c"))), fe("d")), []string{"a", "b", "c", "d"}},
		{"multi %w", fmt.Errorf("%w and %w", fe("p"), fe("q")), []string{"p", "q"}},
		{"field wrapping a field", &config.FieldError{Field: "outer", Err: fe("inner")}, []string{"outer", "inner"}},
	}
	for _, tt := range tests {
		if got := config.Fields(tt.err); !slices.Equal(got, tt.want) {
			t.Errorf("Fields(%s: %v) = %q, want %q", tt.name, tt.err, got, tt.want)
		}
	}
}

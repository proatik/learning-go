package signup

import (
	"errors"
	"fmt"
	"slices"
	"testing"
)

var all = []error{
	ErrUsernameRequired, ErrUsernameInvalid, ErrEmailInvalid, ErrPasswordShort,
	ErrPasswordNoDigit, ErrPasswordMismatch, ErrUnderage, ErrTerms,
}

func good() Form {
	return Form{Username: "ann_99", Email: "ann@go.dev", Password: "s3cretpass", Confirm: "s3cretpass", Age: 30, Terms: true}
}

func TestValidateOK(t *testing.T) {
	forms := []Form{good()}
	f := good()
	f.Username, f.Age = "Bob", 18
	forms = append(forms, f)
	f = good()
	f.Username = "abcdefghijklmnop" // 16
	f.Password, f.Confirm = "пароль12", "пароль12"
	forms = append(forms, f)
	for _, f := range forms {
		if err := Validate(f); err != nil {
			t.Errorf("Validate(%+v) error = %v, want <nil>", f, err)
		}
	}
}

func TestValidateEachProblem(t *testing.T) {
	tests := []struct {
		name   string
		change func(*Form)
		want   []error
	}{
		{"empty username", func(f *Form) { f.Username = "" }, []error{ErrUsernameRequired}},
		{"short username", func(f *Form) { f.Username = "ab" }, []error{ErrUsernameInvalid}},
		{"long username", func(f *Form) { f.Username = "abcdefghijklmnopq" }, []error{ErrUsernameInvalid}},
		{"bad char", func(f *Form) { f.Username = "ann-99" }, []error{ErrUsernameInvalid}},
		{"non-ASCII", func(f *Form) { f.Username = "zoë" }, []error{ErrUsernameInvalid}},
		{"no at", func(f *Form) { f.Email = "ann.go.dev" }, []error{ErrEmailInvalid}},
		{"two ats", func(f *Form) { f.Email = "a@b@go.dev" }, []error{ErrEmailInvalid}},
		{"no user", func(f *Form) { f.Email = "@go.dev" }, []error{ErrEmailInvalid}},
		{"no dot", func(f *Form) { f.Email = "ann@localhost" }, []error{ErrEmailInvalid}},
		{"empty email", func(f *Form) { f.Email = "" }, []error{ErrEmailInvalid}},
		{"short pw", func(f *Form) { f.Password, f.Confirm = "abc1", "abc1" }, []error{ErrPasswordShort}},
		{"no digit", func(f *Form) { f.Password, f.Confirm = "abcdefgh", "abcdefgh" }, []error{ErrPasswordNoDigit}},
		{"short + no digit", func(f *Form) { f.Password, f.Confirm = "abc", "abc" }, []error{ErrPasswordShort, ErrPasswordNoDigit}},
		{"mismatch", func(f *Form) { f.Confirm = "s3cretpasS" }, []error{ErrPasswordMismatch}},
		{"underage", func(f *Form) { f.Age = 17 }, []error{ErrUnderage}},
		{"terms", func(f *Form) { f.Terms = false }, []error{ErrTerms}},
		{"everything", func(f *Form) { *f = Form{Username: "a!", Email: "nope", Password: "short", Confirm: "other", Age: 12} },
			[]error{ErrUsernameInvalid, ErrEmailInvalid, ErrPasswordShort, ErrPasswordNoDigit, ErrPasswordMismatch, ErrUnderage, ErrTerms}},
		{"zero form", func(f *Form) { *f = Form{} },
			[]error{ErrUsernameRequired, ErrEmailInvalid, ErrPasswordShort, ErrPasswordNoDigit, ErrUnderage, ErrTerms}},
	}
	for _, tt := range tests {
		f := good()
		tt.change(&f)
		err := Validate(f)
		for _, s := range all {
			wantIt := slices.Contains(tt.want, s)
			if got := errors.Is(err, s); got != wantIt {
				if wantIt {
					t.Errorf("Validate(%+v) error = %v, want %v", f, err, s)
				} else {
					t.Errorf("Validate(%+v) error = %v, should not match %v", f, err, s)
				}
			}
		}
	}
}

func TestValidateMessages(t *testing.T) {
	f := Form{Username: "ann", Email: "nope", Password: "abcdefgh1", Confirm: "abcdefgh1", Age: 17, Terms: true}
	err := Validate(f)
	want := "email \"nope\": invalid email\nage 17: must be at least 18"
	if err == nil {
		t.Fatalf("Validate(%+v) error = <nil>, want %q", f, want)
	}
	if got := err.Error(); got != want {
		t.Errorf("Validate(%+v) error = %q, want %q", f, got, want)
	}
	gotMsgs := Messages(err)
	wantMsgs := []string{`email "nope": invalid email`, "age 17: must be at least 18"}
	if !slices.Equal(gotMsgs, wantMsgs) {
		t.Errorf("Messages(%q) = %q, want %q", err.Error(), gotMsgs, wantMsgs)
	}
}

func TestMessages(t *testing.T) {
	a, b, c := errors.New("a"), errors.New("b"), errors.New("c")
	tests := []struct {
		err  error
		want []string
	}{
		{nil, nil},
		{a, []string{"a"}},
		{fmt.Errorf("ctx: %w", a), []string{"ctx: a"}},
		{errors.Join(a, b, c), []string{"a", "b", "c"}},
		{errors.Join(a, nil, fmt.Errorf("x: %w", b)), []string{"a", "x: b"}},
		{fmt.Errorf("%w + %w", a, b), []string{"a", "b"}},
	}
	for _, tt := range tests {
		got := Messages(tt.err)
		if !slices.Equal(got, tt.want) || (tt.want == nil) != (got == nil) {
			t.Errorf("Messages(%v) = %q, want %q", tt.err, got, tt.want)
		}
	}
}

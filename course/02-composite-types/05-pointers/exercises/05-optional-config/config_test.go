package config

import "testing"

func TestEffective(t *testing.T) {
	tests := []struct {
		name string
		cfg  Config
		want Settings
	}{
		{"all defaults", Config{}, Settings{3, true, "app"}},
		{"explicit zero values", Config{Retries: new(0), Cache: new(false), Name: new("")}, Settings{0, false, ""}},
		{"only retries", Config{Retries: new(10)}, Settings{10, true, "app"}},
		{"only cache", Config{Cache: new(false)}, Settings{3, false, "app"}},
		{"only name", Config{Name: new("api")}, Settings{3, true, "api"}},
		{"cache explicitly true", Config{Cache: new(true)}, Settings{3, true, "app"}},
		{"all set", Config{Retries: new(1), Cache: new(false), Name: new("db")}, Settings{1, false, "db"}},
	}
	for _, tt := range tests {
		if got := Effective(tt.cfg); got != tt.want {
			t.Errorf("%s: Effective(...) = %+v, want %+v", tt.name, got, tt.want)
		}
	}
}

func TestEffectiveReadsCurrentValue(t *testing.T) {
	n := 1
	c := Config{Retries: &n}
	n = 7
	if got := Effective(c).Retries; got != 7 {
		t.Errorf("Effective after changing the pointed-to int to 7: Retries = %d, want 7", got)
	}
}

func TestMerge(t *testing.T) {
	tests := []struct {
		name           string
		base, override Config
		want           Settings
	}{
		{"both empty", Config{}, Config{}, Settings{3, true, "app"}},
		{"base only", Config{Retries: new(5)}, Config{}, Settings{5, true, "app"}},
		{"override only", Config{}, Config{Name: new("api")}, Settings{3, true, "api"}},
		{"disjoint fields", Config{Retries: new(5)}, Config{Name: new("api")}, Settings{5, true, "api"}},
		{"override wins", Config{Retries: new(5), Cache: new(true)}, Config{Retries: new(0), Cache: new(false)}, Settings{0, false, "app"}},
		{"nil override keeps base", Config{Name: new("base")}, Config{Retries: new(2)}, Settings{2, true, "base"}},
	}
	for _, tt := range tests {
		got := Effective(Merge(tt.base, tt.override))
		if got != tt.want {
			t.Errorf("%s: Effective(Merge(...)) = %+v, want %+v", tt.name, got, tt.want)
		}
	}
}

func TestMergeKeepsNil(t *testing.T) {
	m := Merge(Config{Retries: new(1)}, Config{})
	if m.Retries == nil {
		t.Fatalf("Merge lost base.Retries: got nil")
	}
	if m.Cache != nil || m.Name != nil {
		t.Errorf("Merge(Config{Retries: ...}, Config{}): Cache=%v Name=%v, want both nil", m.Cache, m.Name)
	}
}

func TestMergeDoesNotModifyInputs(t *testing.T) {
	base := Config{Retries: new(5)}
	override := Config{Retries: new(9), Name: new("x")}
	Merge(base, override)
	if base.Retries == nil || *base.Retries != 5 || base.Name != nil {
		t.Errorf("Merge modified base")
	}
	if override.Retries == nil || *override.Retries != 9 {
		t.Errorf("Merge modified override")
	}
}

func TestFromSettings(t *testing.T) {
	for _, s := range []Settings{{0, false, ""}, {3, true, "app"}, {42, false, "svc"}} {
		c := FromSettings(s)
		if c.Retries == nil || c.Cache == nil || c.Name == nil {
			t.Fatalf("FromSettings(%+v) left a field nil: %+v", s, c)
		}
		if *c.Retries != s.Retries || *c.Cache != s.Cache || *c.Name != s.Name {
			t.Errorf("FromSettings(%+v) = {%d %v %q}, want {%d %v %q}", s, *c.Retries, *c.Cache, *c.Name, s.Retries, s.Cache, s.Name)
		}
		if got := Effective(c); got != s {
			t.Errorf("Effective(FromSettings(%+v)) = %+v, want the same settings back", s, got)
		}
	}
}

func TestFromSettingsFreshPointers(t *testing.T) {
	a := FromSettings(Settings{1, true, "a"})
	b := FromSettings(Settings{1, true, "a"})
	if a.Retries == nil || b.Retries == nil {
		t.Fatalf("FromSettings returned nil Retries")
	}
	if a.Retries == b.Retries {
		t.Errorf("two FromSettings calls returned the same *int; each call needs its own variable")
	}
	*a.Retries = 100
	if *b.Retries != 1 {
		t.Errorf("changing one Config's Retries changed another's")
	}
}

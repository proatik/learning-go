package config

// Config holds optional settings. A nil field means "not set: use the default".
type Config struct {
	Retries *int    // default 3
	Cache   *bool   // default true
	Name    *string // default "app"
}

// Settings holds the final, effective values.
type Settings struct {
	Retries int
	Cache   bool
	Name    string
}

// Effective returns the settings described by c, using the default for
// every nil field.
func Effective(c Config) Settings {
	s := Settings{Retries: 3, Cache: true, Name: "app"} // start from the defaults
	if c.Retries != nil {
		s.Retries = *c.Retries
	}
	if c.Cache != nil {
		s.Cache = *c.Cache
	}
	if c.Name != nil {
		s.Name = *c.Name
	}
	return s
}

// Merge returns a Config whose fields come from override when they are
// non-nil and from base otherwise.
func Merge(base, override Config) Config {
	// base is a copy, so changing its fields does not affect the caller.
	// Only the pointers are copied: the result shares the pointed-to values.
	if override.Retries != nil {
		base.Retries = override.Retries
	}
	if override.Cache != nil {
		base.Cache = override.Cache
	}
	if override.Name != nil {
		base.Name = override.Name
	}
	return base
}

// FromSettings returns a Config with every field set to the value in s.
func FromSettings(s Settings) Config {
	// new(expr) (Go 1.26+) allocates a variable initialized to expr's value.
	return Config{
		Retries: new(s.Retries),
		Cache:   new(s.Cache),
		Name:    new(s.Name),
	}
}

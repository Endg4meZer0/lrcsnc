package romanization

type Config struct {
	Korean  bool `toml:"korean"`
	Chinese bool `toml:"chinese"`
}

// IsEnabled returns true if at least one of the supported romanization options
// is turned on.
func (c *Config) IsEnabled() bool {
	return c.Korean || c.Chinese
}

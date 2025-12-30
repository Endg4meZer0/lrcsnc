package romanization

type Config struct {
	Japanese bool `toml:"japanese"`
	Chinese  bool `toml:"chinese"`
	Korean   bool `toml:"korean"`
}

// IsEnabled returns true if at least one of the supported romanization options
// is turned on.
func (c *Config) IsEnabled() bool {
	return c.Japanese || c.Chinese || c.Korean
}
package client

import "lrcsnc/internal/config/struct/client/format"

type Config struct {
	Destination   string        `toml:"destination"`
	Template      string        `toml:"template"`
	InsertNewline bool          `toml:"insert-newline"`
	Format        format.Config `toml:"format"`
}

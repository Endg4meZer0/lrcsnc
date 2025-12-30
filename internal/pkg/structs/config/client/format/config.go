package format

type Config struct {
	Lyric          string             `toml:"lyric"`
	Multiplier     string             `toml:"multiplier"`
	NotPlaying     string             `toml:"not-playing"`
	NoLyrics       string             `toml:"no-lyrics"`
	NoSyncedLyrics string             `toml:"no-synced-lyrics"`
	LoadingLyrics  string             `toml:"loading-lyrics"`
	ErrorMessage   string             `toml:"error-message"`
	Instrumental   InstrumentalConfig `toml:"instrumental"`
}

type InstrumentalConfig struct {
	Interval   float64 `toml:"interval"`
	Symbol     string  `toml:"symbol"`
	MaxSymbols uint    `toml:"max-symbols"`
}

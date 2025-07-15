package structs

import (
	"lrcsnc/internal/pkg/types"
)

// LEVEL 0

type Config struct {
	// Player config is for player related things. Currently it is used
	// for specifying included/excluded players for the watcher.
	Player PlayerConfig `toml:"player"`
	// Lyrics config currently has stuff to do with lyrics providers,
	// time offset and romanization
	Lyrics LyricsConfig `toml:"lyrics"`
	// Cache config has an "enabled" toggle, dir path and life span
	Cache CacheConfig `toml:"cache"`
	// Output config has... a lot of personalized settings.
	Output OutputConfig `toml:"output"`
}

// LEVEL 1

type PlayerConfig struct {
	IncludedPlayers []string `toml:"included-players"`
	ExcludedPlayers []string `toml:"excluded-players"`
}

type LyricsConfig struct {
	Provider        types.LyricsProviderType `toml:"provider"`
	TimestampOffset float64                  `toml:"timestamp-offset"`
	Romanization    RomanizationConfig       `toml:"romanization"`
}

type CacheConfig struct {
	Enabled        bool                      `toml:"enabled"`
	Dir            string                    `toml:"dir"`
	LifeSpan       uint                      `toml:"life-span"`
	StoreCondition CacheStoreConditionConfig `toml:"store-condition"`
}

type OutputConfig struct {
	Type  types.OutputType  `toml:"type"`
	Piped PipedOutputConfig `toml:"piped"`
}

// LEVEL 2

type RomanizationConfig struct {
	Japanese bool `toml:"japanese"`
	Chinese  bool `toml:"chinese"`
	Korean   bool `toml:"korean"`
}

// IsEnabled returns true if at least one of the supported romanization options 
// is turned on.
func (r *RomanizationConfig) IsEnabled() bool {
	return r.Japanese || r.Chinese || r.Korean
}

type CacheStoreConditionConfig struct {
	IfSynced       bool `toml:"if-synced"`       // LyricsState.Synced
	IfPlain        bool `toml:"if-plain"`        // LyricsState.Plain
	IfInstrumental bool `toml:"if-instrumental"` // LyricsState.Instrumental
}

// IsEnabledFor returns true for a LyricsState `ls`,
// if the option, corresponding to it, is enabled. Otherwise, returns false.
func (c *CacheStoreConditionConfig) IsEnabledFor(ls types.LyricsState) bool {
	switch ls {
	case types.LyricsStateSynced:
		return c.IfSynced
	case types.LyricsStatePlain:
		return c.IfPlain
	case types.LyricsStateInstrumental:
		return c.IfInstrumental
	default:
		return false
	}
}

type PipedOutputConfig struct {
	Destination    string                 `toml:"destination"`
	JSON           types.JSONOutputType   `toml:"json"`
	JSONWaybar     JSONWaybarOutputConfig `toml:"json-waybar"`
	InsertNewline  bool                   `toml:"insert-newline"`
	Text           FormatOutputConfig     `toml:"text"`
	Multiplier     FormatOutputConfig     `toml:"multiplier"`
	Lyric          LyricOutputConfig      `toml:"lyric"`
	NotPlaying     NotPlayingOutputConfig `toml:"not-playing"`
	SongNotFound   MessageOutputConfig    `toml:"song-not-found"`
	NoSyncedLyrics MessageOutputConfig    `toml:"no-synced-lyrics"`
	LoadingLyrics  MessageOutputConfig    `toml:"loading-lyrics"`
	ErrorMessage   MessageOutputConfig    `toml:"error-message"`
	Instrumental   InstrumentalConfig     `toml:"instrumental"`
}

// LEVEL 3

type JSONWaybarOutputConfig struct {
	Alt     string `toml:"alt"`
	Tooltip string `toml:"tooltip"`
	Class   string `toml:"class"`
}

type FormatOutputConfig struct {
	Format string `toml:"format"`
}

type LyricOutputConfig struct {
	Icon string `toml:"icon"`
}

type NotPlayingOutputConfig struct {
	Text string `toml:"text"`
}

type MessageOutputConfig struct {
	Enabled bool   `toml:"enabled"`
	Icon    string `toml:"icon"`
	Text    string `toml:"text"`
}

type InstrumentalConfig struct {
	Interval   float64 `toml:"interval"`
	Symbol     string  `toml:"symbol"`
	MaxSymbols uint    `toml:"max-symbols"`
}

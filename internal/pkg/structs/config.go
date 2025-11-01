package structs

import (
	"lrcsnc/internal/pkg/types"
)

// LEVEL 0

type Config struct {
	// Player config is for player related things. Currently it is used
	// for specifying included/excluded players for the watcher.
	Player PlayerConfig `toml:"player"`
	// Net config is for client-server mode settings.
	Net NetConfig `toml:"net"`
	// Lyrics config currently has stuff to do with lyrics providers,
	// time offset and romanization
	Lyrics LyricsConfig `toml:"lyrics"`
	// Cache config has an "enabled" toggle, dir path and life span
	Cache CacheConfig `toml:"cache"`
	// Client config has lots of personalized settings for the
	// lrcsnc's native client that is used in standalone and client modes.
	// Stuff like what kind of output should be going from lrcsnc is described here.
	Client ClientConfig `toml:"client"`
}

// LEVEL 1

type PlayerConfig struct {
	IncludedPlayers []string `toml:"included-players"`
	ExcludedPlayers []string `toml:"excluded-players"`
}

type NetConfig struct {
	// IsServer is a boolean that exists mostly for the purposes
	// of making an easier way to code the switch to server mode.
	// It can be set to true in the config file as well:
	// though it is not recommended to do that in your main config file.
	IsServer bool `toml:"is-server"`
	// Protocol represents the protocol used to connect or to start a server
	// on. Can use values available in net.Listen (so "unix", "tcp" and some other variations)
	// and empty for standalone mode.
	Protocol string `toml:"protocol"`
	// ListenAt represents the path of a server to listen to.
	// If lrcsnc fails to connect to a server at this path, then it crashes.
	ListenAt string `toml:"listen-at"`
}

type LyricsConfig struct {
	Provider     types.LyricsProviderType `toml:"provider"`
	TimingOffset float64                  `toml:"timing-offset"`
	Romanization RomanizationConfig       `toml:"romanization"`
}

type CacheConfig struct {
	Enabled        bool                      `toml:"enabled"`
	Dir            string                    `toml:"dir"`
	LifeSpan       uint                      `toml:"life-span"`
	StoreCondition CacheStoreConditionConfig `toml:"store-condition"`
}

type ClientConfig struct {
	Destination   string             `toml:"destination"`
	Template      string             `toml:"template"`
	InsertNewline bool               `toml:"insert-newline"`
	Format        FormatOutputConfig `toml:"format"`
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

type FormatOutputConfig struct {
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

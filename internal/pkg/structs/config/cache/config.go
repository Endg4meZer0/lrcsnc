package cache

import "lrcsnc/internal/pkg/types"

type Config struct {
	Enabled        bool                 `toml:"enabled"`
	Dir            string               `toml:"dir"`
	LifeSpan       uint                 `toml:"life-span"`
	StoreCondition StoreConditionConfig `toml:"store-condition"`
}

type StoreConditionConfig struct {
	IfSynced       bool `toml:"if-synced"`       // LyricsState.Synced
	IfPlain        bool `toml:"if-plain"`        // LyricsState.Plain
	IfInstrumental bool `toml:"if-instrumental"` // LyricsState.Instrumental
}

// IsEnabledFor returns true for a LyricsState `ls`,
// if the option, corresponding to it, is enabled. Otherwise, returns false.
func (c *StoreConditionConfig) IsEnabledFor(ls types.LyricsState) bool {
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

package providers

import (
	"lrcsnc/internal/lyrics/providers/local"
	lrclib "lrcsnc/internal/lyrics/providers/lrclib"

	lyricStruct "lrcsnc/internal/lyrics/struct"
	"lrcsnc/internal/pkg/types"
	playerStruct "lrcsnc/internal/player/struct"
)

type Provider interface {
	// Get returns the lyrics of a song in form of LyricsData
	Get(playerStruct.Song) (lyricStruct.LyricsData, error)
}

var Providers = map[types.LyricsProviderType]Provider{
	types.LyricsProviderLocal:  local.Provider{},
	types.LyricsProviderLrclib: lrclib.Provider{},
}

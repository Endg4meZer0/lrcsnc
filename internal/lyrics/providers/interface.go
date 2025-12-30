package providers

import (
	lrclib "lrcsnc/internal/lyrics/providers/lrclib"

	playerStructs "lrcsnc/internal/pkg/structs/player"
	"lrcsnc/internal/pkg/types"
)

type Provider interface {
	// Get returns the lyrics of a song in form of LyricsData
	Get(playerStructs.Song) (playerStructs.LyricsData, error)
}

var Providers = map[types.LyricsProviderType]Provider{
	types.LyricsProviderLrclib: lrclib.Provider{},
}

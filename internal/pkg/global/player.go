package global

import (
	lyricStruct "lrcsnc/internal/lyrics/struct"
	"lrcsnc/internal/pkg/types"
	playerStruct "lrcsnc/internal/player/struct"
	"sync"
)

var Player = struct {
	M sync.Mutex
	P playerStruct.Player
}{
	P: playerStruct.Player{
		PlaybackStatus:       types.PlaybackStatusStopped,
		LastDetectedPosition: 0.0,
		Rate:                 1.0,
		Song: playerStruct.Song{
			LyricsData: lyricStruct.LyricsData{
				LyricsState: types.LyricsStateUnknown,
			},
		},
	},
}

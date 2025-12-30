package global

import (
	playerStructs "lrcsnc/internal/pkg/structs/player"
	"lrcsnc/internal/pkg/types"
	"sync"

	"github.com/Endg4meZer0/go-mpris"
)

var Player = struct {
	M sync.Mutex
	P playerStructs.Player
}{
	P: playerStructs.Player{
		PlaybackStatus: mpris.PlaybackStopped,
		Position:       0.0,
		Rate:           1.0,
		Song: playerStructs.Song{
			LyricsData: playerStructs.LyricsData{
				LyricsState: types.LyricsStateUnknown,
			},
		},
	},
}

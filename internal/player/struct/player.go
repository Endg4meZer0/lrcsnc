package playerStruct

import (
	"lrcsnc/internal/pkg/types"
)

type Player struct {
	Name                 string
	PlaybackStatus       types.PlaybackStatus
	LastDetectedPosition float64
	Rate                 float64
	Song                 Song
}

func (p *Player) Update(name string, ps types.PlaybackStatus, ldp, rate float64, md SongMetadata) {
	p.Name = name
	p.PlaybackStatus = ps
	p.LastDetectedPosition = ldp
	p.Rate = rate
	p.Song.Metadata = md
	p.Song.LyricsData.LyricsState = types.LyricsStateLoading
}

package event

import (
	"lrcsnc/internal/pkg/structs"
	"lrcsnc/internal/pkg/types"

	"github.com/Endg4meZer0/go-mpris"
)

type EventTypeActiveLyricChangedData struct {
	// Index is the index of the new active lyric.
	Index int
	// Lyric is lyric itself. If the lyric is empty, it is
	// considered an instrumental lyric.
	Lyric structs.Lyric
	// Multiplier is the number of times the lyric repeated itself
	// in the text up to this moment.
	Multiplier int
	// TimeUntilEnd is calculated here using just the timing
	// of the current lyric and the next lyric (in seconds).
	TimeUntilEnd float64
	// Resync represents if a resynchronization happened
	// server-side. The client itself should decide onto
	// how to update the active lyric based on this and
	// lyric itself (e. g. if it changed due to resync).
	Resync bool
}

type EventTypeSongChangedData struct {
	Title    string
	Artists  []string
	Album    string
	Duration float64
}

type EventTypePlayerChangedData struct {
	Name string
}

type EventTypePlaybackStatusChangedData struct {
	PlaybackStatus mpris.PlaybackStatus
}

type EventTypeRateChangedData struct {
	Rate float64
}

type EventTypeLyricsStateChangedData struct {
	State types.LyricsState
}

type EventTypeLyricsChangedData struct {
	Lyrics []structs.Lyric
}

type EventTypeOverwriteRequiredData struct {
	Overwrite string
}

type EventTypeServerClosedData struct{}

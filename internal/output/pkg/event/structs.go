package event

import (
	lyricStruct "lrcsnc/internal/lyrics/struct"
	"lrcsnc/internal/pkg/types"
)

type EventTypeActiveLyricChangedData struct {
	// Index is the index of the new active lyric.
	Index int
	// Lyric is lyric itself. If the lyric is empty, it is
	// considered an instrumental lyric.
	Lyric lyricStruct.Lyric
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
	PlaybackStatus types.PlaybackStatus
}

type EventTypeRateChangedData struct {
	Rate float64
}

type EventTypeLyricsStateChangedData struct {
	State types.LyricsState
}

type EventTypeLyricsChangedData struct {
	Lyrics lyricStruct.Lyrics
}

type EventTypeOverwriteRequiredData struct {
	Overwrite string
}

type EventTypeServerClosedData struct{}

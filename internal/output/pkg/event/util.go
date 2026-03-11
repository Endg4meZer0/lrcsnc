package event

import (
	lyricStruct "lrcsnc/internal/lyrics/struct"
	"lrcsnc/internal/pkg/types"
)

func MakeDefaultEvent(et EventType) (e Event) {
	e.Type = et

	switch et {
	case EventTypePlayerChanged:
		e.Data = EventTypePlayerChangedData{
			Name: "",
		}
	case EventTypeRateChanged:
		e.Data = EventTypeRateChangedData{
			Rate: 1.00,
		}
	case EventTypeSongChanged:
		e.Data = EventTypeSongChangedData{
			Title:    "",
			Artists:  []string{""},
			Album:    "",
			Duration: 0,
		}
	case EventTypePlaybackStatusChanged:
		e.Data = EventTypePlaybackStatusChangedData{
			PlaybackStatus: types.PlaybackStatusStopped,
		}
	case EventTypeLyricsStateChanged:
		e.Data = EventTypeLyricsStateChangedData{
			State: types.LyricsStateLoading,
		}
	case EventTypeActiveLyricChanged:
		e.Data = EventTypeActiveLyricChangedData{
			Index:        -1,
			Lyric:        lyricStruct.Lyric{Timing: 0, Text: ""},
			Multiplier:   0,
			TimeUntilEnd: 0,
			Resync:       false,
		}
	case EventTypeOverwriteRequired:
		e.Data = EventTypeOverwriteRequiredData{
			Overwrite: "",
		}
	}

	return
}

package event

import (
	"lrcsnc/internal/pkg/structs"
	"lrcsnc/internal/pkg/types"

	"github.com/Endg4meZer0/go-mpris"
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
			PlaybackStatus: mpris.PlaybackStopped,
		}
	case EventTypeLyricsStateChanged:
		e.Data = EventTypeLyricsStateChangedData{
			State: types.LyricsStateLoading,
		}
	case EventTypeActiveLyricChanged:
		e.Data = EventTypeActiveLyricChangedData{
			Index:        -1,
			Lyric:        structs.Lyric{Timing: 0, Text: ""},
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

package sync

import (
	"errors"
	"time"

	"lrcsnc/internal/lyrics"
	errs "lrcsnc/internal/lyrics/errors"
	"lrcsnc/internal/output/pkg/event"
	"lrcsnc/internal/output/server"
	"lrcsnc/internal/pkg/global"
	"lrcsnc/internal/pkg/types"
)

var songChanged chan bool = make(chan bool)
var lastDownloadStart time.Time

func lyricFetcher() {
	for {
		<-songChanged

		go server.ReceiveEvent(event.Event{
			Type: event.EventTypeSongChanged,
			Data: event.EventTypeSongChangedData{
				Title:    global.Player.P.Song.Title,
				Artists:  global.Player.P.Song.Artists,
				Album:    global.Player.P.Song.Album,
				Duration: global.Player.P.Song.Duration,
			},
		})

		go server.ReceiveEvent(event.Event{
			Type: event.EventTypeLyricsStateChanged,
			Data: event.EventTypeLyricsStateChangedData{
				State: types.LyricsStateLoading,
			},
		})

		// This value will change on each new song changed event
		// so if the download takes too long and the song was switched
		// it can just store the necessary data in cache and forget about it
		lastDownloadStart = time.Now()

		go func() {
			thisDownloadStart := lastDownloadStart

			lyricsData, err := lyrics.Fetch()
			if err != nil && !errors.Is(err, errs.NotFound) {
				return
			}

			if thisDownloadStart != lastDownloadStart {
				return
			}

			lyrics.Configure(&lyricsData)

			global.Player.M.Lock()
			global.Player.P.Song.LyricsData = lyricsData
			global.Player.M.Unlock()

			go server.ReceiveEvent(event.Event{Type: event.EventTypeLyricsStateChanged, Data: event.EventTypeLyricsStateChangedData{
				State: global.Player.P.Song.LyricsData.LyricsState,
			}})

			go server.ReceiveEvent(event.Event{Type: event.EventTypeLyricsChanged, Data: event.EventTypeLyricsChangedData{
				Lyrics: global.Player.P.Song.LyricsData.Lyrics,
			}})

			// And finally, it ends with a position sync
			AskForPositionSync()
		}()
	}
}

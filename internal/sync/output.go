package sync

import (
	"math"
	"time"

	"lrcsnc/internal/output/pkg/event"
	"lrcsnc/internal/output/server"
	"lrcsnc/internal/pkg/global"
	playerStructs "lrcsnc/internal/pkg/structs/player"
	"lrcsnc/internal/pkg/types"
	"lrcsnc/internal/pkg/util"

	"github.com/Endg4meZer0/go-mpris"
)

var lyricsTimer = time.NewTimer(5 * time.Minute)
var writtenTiming float64
var resyncFlag = false

func resyncLyrics() {
	resyncFlag = true
	lyricsTimer.Reset(1)
}

func stopLyricsSync() {
	lyricsTimer.Stop()
}

func lyricsSynchronizer() {
	for {
		<-lyricsTimer.C
		if global.Player.P.Song.LyricsData.LyricsState != types.LyricsStateSynced {
			go server.ReceiveEvent(event.Event{
				Type: event.EventTypeActiveLyricChanged,
				Data: event.EventTypeActiveLyricChangedData{
					Index:        -1,
					Lyric:        playerStructs.Lyric{Timing: 0, Text: ""},
					Multiplier:   0,
					TimeUntilEnd: 0,
					Resync:       resyncFlag,
				},
			})
		} else {
			// 5999.99s is basically the maximum limit of .lrc files' timestamps AFAIK, so 6000s is unreachable
			currentLyricTiming := -1.0
			nextLyricTiming := 6000.0
			newLyricIndex := -1

			for i, lyric := range global.Player.P.Song.LyricsData.Lyrics {
				if lyric.Timing+global.Config.C.Lyrics.TimingOffset <= global.Player.P.Position && currentLyricTiming <= lyric.Timing+global.Config.C.Lyrics.TimingOffset {
					currentLyricTiming = lyric.Timing + global.Config.C.Lyrics.TimingOffset
					newLyricIndex = i
				}
			}

			if newLyricIndex != len(global.Player.P.Song.LyricsData.Lyrics)-1 {
				nextLyricTiming = global.Player.P.Song.LyricsData.Lyrics[newLyricIndex+1].Timing + global.Config.C.Lyrics.TimingOffset
			}

			lyricsTimerDuration := time.Duration(int64(math.Abs(nextLyricTiming-global.Player.P.Position)*1000)) * time.Millisecond

			if currentLyricTiming == -1 || (global.Player.P.PlaybackStatus == mpris.PlaybackPlaying && writtenTiming != currentLyricTiming) {
				lyric := playerStructs.Lyric{Timing: 0, Text: ""}
				if newLyricIndex >= 0 && newLyricIndex < len(global.Player.P.Song.LyricsData.Lyrics) {
					lyric = global.Player.P.Song.LyricsData.Lyrics[newLyricIndex]
				}
				go server.ReceiveEvent(event.Event{
					Type: event.EventTypeActiveLyricChanged,
					Data: event.EventTypeActiveLyricChangedData{
						Index:        newLyricIndex,
						Lyric:        lyric,
						Multiplier:   util.CalculateMultiplier(global.Player.P.Song.LyricsData.Lyrics, newLyricIndex),
						TimeUntilEnd: nextLyricTiming - global.Player.P.Position,
						Resync:       resyncFlag,
					},
				})
			}

			writtenTiming = currentLyricTiming
			global.Player.P.Position = nextLyricTiming
			lyricsTimer.Reset(lyricsTimerDuration)
		}
		resyncFlag = false
	}
}

package sync

import (
	"math"
	"time"

	lyricStruct "lrcsnc/internal/lyrics/struct"
	"lrcsnc/internal/output/pkg/event"
	"lrcsnc/internal/output/server"
	"lrcsnc/internal/pkg/global"
	"lrcsnc/internal/pkg/types"
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
					Lyric:        lyricStruct.Lyric{Timing: 0, Text: ""},
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
				if lyric.Timing+global.Config.C.Lyrics.TimingOffset <= global.Player.P.LastDetectedPosition && currentLyricTiming <= lyric.Timing+global.Config.C.Lyrics.TimingOffset {
					currentLyricTiming = lyric.Timing + global.Config.C.Lyrics.TimingOffset
					newLyricIndex = i
				}
			}

			if newLyricIndex != len(global.Player.P.Song.LyricsData.Lyrics)-1 {
				nextLyricTiming = global.Player.P.Song.LyricsData.Lyrics[newLyricIndex+1].Timing + global.Config.C.Lyrics.TimingOffset
			}

			lyricsTimerDuration := time.Duration(int64(math.Abs(nextLyricTiming-global.Player.P.LastDetectedPosition)*1000)) * time.Millisecond

			if currentLyricTiming == -1 || (global.Player.P.PlaybackStatus == types.PlaybackStatusPlaying && writtenTiming != currentLyricTiming) {
				lyric := lyricStruct.Lyric{Timing: 0, Text: ""}
				if newLyricIndex >= 0 && newLyricIndex < len(global.Player.P.Song.LyricsData.Lyrics) {
					lyric = global.Player.P.Song.LyricsData.Lyrics[newLyricIndex]
				}
				go server.ReceiveEvent(event.Event{
					Type: event.EventTypeActiveLyricChanged,
					Data: event.EventTypeActiveLyricChangedData{
						Index:        newLyricIndex,
						Lyric:        lyric,
						Multiplier:   global.Player.P.Song.LyricsData.Lyrics.CalculateMultiplierFor(newLyricIndex),
						TimeUntilEnd: nextLyricTiming - global.Player.P.LastDetectedPosition,
						Resync:       resyncFlag,
					},
				})
			}

			writtenTiming = currentLyricTiming
			global.Player.P.LastDetectedPosition = nextLyricTiming
			lyricsTimer.Reset(lyricsTimerDuration)
		}
		resyncFlag = false
	}
}

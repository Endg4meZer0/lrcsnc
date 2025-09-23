package sync

import (
	"math"
	"time"

	"lrcsnc/internal/output"
	"lrcsnc/internal/pkg/global"
	"lrcsnc/internal/pkg/types"

	"github.com/Endg4meZer0/go-mpris"
)

var lyricsTimer = time.NewTimer(5 * time.Minute)
var lyricIndex = -1
var writtenTiming float64

func resyncLyrics() {
	lyricsTimer.Reset(1)
}

func stopLyricsSync() {
	lyricsTimer.Stop()
}

func outputUpdate() {
	if global.Player.P.Song.LyricsData.LyricsState != types.LyricsStateSynced {
		output.Controllers[global.Config.C.Output.Type].DisplayLyric(-1)
	} else {
		output.Controllers[global.Config.C.Output.Type].DisplayLyric(lyricIndex)
	}
}

func lyricsSynchronizer() {
	for {
		<-lyricsTimer.C
		if global.Player.P.Song.LyricsData.LyricsState != types.LyricsStateSynced {
			output.Controllers[global.Config.C.Output.Type].DisplayLyric(-1)
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
				output.Controllers[global.Config.C.Output.Type].DisplayLyric(newLyricIndex)
			}

			lyricIndex = newLyricIndex
			writtenTiming = currentLyricTiming
			global.Player.P.Position = nextLyricTiming
			lyricsTimer.Reset(lyricsTimerDuration)
		}
	}
}

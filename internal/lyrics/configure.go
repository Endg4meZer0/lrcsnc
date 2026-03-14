package lyrics

import (
	"lrcsnc/internal/lyrics/romanization"
	lyricStruct "lrcsnc/internal/lyrics/struct"
	"lrcsnc/internal/pkg/log"
)

// Configure sets up the lyrics data by applying necessary configurations.
// It should be called after format decrypt and before the data is sent to
// the main sync goroutines by channels.
// Currently, it only applies romanization to the lyrics data.
// May be extended in the future.
//
// Every function/method/module/whatever needs to lock the mutex
// by themselves and only themselves.
// No locking a mutex in THIS function.
func Configure(lyricsData *lyricStruct.LyricsData) {
	log.Debug("lyrics/configure", "Starting configuring the received lyrics")

	// Romanization
	log.Debug("lyrics/configure", "Applying romanization if enabled and necessary")
	romanization.Romanize(lyricsData.Lyrics)

	log.Debug("lyrics/configure", "Done")
}

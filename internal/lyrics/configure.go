package lyrics

import (
	"sync"

	"github.com/longbridgeapp/opencc"

	"lrcsnc/internal/pkg/global"
	"lrcsnc/internal/pkg/log"
	playerStructs "lrcsnc/internal/pkg/structs/player"
	"lrcsnc/internal/romanization"
)

var (
	t2s        *opencc.OpenCC
	t2sOnce    sync.Once
	t2sInitErr error
)

func getT2S() *opencc.OpenCC {
	t2sOnce.Do(func() {
		t2s, t2sInitErr = opencc.New("t2s")
		if t2sInitErr != nil {
			log.Warn("lyrics/configure", "Failed to initialize OpenCC: "+t2sInitErr.Error())
		}
	})
	return t2s
}

// Configure sets up the lyrics data by applying necessary configurations.
// It should be called after format decrypt and before the data is sent to
// the main sync goroutines by channels.
// Currently, it only applies romanization to the lyrics data.
// May be extended in the future.
//
// Every function/method/module/whatever needs to lock the mutex
// by themselves and only themselves.
// No locking a mutex in THIS function.
func Configure(lyricsData *playerStructs.LyricsData) {
	log.Debug("lyrics/configure", "Starting configuring the received lyrics")

	// Traditional to Simplified
	if global.Config.C.Lyrics.T2S {
		if t := getT2S(); t != nil {
			log.Debug("lyrics/configure", "Applying Traditional to Simplified Chinese conversion")
			for i := range lyricsData.Lyrics {
				if converted, err := t.Convert(lyricsData.Lyrics[i].Text); err == nil {
					lyricsData.Lyrics[i].Text = converted
				}
			}
		}
	}

	// Romanization
	log.Debug("lyrics/configure", "Applying romanization if enabled and necessary")
	romanization.Romanize(lyricsData.Lyrics)

	log.Debug("lyrics/configure", "Done")
}

package sync

import (
	"lrcsnc/internal/pkg/log"
	"time"

	"lrcsnc/internal/mpris"
	"lrcsnc/internal/pkg/global"

	mprislib "github.com/Endg4meZer0/go-mpris"
)

var needsSynchronization chan bool = make(chan bool, 1)
var isSynchronizing bool = false

// This is a position synchronizer.
// It is triggered by AskForPositionSync function 
// on any Seeked signals, PlaybackStatus changes and the lyrics fetching.
// 
// It is needed for more precise syncing of the lyrics
// (to prevent a possible mismatch of the position from our data
// and the actual player's position).
//
// To prevent multiple synchronizations from taking place at the same time
// the `needsSynchronization` channel is buffered at 1.
func positionSynchronizer() {
	ticker := time.NewTicker(50 * time.Millisecond) // 0.05 seconds as delta time for sync
	ticker.Stop()                                   // ticker should not fire just yet
	for {
		<-needsSynchronization

		if global.Player.P.PlaybackStatus != mprislib.PlaybackPlaying {
			stopLyricsSync()
			continue
		}

		isSynchronizing = true
		oldPos, err := mpris.GetPosition()
		if err != nil {
			log.Error("sync", "positionSynchronizer failed to get position: " + err.Error())
			continue
		}
		ticker.Reset(50 * time.Millisecond)
		for {
			<-ticker.C
			newPos, err := mpris.GetPosition()
			if err != nil {
			log.Error("sync", "positionSynchronizer failed to get position: " + err.Error())
				break
			}
			if newPos != oldPos {
				global.Player.M.Lock()
				global.Player.P.Position = newPos
				global.Player.M.Unlock()

				break
			}
		}
		ticker.Stop()
		resyncLyrics()
		isSynchronizing = false
	}
}

func AskForPositionSync() {
	if !isSynchronizing {
		needsSynchronization <- true
	}
}

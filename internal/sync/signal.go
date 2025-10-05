package sync

import (
	"fmt"
	"lrcsnc/internal/mpris"
	"lrcsnc/internal/output/pkg/event"
	"lrcsnc/internal/output/server"
	"lrcsnc/internal/pkg/global"
	"lrcsnc/internal/pkg/log"

	mprislib "github.com/Endg4meZer0/go-mpris"
)

func mprisMessageReceiver() {
	for msg := range mpris.MPRISMessageChannel {
		log.Debug("sync/mprisMessageReceiver", fmt.Sprintf("Received message: %v", msg))
		switch msg.Type {
		case mpris.SignalReady, mpris.SignalPlayerChanged:
			go server.ReceiveEvent(event.Event{
				Type: event.EventTypePlayerChanged,
				Data: event.EventTypePlayerChangedData{
					Name: global.Player.P.Name,
				},
			})
			if global.Player.P.PlaybackStatus != mprislib.PlaybackStopped {
				songChanged <- true
			}
			go server.ReceiveEvent(event.Event{
				Type: event.EventTypePlaybackStatusChanged,
				Data: event.EventTypePlaybackStatusChangedData{
					PlaybackStatus: global.Player.P.PlaybackStatus,
				},
			})
		case mpris.SignalSeeked:
			// On seeked signal we only update the position,
			// and we're not sending any specific events
			// to server...
			global.Player.M.Lock()
			global.Player.P.Position = msg.Data.(float64) / 1000 / 1000
			global.Player.M.Unlock()
			// ...and, of course, ask for a position sync
			AskForPositionSync()
		case mpris.SignalPlaybackStatusChanged:
			// If the playback status has changed...
			global.Player.M.Lock()
			global.Player.P.PlaybackStatus = msg.Data.(mprislib.PlaybackStatus)
			// ...and if it's stopped then also reset the position
			if global.Player.P.PlaybackStatus == mprislib.PlaybackStopped {
				global.Player.P.Position = 0
			}

			go server.ReceiveEvent(event.Event{
				Type: event.EventTypePlaybackStatusChanged,
				Data: event.EventTypePlaybackStatusChangedData{
					PlaybackStatus: global.Player.P.PlaybackStatus,
				},
			})

			global.Player.M.Unlock()

			// And ask for a position sync to be sure
			AskForPositionSync()

		case mpris.SignalRateChanged:
			// If the rate has changed...
			global.Player.M.Lock()
			global.Player.P.Rate = msg.Data.(float64)
			go server.ReceiveEvent(event.Event{
				Type: event.EventTypeRateChanged,
				Data: event.EventTypeRateChangedData{
					Rate: global.Player.P.Rate,
				},
			})
			global.Player.M.Unlock()

			// ...we will absolutely ask for a sync, since the old sync is now invalid
			AskForPositionSync()

		case mpris.SignalMetadataChanged:
			// If the metadata has changed...
			global.Player.M.Lock()
			err := mpris.ApplyMetadataOntoGlobal(msg.Data.(mprislib.Metadata))
			if err != nil {
				log.Error("sync/mprisMessageReceiver", "Couldn't parse metadata: "+err.Error())
			}
			global.Player.M.Unlock()

			// Send a signal that the song has changed and we need
			// to fetch some new lyrics
			songChanged <- true

			// If the song changed (and we know it did), the position was probably set to 0
			// and we can also ask for a position sync
			AskForPositionSync()
		}
	}
	log.Error("sync/mprisMessageReceiver", "The MPRIS message channel was closed. What happened?")
}

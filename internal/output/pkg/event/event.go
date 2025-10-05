package event

type Event struct {
	Type EventType
	Data any
}

// EventType represents the type of the event (received or sent)
type EventType string

const (
	EventTypeActiveLyricChanged    EventType = "ActiveLyricChanged"
	EventTypeSongChanged           EventType = "SongChanged"
	EventTypePlayerChanged         EventType = "PlayerChanged"
	EventTypePlaybackStatusChanged EventType = "PlaybackStatusChanged"
	EventTypeRateChanged           EventType = "RateChanged"
	EventTypeLyricsStateChanged    EventType = "LyricsStateChanged"
	EventTypeLyricsChanged         EventType = "LyricsChanged"
	EventTypeOverwriteRequired     EventType = "OverwriteRequired"
	EventTypeServerClosed          EventType = "ServerClosed"
	EventTypeConfigReloaded        EventType = "ConfigReloaded" // only for client
)

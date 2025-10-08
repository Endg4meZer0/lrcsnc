package server

import "lrcsnc/internal/output/pkg/event"

var lastEventTypes []event.EventType = []event.EventType{
	event.EventTypePlayerChanged,
	event.EventTypeRateChanged,
	event.EventTypeSongChanged,
	event.EventTypePlaybackStatusChanged,
	event.EventTypeLyricsStateChanged,
	event.EventTypeLyricsChanged,
	event.EventTypeActiveLyricChanged,
	event.EventTypeOverwriteRequired,
}

func makeDefaultLastEvents() []event.Event {
	events := make([]event.Event, 8)

	for i, et := range lastEventTypes {
		events[i] = event.MakeDefaultEvent(et)
	}

	return events
}

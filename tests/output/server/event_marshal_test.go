package server

import (
	"encoding/json"
	"lrcsnc/internal/output/pkg/event"
	"lrcsnc/internal/pkg/structs"
	"testing"
)

func TestEventMarshal(t *testing.T) {

	lyrCngEvent := event.Event{Type: event.EventTypeActiveLyricChanged, Data: event.EventTypeActiveLyricChangedData{
		Lyric: structs.Lyric{
			Timing: 123.11,
			Text:   "yo~!",
		},
		TimeUntilEnd: 3.11,
		Resync:       false,
	}}

	out, err := json.Marshal(lyrCngEvent)
	if err != nil {
		t.Error("[tests/output/server/TestEventMarshal] ERROR: " + err.Error())
	}

	expected := `{"Type":"ActiveLyricChanged","Data":{"Index":0,"Lyric":{"Timing":123.11,"Text":"yo~!"},"Multiplier":0,"TimeUntilEnd":3.11,"Resync":false}}`
	if string(out) != expected {
		t.Errorf("[tests/output/server/TestEventMarshal] ERROR: wrong output, received %s instead of %s.", string(out), expected)
	}
}

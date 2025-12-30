package player

import (
	"github.com/Endg4meZer0/go-mpris"
)

type Player struct {
	Name           string
	PlaybackStatus mpris.PlaybackStatus
	Position       float64
	Rate           float64
	Song           Song
}

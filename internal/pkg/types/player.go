package types

type PlaybackStatus string

const (
	PlaybackStatusPlaying PlaybackStatus = "Playing"
	PlaybackStatusPaused  PlaybackStatus = "Paused"
	PlaybackStatusStopped PlaybackStatus = "Stopped"
)

package client

import (
	"lrcsnc/internal/output/pkg/event"
	"lrcsnc/internal/pkg/global"
	"lrcsnc/internal/pkg/log"
	"lrcsnc/internal/pkg/types"
	"os"
	"time"
)

func (c *client) handleActiveLyricChanged(d event.EventTypeActiveLyricChangedData) {
	// Set new pending text and index to the received ones.
	c.pendingText = d.Lyric.Text

	if c.pendingText != "" {
		c.pendingText = c.formatLyric(d.Multiplier)
	}

	// If overwrite is not empty yet, once its duration ends,
	// it will restore lyric and index from pending variables.
	// So there is no need to act any further in this case.
	if c.overwrite != "" {
		return
	}

	if c.conn != nil {
		global.Player.M.Lock()
		global.Player.P.LastDetectedPosition = d.Lyric.Timing
		global.Player.M.Unlock()
	}

	c.output()
}

func (c *client) handleSongChanged(d event.EventTypeSongChangedData) {
	if c.conn != nil {
		global.Player.M.Lock()
		global.Player.P.Song.Metadata.Title = d.Title
		global.Player.P.Song.Metadata.Artists = d.Artists
		global.Player.P.Song.Metadata.Album = d.Album
		global.Player.P.Song.Metadata.Duration = d.Duration
		global.Player.M.Unlock()
	}
}

func (c *client) handlePlayerChanged(d event.EventTypePlayerChangedData) {
	if c.conn != nil {
		global.Player.M.Lock()
		global.Player.P.Name = d.Name
		global.Player.M.Unlock()
	}
}

func (c *client) handlePlaybackStatusChanged(d event.EventTypePlaybackStatusChangedData) {
	if c.conn != nil {
		global.Player.M.Lock()
		global.Player.P.PlaybackStatus = d.PlaybackStatus
		global.Player.M.Unlock()
	}

	if global.Player.P.PlaybackStatus == types.PlaybackStatusStopped {
		c.pendingText = ""
	}

	c.output()
}

func (c *client) handleRateChanged(d event.EventTypeRateChangedData) {
	if c.conn != nil {
		global.Player.M.Lock()
		global.Player.P.Rate = d.Rate
		global.Player.M.Unlock()
	}
}

func (c *client) handleLyricsStateChanged(d event.EventTypeLyricsStateChangedData) {
	if c.conn != nil {
		global.Player.M.Lock()
		global.Player.P.Song.LyricsData.LyricsState = d.State
		global.Player.M.Unlock()
	}

	c.pendingText = ""
	c.output()
}

func (c *client) handleLyricsChanged(d event.EventTypeLyricsChangedData) {
	if c.conn != nil {
		global.Player.M.Lock()
		global.Player.P.Song.LyricsData.Lyrics = d.Lyrics
		global.Player.M.Unlock()
	}

	c.pendingText = ""
	c.output()
}

func (c *client) handleOverwriteRequired(d event.EventTypeOverwriteRequiredData) {
	c.pendingText = c.activeText
	c.overwrite = d.Overwrite
	c.output()

	go func() {
		<-time.After(5 * time.Second)
		c.overwrite = ""
		c.output()
	}()
}

func (c *client) handleServerClosed(_ event.EventTypeServerClosedData) {
	c.close()
	log.Info("output/client", "Received ServerClosed event, disconnecting and stopping client...")
	os.Exit(1)
}

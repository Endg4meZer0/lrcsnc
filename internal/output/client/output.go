package client

import (
	"os"
	"strings"
	"time"

	"lrcsnc/internal/pkg/global"
	"lrcsnc/internal/pkg/log"
	"lrcsnc/internal/pkg/types"

	mprislib "github.com/Endg4meZer0/go-mpris"
)

// output sets the client's pendingText as activeText and
// writes it to outputDestination. If the outputDestination
// is not related to std, then does its best to ensure the
// output is an atomic operation by using temp files.
func (c *Client) output() {
	if c.pendingText == "" {
		if !c.instrumentalActive {
			c.instrumentalActive = true
			go c.instrumentalLoop()
		}
		return
	}

	if c.overwrite != "" {
		c.instrumentalActive = false
		c.activeText = c.overwrite
	} else {
		c.instrumentalActive = false
		c.activeText = c.pendingText
	}

	c.write()
}

func (c *Client) write() {
	global.Config.M.Lock()
	global.Player.M.Lock()

	s := c.formatToTemplate()
	if global.Config.C.ClientOutput.InsertNewline {
		s = s + "\n"
	}

	global.Player.M.Unlock()
	global.Config.M.Unlock()

	if !c.isOutputStd() {
		if tempDestination, err := os.OpenFile(c.tempFile, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644); err == nil {
			// Atomic copy (for better support of something like obs-text-pthread)
			tempDestination.Truncate(0)
			tempDestination.Seek(0, 0)
			tempDestination.WriteString(s)
			err = os.Rename(c.tempFile, c.outputPath)
			if err != nil {
				log.Error("output/client", "Failed to move the temp file onto output destination: "+err.Error())
				return
			}
		} else {
			// If temp destination is unavailable, revert to basic handling (not atomic)
			c.outputDestination.Truncate(0)
			c.outputDestination.Seek(0, 0)
			c.outputDestination.WriteString(s)
		}
	} else {
		c.outputDestination.WriteString(s)
	}
}

// instrumentalLoop starts a loop of writing instrumental notes
// until the next lyric shows up.
func (c *Client) instrumentalLoop() {
	i := 1
	instrumentalGen := func() (out string, stopWritingAfter bool) {
		global.Config.M.Lock()
		global.Player.M.Lock()

		note := global.Config.C.ClientOutput.Format.Instrumental.Symbol
		j := int(global.Config.C.ClientOutput.Format.Instrumental.MaxSymbols + 1)

		// Only update instrumental stuff if there is an active song
		if global.Player.P.PlaybackStatus != mprislib.PlaybackStopped {
			var stringToPrint string

			switch global.Player.P.Song.LyricsData.LyricsState {
			case types.LyricsStateSynced, types.LyricsStateInstrumental:
				stringToPrint = strings.NewReplacer("%lyric%", "", "%multiplier%", "").Replace(global.Config.C.ClientOutput.Format.Lyric)
			case types.LyricsStatePlain:
				stringToPrint = global.Config.C.ClientOutput.Format.NoSyncedLyrics
			case types.LyricsStateNotFound:
				stringToPrint = global.Config.C.ClientOutput.Format.NoLyrics
			case types.LyricsStateLoading:
				stringToPrint = global.Config.C.ClientOutput.Format.LoadingLyrics
			default:
				stringToPrint = global.Config.C.ClientOutput.Format.ErrorMessage
			}

			if len(stringToPrint) != 0 {
				stringToPrint += " "
			}
			stringToPrint += strings.Repeat(note, i%j)

			out = stringToPrint
			stopWritingAfter = global.Player.P.PlaybackStatus == mprislib.PlaybackPaused

			i++
			if i >= j {
				i = 1
			}
		} else {
			out = global.Config.C.ClientOutput.Format.NotPlaying
			stopWritingAfter = true
		}
		global.Player.M.Unlock()
		global.Config.M.Unlock()

		return
	}

	// A preemptive c.write() 'cause new tickers always wait
	// their duration before firing a first tick
	ig, stopWritingAfter := instrumentalGen()
	c.activeText = ig
	c.write()
	if stopWritingAfter {
		c.instrumentalActive = false
		return
	}

	instrumentalTicker := time.NewTicker(time.Duration(global.Config.C.ClientOutput.Format.Instrumental.Interval*1000) * time.Millisecond)
	for range instrumentalTicker.C {
		if !c.instrumentalActive {
			instrumentalTicker.Stop()
			break
		}

		ig, stopWritingAfter := instrumentalGen()
		c.activeText = ig
		c.write()
		if stopWritingAfter {
			c.instrumentalActive = false
			break
		}
	}
}

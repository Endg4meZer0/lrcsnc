package client

import (
	"fmt"
	"lrcsnc/internal/pkg/global"
	"lrcsnc/internal/pkg/log"
	"lrcsnc/internal/pkg/util/dynreplacer"
	"os"
	"strconv"
	"strings"
)

func (c *Client) isOutputStd() bool {
	return c.outputDestination == os.Stdout ||
		c.outputDestination == os.Stderr ||
		c.outputDestination == os.Stdin
}

// changeOutput changes the output destination to the specified path.
// The write check is usually performed at config validation step,
// but it's good to have it here too.
func (c *Client) changeOutput() error {
	newDest, err := os.OpenFile(global.Config.C.ClientOutput.Destination, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		log.Error("output/client", "Error opening new output destination file, falling back to stdout. More: "+err.Error())
		global.Config.C.ClientOutput.Destination = c.outputPath
		return err
	} else {
		r := strings.NewReplacer(
			"{pid}", strconv.Itoa(os.Getpid()),
		)
		c.tempFile = r.Replace(global.Config.C.ClientOutput.Destination + ".{pid}.tmp")
		// We'll try to open temp file here to see if it even works
		_, err := os.OpenFile(c.tempFile, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
		if err != nil {
			log.Error("output/client", fmt.Sprintf("Failed to open a temp write file (%s). The writes may not be atomic.", c.tempFile))
		}

		c.outputDestination = newDest
		c.outputPath = global.Config.C.ClientOutput.Destination
	}

	return nil
}

func (c *Client) setTemplateReplacer() {
	// The mutexes should be locked before replacer takes action.
	c.templateReplacer = dynreplacer.NewDynamicReplacer(
		map[string]func() string{
			"text": func() string { return c.activeText },
			"artist": func() string {
				if len(global.Player.P.Song.Artists) > 0 {
					return global.Player.P.Song.Artists[0]
				} else {
					return ""
				}
			},
			"artists": func() string {
				return strings.Join(global.Player.P.Song.Artists, ", ")
			},
			"title": func() string {
				return global.Player.P.Song.Title
			},
			"album": func() string {
				return global.Player.P.Song.Album
			},
			"position": func() string {
				return fmt.Sprintf("%02d:%02d", int(global.Player.P.Position)/60, int(global.Player.P.Position)%60)
			},
			"duration": func() string {
				return fmt.Sprintf("%02d:%02d", int(global.Player.P.Song.Duration)/60, int(global.Player.P.Song.Duration)%60)
			},
			"rate": func() string {
				return fmt.Sprintf("%1.2fx", global.Player.P.Rate)
			},
			"playback-status": func() string {
				return strings.ToLower(string(global.Player.P.PlaybackStatus))
			},
			"lyrics-status": func() string {
				return global.Player.P.Song.LyricsData.LyricsState.String()
			},
			"player-name": func() string {
				return global.Player.P.Name
			},
		},
	)
}

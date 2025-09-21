package piped

import (
	"fmt"
	"lrcsnc/internal/pkg/global"
	"strings"
)

// FormatToTemplate formats the output to the Template configuration option.
// Does NOT lock player's mutex, so make sure it is locked beforehand.
//
// text - the outcoming text meant for display (e.g. lyric or error message)
func FormatToTemplate(text string) (out string) {
	var artist string
	if len(global.Player.P.Song.Artists) > 0 {
		artist = global.Player.P.Song.Artists[0]
	}

	formatReplacer := strings.NewReplacer(
		"{text}", text,
		"{artist}", artist,
		"{artists}", strings.Join(global.Player.P.Song.Artists, ", "),
		"{title}", global.Player.P.Song.Title,
		"{album}", global.Player.P.Song.Album,
		"{position}", fmt.Sprintf("%02d:%02d", int(global.Player.P.Position)/60, int(global.Player.P.Position)%60),
		"{duration}", fmt.Sprintf("%02d:%02d", int(global.Player.P.Song.Duration)/60, int(global.Player.P.Song.Duration)%60),
		"{rate}", fmt.Sprintf("%1.2fx", global.Player.P.Rate),
		"{playback-status}", strings.ToLower(string(global.Player.P.PlaybackStatus)),
		"{lyrics-status}", strings.ToLower(global.Player.P.Song.LyricsData.LyricsState.String()),
		"{player-name}", global.Player.P.Name,
	)

	return formatReplacer.Replace(global.Config.C.Output.Piped.Template)
}

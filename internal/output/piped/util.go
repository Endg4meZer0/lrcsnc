package piped

import (
	"lrcsnc/internal/pkg/structs"
)

func lyricIndexToString(l int, lyricsData []structs.Lyric) string {
	if l < 0 || l >= len(lyricsData) {
		return ""
	} else {
		return lyricsData[l].Text
	}
}

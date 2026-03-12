package lrc

import (
	lyricStruct "lrcsnc/internal/lyrics/struct"
	"strings"
)

type LyricsFormatLrcPlain struct{}

func (LyricsFormatLrcPlain) Convert(data string) lyricStruct.Lyrics {
	lines := strings.Split(data, "\n")
	out := make(lyricStruct.Lyrics, 0, len(lines))
	for _, lyric := range lines {
		out = append(out, lyricStruct.Lyric{
			Timing: 0,
			Text:   strings.TrimSpace(lyric),
		})
	}
	return out
}

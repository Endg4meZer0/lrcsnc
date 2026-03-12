package formats

import lyricStruct "lrcsnc/internal/lyrics/struct"

type LyricsFormat interface {
	Convert(data string) lyricStruct.Lyrics
}

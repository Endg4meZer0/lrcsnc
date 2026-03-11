package playerStruct

import lyricStruct "lrcsnc/internal/lyrics/struct"

type Song struct {
	Metadata   SongMetadata
	LyricsData lyricStruct.LyricsData
}

type SongMetadata struct {
	Title        string
	Artists      []string
	Album        string
	AlbumArtists []string
	AlbumArt     string // link to the album art, see xesam:artUrl
	Url          string // where the file is located, see xesam:url
	Duration     float64
}

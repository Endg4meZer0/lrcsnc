package player

type Song struct {
	Title        string
	Artists      []string
	Album        string
	AlbumArtists []string
	Duration     float64
	LyricsData   LyricsData
}

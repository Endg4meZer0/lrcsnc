package player

import (
	"hash/fnv"
	"strconv"
	"strings"
)

type Song struct {
	Title        string
	Artists      []string
	Album        string
	AlbumArtists []string
	Duration     float64
	LyricsData   LyricsData
}

func (s *Song) ID() uint64 {
	h := fnv.New64a()
	h.Write([]byte(s.Title))
	h.Write([]byte(strings.Join(s.Artists, ", ")))
	h.Write([]byte(s.Album))
	h.Write([]byte(strconv.FormatFloat(s.Duration, 'f', 1, 64)))
	return h.Sum64()
}

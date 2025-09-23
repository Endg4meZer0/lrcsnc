package structs

import (
	"hash/fnv"
	"lrcsnc/internal/pkg/types"
	"strconv"
	"strings"

	"github.com/Endg4meZer0/go-mpris"
)

type Player struct {
	Name           string
	PlaybackStatus mpris.PlaybackStatus
	Position       float64
	Rate           float64
	Song           Song
}

type Song struct {
	Title      string
	Artists    []string
	Album      string
	Duration   float64
	LyricsData LyricsData
}

type LyricsData struct {
	Lyrics      []Lyric
	LyricsState types.LyricsState
}

type Lyric struct {
	Timing float64
	Text string
}

func (s *Song) ID() uint64 {
	h := fnv.New64a()
	h.Write([]byte(s.Title))
	h.Write([]byte(strings.Join(s.Artists, ", ")))
	h.Write([]byte(s.Album))
	h.Write([]byte(strconv.FormatFloat(s.Duration, 'f', 1, 64)))
	return h.Sum64()
}

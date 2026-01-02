package lrclib

import (
	errs "lrcsnc/internal/lyrics/errors"
	"lrcsnc/internal/pkg/log"
	playerStructs "lrcsnc/internal/pkg/structs/player"
	"lrcsnc/internal/pkg/types"
	"strings"
)

func (l Provider) Get(song playerStructs.Song) (playerStructs.LyricsData, error) {
	var body []byte
	var err error
	var res playerStructs.LyricsData

	log.Debug("lyrics/providers/lrclib/get", "Trying to fetch lyrics...")
	body, err = requestLyrics(song.Title, strings.Join(song.AlbumArtists, ", "))
	if err == nil {
		res, err = responseListToLyricsData(&song, body)
	}
	if err != errs.NotFound {
		return res, err
	}

	log.Debug("lyrics/providers/lrclib/get", "Failed; the lyrics for this song don't exist")

	// If nothing is found, return a not found state
	return playerStructs.LyricsData{LyricsState: types.LyricsStateNotFound}, errs.NotFound
}

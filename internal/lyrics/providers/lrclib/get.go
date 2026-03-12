package lrclib

import (
	errs "lrcsnc/internal/lyrics/errors"
	lyricStruct "lrcsnc/internal/lyrics/struct"
	"lrcsnc/internal/pkg/global"
	"lrcsnc/internal/pkg/log"
	"lrcsnc/internal/pkg/types"
	playerStruct "lrcsnc/internal/player/struct"
	"strings"
)

func (l Provider) Get(song playerStruct.Song) (lyricStruct.LyricsData, error) {
	var body []byte
	var err error
	var res lyricStruct.LyricsData = lyricStruct.LyricsData{LyricsState: types.LyricsStateNotFound}

	if global.Config.C.Lyrics.LrcLibProviderConfig.EnableFirstGetRequest {
		log.Debug("lyrics/providers/lrclib/get", "Trying to get lyrics directly...")
		body, err = getLyrics(song.Metadata.Title, strings.Join(song.Metadata.Artists, ", "), song.Metadata.Album, song.Metadata.Duration)
		if err == nil {
			outs, err := parseResps(body)
			if err == nil && outs[0].toLyricsData().LyricsState != types.LyricsStatePlain {
				return outs[0].toLyricsData(), nil
			}
		}
		log.Debug("lyrics/providers/lrclib/get", "Failed to get lyrics directly")
	}

	log.Debug("lyrics/providers/lrclib/get", "Trying to search for lyrics broadly...")
	body, err = searchLyrics(song.Metadata.Title, strings.Join(song.Metadata.AlbumArtists, ", "))
	if err == nil {
		res, err = responseListToLyricsData(&song, body)
	}
	if err != errs.NotFound {
		return res, err
	}

	log.Debug("lyrics/providers/lrclib/get", "Failed; the lyrics for this song don't exist")

	// If nothing is found, return a not found state
	return res, errs.NotFound
}

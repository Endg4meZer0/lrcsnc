package lrclib

import (
	"net/url"

	"lrcsnc/internal/pkg/errors"
	"lrcsnc/internal/pkg/log"
	playerStructs "lrcsnc/internal/pkg/structs/player"
	"lrcsnc/internal/pkg/types"
)

func (l Provider) Get(song playerStructs.Song) (playerStructs.LyricsData, error) {
	var getURL *url.URL
	var body []byte
	var err error
	var res playerStructs.LyricsData

	log.Debug("lyrics/providers/lrclib/Get", "Trying to fetch lyrics with a /get request full with details")

	// Try to get the lyrics with everything exact: artists, album, duration
	if song.Duration != 0 {
		getURL = makeURL(song, lrcLibURLTypeGet)
		body, err = sendRequest(getURL)
	}
	if err == nil {
		res, err = dtoListToLyricsData(song, body)
	}
	if err != errors.ErrLyricsNotFound {
		return res, err
	}

	// If the song has more than one artist,
	// try to get the lyrics with everything exact except pick only the first artist
	if len(song.Artists) > 1 {
		log.Debug("lyrics/providers/lrclib/Get", "Failed; trying to fetch lyrics with a /get request with all details except pick only the first artist")

		getURL = makeURL(song, lrcLibURLTypeGetWithSingleArtist)
		body, err = sendRequest(getURL)
		if err == nil {
			res, err = dtoListToLyricsData(song, body)
		}
		if err != errors.ErrLyricsNotFound {
			return res, err
		}
	}

	log.Debug("lyrics/providers/lrclib/Get", "Failed; trying to fetch lyrics with a /search request with all details")

	// Try to search for lyrics with exact album and artists
	getURL = makeURL(song, lrcLibURLTypeSearchWithAlbum)
	body, err = sendRequest(getURL)
	if err == nil {
		res, err = dtoListToLyricsData(song, body)
	}
	if err != errors.ErrLyricsNotFound {
		return res, err
	}

	// If the song has more than one artist,
	// try to search for lyrics with exact album, but only the first artist
	if len(song.Artists) > 1 {
		log.Debug("lyrics/providers/lrclib/Get", "Failed; trying to fetch lyrics with a /search request with all details except pick only the first artist")
		getURL = makeURL(song, lrcLibURLTypeSearchWithSingleArtistAndAlbum)
		body, err = sendRequest(getURL)
		if err == nil {
			res, err = dtoListToLyricsData(song, body)
		}
		if err != errors.ErrLyricsNotFound {
			return res, err
		}
	}

	log.Debug("lyrics/providers/lrclib/Get", "Failed; trying to fetch lyrics with a /search request without album")

	// Try to search for lyrics with only the title and all artists
	getURL = makeURL(song, lrcLibURLTypeSearch)
	body, err = sendRequest(getURL)
	if err == nil {
		res, err = dtoListToLyricsData(song, body)
	}
	if err != errors.ErrLyricsNotFound {
		return res, err
	}

	// If the song has more than one artist,
	// try to search for lyrics with only the title and the first artist
	if len(song.Artists) > 1 {
		log.Debug("lyrics/providers/lrclib/Get", "Failed; trying to fetch lyrics with a /search request without album and picking only the first artist")
		getURL = makeURL(song, lrcLibURLTypeSearchWithSingleArtist)
		body, err = sendRequest(getURL)
		if err == nil {
			res, err = dtoListToLyricsData(song, body)
		}
		if err != errors.ErrLyricsNotFound {
			return res, err
		}
	}

	log.Debug("lyrics/providers/lrclib/Get", "Failed; the lyrics for this song don't exist")

	// If nothing is found, return a not found state
	return playerStructs.LyricsData{LyricsState: types.LyricsStateNotFound}, errors.ErrLyricsNotFound
}

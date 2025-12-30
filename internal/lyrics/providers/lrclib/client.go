package lrclib

import (
	"fmt"
	"io"
	"math"
	"net/http"
	"net/url"
	"strings"

	"lrcsnc/internal/pkg/errors"
	"lrcsnc/internal/pkg/log"
	playerStructs "lrcsnc/internal/pkg/structs/player"
)

type lrcLibURLType int

const (
	lrcLibURLTypeGet lrcLibURLType = iota
	lrcLibURLTypeGetWithSingleArtist
	lrcLibURLTypeSearchWithAlbum
	lrcLibURLTypeSearchWithSingleArtistAndAlbum
	lrcLibURLTypeSearch
	lrcLibURLTypeSearchWithSingleArtist
)

func makeURL(song playerStructs.Song, t lrcLibURLType) (out *url.URL) {
	rawURL := "http://lrclib.net/api/"
	switch t {
	case lrcLibURLTypeGet:
		rawURL += "get?" + url.PathEscape(fmt.Sprintf("track_name=%v&artist_name=%v&album_name=%v&duration=%v", song.Title, strings.Join(song.Artists, ", "), song.Album, int(math.Ceil(song.Duration))))
	case lrcLibURLTypeGetWithSingleArtist:
		rawURL += "get?" + url.PathEscape(fmt.Sprintf("track_name=%v&artist_name=%v&album_name=%v&duration=%v", song.Title, song.Artists[0], song.Album, int(math.Ceil(song.Duration))))
	case lrcLibURLTypeSearchWithAlbum:
		rawURL += "search?" + url.PathEscape(fmt.Sprintf("track_name=%v&artist_name=%v&album_name=%v", song.Title, strings.Join(song.Artists, ", "), song.Album))
	case lrcLibURLTypeSearchWithSingleArtistAndAlbum:
		rawURL += "search?" + url.PathEscape(fmt.Sprintf("track_name=%v&artist_name=%v&album_name=%v", song.Title, song.Artists[0], song.Album))
	case lrcLibURLTypeSearch:
		rawURL += "search?" + url.PathEscape(fmt.Sprintf("track_name=%v&artist_name=%v", song.Title, strings.Join(song.Artists, ", ")))
	case lrcLibURLTypeSearchWithSingleArtist:
		rawURL += "search?" + url.PathEscape(fmt.Sprintf("track_name=%v&artist_name=%v", song.Title, song.Artists[0]))
	default:
		return nil
	}
	out, err := url.Parse(rawURL)
	if err != nil {
		log.Error("lyrics/providers/lrclib/makeURL", fmt.Sprintf("Failed to parse string (%v) to URL", rawURL))
	}
	return
}

func sendRequest(link *url.URL) ([]byte, error) {
	resp, err := http.Get((*link).String())
	if resp.StatusCode == 404 {
		return nil, errors.ErrLyricsNotFound
	}
	if err != nil || resp.StatusCode != 200 {
		return nil, errors.ErrLyricsServerError
	}

	body, err := io.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		return nil, errors.ErrLyricsBodyReadFail
	}
	return body, nil
}

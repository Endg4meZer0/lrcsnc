package local

import (
	errs "lrcsnc/internal/lyrics/errors"
	"lrcsnc/internal/lyrics/formats"
	lyricStruct "lrcsnc/internal/lyrics/struct"
	"lrcsnc/internal/pkg/global"
	"lrcsnc/internal/pkg/log"
	"lrcsnc/internal/pkg/types"
	playerStruct "lrcsnc/internal/player/struct"
	"os"
	"strings"
)

func (l Provider) Get(song playerStruct.Song) (lyricStruct.LyricsData, error) {
	var data string
	var err error
	var res lyricStruct.LyricsData = lyricStruct.LyricsData{LyricsState: types.LyricsStateNotFound}

	if song.Metadata.Url == "" ||
		strings.HasPrefix(song.Metadata.Url, "http://") ||
		strings.HasPrefix(song.Metadata.Url, "https://") {
		log.Debug("lyrics/providers/local/get", "Skipping local provider; the song is played online or has no URL field in metadata.")
		return res, errs.NotFound
	}

	songPath, _ := strings.CutPrefix(song.Metadata.Url, "file://")

	// this thing needs to be changed wth i'm such a bad coder
	if global.Config.C.Lyrics.LocalProviderConfig.TryEmbeddedFirst {
		log.Debug("lyrics/providers/local/get", "Trying to get embedded lyrics...")
		data, err = getFromTags(songPath)
		if err == nil && len(data) != 0 {
			format := formats.DetectFormat(data)
			res.Lyrics = formats.Formats[format].Convert(data)
			res.LyricsState = format.ToLyricsState()
			return res, nil
		} else if err == nil {
			log.Debug("lyrics/providers/local/get", "Failed to get embedded lyrics: no lyrics found (got empty string).")
		} else {
			log.Debug("lyrics/providers/local/get", "Failed to get embedded lyrics: "+err.Error())
		}
	}

	log.Debug("lyrics/providers/local/get", "Trying to get lyrics from a nearby LRC file...")
	data, err = getFromLrcFile(songPath)
	if err == nil && len(data) != 0 {
		format := formats.DetectFormat(data)
		res.Lyrics = formats.Formats[format].Convert(data)
		res.LyricsState = format.ToLyricsState()
		return res, nil
	} else if err == nil || os.IsNotExist(err) {
		log.Debug("lyrics/providers/local/get", "Failed to get lyrics from LRC file: no lyrics found (got empty string).")
	} else {
		log.Debug("lyrics/providers/local/get", "Failed to get lyrics from LRC file: "+err.Error())
	}

	if !global.Config.C.Lyrics.LocalProviderConfig.TryEmbeddedFirst {
		log.Debug("lyrics/providers/local/get", "Trying to get embedded lyrics...")
		data, err = getFromTags(songPath)
		if err == nil && len(data) != 0 {
			format := formats.DetectFormat(data)
			res.Lyrics = formats.Formats[format].Convert(data)
			res.LyricsState = format.ToLyricsState()
			return res, nil
		} else if err == nil {
			log.Debug("lyrics/providers/local/get", "Failed to get embedded lyrics: no lyrics found (got empty string).")
		} else {
			log.Debug("lyrics/providers/local/get", "Failed to get embedded lyrics: "+err.Error())
		}
	}

	log.Debug("lyrics/providers/local/get", "Failed; the lyrics for this song don't exist locally")

	// If nothing is found, return a not found state
	return res, errs.NotFound
}

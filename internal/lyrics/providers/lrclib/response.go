package lrclib

import (
	"encoding/json"
	"math"
	"strings"

	errs "lrcsnc/internal/lyrics/errors"
	"lrcsnc/internal/lyrics/formats"
	lyricStruct "lrcsnc/internal/lyrics/struct"
	"lrcsnc/internal/pkg/errors"
	"lrcsnc/internal/pkg/types"
	playerStruct "lrcsnc/internal/player/struct"
)

type LrcLibResponse struct {
	Title        string  `json:"trackName"`
	Artist       string  `json:"artistName"`
	Album        string  `json:"albumName"`
	Duration     float64 `json:"duration"`
	Instrumental bool    `json:"instrumental"`
	PlainLyrics  string  `json:"plainLyrics"`
	SyncedLyrics string  `json:"syncedLyrics"`
}

type LrcLibResponses []LrcLibResponse

func responseListToLyricsData(song *playerStruct.Song, bytes []byte) (lyricStruct.LyricsData, error) {
	resps, err := parseResps(bytes)
	if err != nil {
		return lyricStruct.LyricsData{LyricsState: types.LyricsStateUnknown}, err
	}

	if len(resps) == 0 {
		return lyricStruct.LyricsData{LyricsState: types.LyricsStateNotFound}, errs.NotFound
	}

	resp := resps.pickBest(song)
	if resp.Title != "" && resp.Artist != "" {
		lyricsData := resp.toLyricsData()
		return lyricsData, nil
	}

	return lyricStruct.LyricsData{LyricsState: types.LyricsStateNotFound}, errs.NotFound
}

func parseResps(data []byte) (LrcLibResponses, error) {
	var out LrcLibResponse
	err := json.Unmarshal(data, &out)
	if err != nil {
		var outs LrcLibResponses
		err = json.Unmarshal(data, &outs)
		if err != nil {
			return nil, errors.UnmarshalFail
		}
		return outs, nil
	}

	return LrcLibResponses{out}, nil
}

func (resp LrcLibResponse) toLyricsData() (out lyricStruct.LyricsData) {
	if !resp.Instrumental && resp.PlainLyrics == "" && resp.SyncedLyrics == "" {
		out.LyricsState = types.LyricsStateUnknown
	}

	if resp.Instrumental {
		out.LyricsState = types.LyricsStateInstrumental
		return
	}

	if resp.PlainLyrics != "" && resp.SyncedLyrics == "" {
		out.Lyrics = formats.Formats[types.LyricsFormatLrcPlain].Convert(resp.PlainLyrics)
		out.LyricsState = types.LyricsStatePlain
		return
	}

	out.Lyrics = formats.Formats[types.LyricsFormatLrcSynced].Convert(resp.SyncedLyrics)
	out.LyricsState = types.LyricsStateSynced

	return
}

func (rs LrcLibResponses) pickBest(song *playerStruct.Song) LrcLibResponse {
	if len(rs) == 0 {
		return LrcLibResponse{}
	}

	if len(rs) == 1 {
		return rs[0]
	}

	// Picking best uses a score system
	// I implemented various checks that increase score for lyrics data.
	// Here's the list (in order as in code):
	// - Having only plain lyrics gets 1 point;
	// 	 having instrumental mark gets 3 points;
	// 	 having synced lyrics gets 5 points
	// - A full match on album name gets 2 points; having a part matching gets 1 point
	// - A full match on artists list gets 3 points; having only the first gets 1 point
	//   The original (first) artist is determined using album-artist field, if possible
	// - A direct duration match (while rounding up) gets 5 points;
	//   a duration match within 2 seconds of difference gets 3 points;
	//   a duration match within 5 seconds of difference gets 1 point.

	maxi := 0
	maxs := 0
	for i, r := range rs {
		score := 0

		// -----
		if r.Instrumental {
			score += 3
		} else if r.SyncedLyrics != "" {
			score += 5
		} else if r.PlainLyrics != "" {
			score += 1
		}

		// -----
		if r.Album == song.Metadata.Album {
			score += 2
		} else if strings.Contains(song.Metadata.Album, r.Album) || strings.Contains(r.Album, song.Metadata.Album) {
			score += 1
		}

		// -----
		artistNameCleaner := strings.NewReplacer(" & ", " ", ", ", " ")
		rartist := artistNameCleaner.Replace(r.Artist)
		sartist := artistNameCleaner.Replace(strings.Join(song.Metadata.Artists, ", "))
		salbartist := artistNameCleaner.Replace(strings.Join(song.Metadata.AlbumArtists, ", "))

		switch rartist {
		case sartist:
			score += 3
		case salbartist:
			score += 1
		}

		// -----
		if math.Floor(r.Duration) == math.Floor(song.Metadata.Duration) {
			score += 5
		} else if math.Abs(r.Duration-song.Metadata.Duration) <= 2 {
			score += 3
		} else if math.Abs(r.Duration-song.Metadata.Duration) <= 5 {
			score += 1
		}

		// ---------
		if score > maxs {
			maxs = score
			maxi = i
		}
	}

	return rs[maxi]
}

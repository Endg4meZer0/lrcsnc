package lrclib

import (
	"cmp"
	"encoding/json"
	"math"
	"regexp"
	"slices"
	"strconv"
	"strings"

	"lrcsnc/internal/pkg/errors"
	playerStructs "lrcsnc/internal/pkg/structs/player"
	"lrcsnc/internal/pkg/types"
)

var timingRegexp = regexp.MustCompile(`(\[[0-9]{2}:[0-9]{2}.[0-9]{1,3}])+`)

type DTO struct {
	Title        string  `json:"trackName"`
	Artist       string  `json:"artistName"`
	Album        string  `json:"albumName"`
	Duration     float64 `json:"duration"`
	Instrumental bool    `json:"instrumental"`
	PlainLyrics  string  `json:"plainLyrics"`
	SyncedLyrics string  `json:"syncedLyrics"`
}

func dtoListToLyricsData(song playerStructs.Song, bytes []byte) (playerStructs.LyricsData, error) {
	dtos, err := parseDTOs(bytes)
	if err != nil {
		return playerStructs.LyricsData{LyricsState: types.LyricsStateUnknown}, err
	}

	if len(dtos) == 0 {
		return playerStructs.LyricsData{LyricsState: types.LyricsStateNotFound}, errors.ErrLyricsNotFound
	}

	dtos = removeMismatches(song, dtos)
	if len(dtos) != 0 {
		lyricsData := dtos[0].toLyricsData()
		return lyricsData, nil
	}

	return playerStructs.LyricsData{LyricsState: types.LyricsStateNotFound}, errors.ErrLyricsNotFound
}

func parseDTOs(data []byte) ([]DTO, error) {
	var out DTO
	err := json.Unmarshal(data, &out)
	if err != nil {
		var outs []DTO
		err = json.Unmarshal(data, &outs)
		if err != nil {
			return nil, errors.ErrUnmarshalFail
		}
		return outs, nil
	}

	return []DTO{out}, nil
}

func (dto DTO) toLyricsData() (out playerStructs.LyricsData) {
	if !dto.Instrumental && dto.PlainLyrics == "" && dto.SyncedLyrics == "" {
		out.LyricsState = types.LyricsStateUnknown
	}

	if dto.Instrumental {
		out.LyricsState = types.LyricsStateInstrumental
		return
	}

	if dto.PlainLyrics != "" && dto.SyncedLyrics == "" {
		lyrics := strings.Split(dto.PlainLyrics, "\n")
		out.Lyrics = make([]playerStructs.Lyric, len(lyrics))
		for i, l := range lyrics {
			out.Lyrics[i] = playerStructs.Lyric{Text: sanitizeLyric(l)}
		}

		out.LyricsState = types.LyricsStatePlain
		return
	}

	out.Lyrics = parseSyncedLyrics(dto.SyncedLyrics)
	out.LyricsState = types.LyricsStateSynced

	return
}

func parseSyncedLyrics(lyrics string) (out []playerStructs.Lyric) {
	hasRepetitiveLyrics := false
	syncedLyrics := strings.Split(lyrics, "\n")

	out = make([]playerStructs.Lyric, 0, len(syncedLyrics))

	for _, lyric := range syncedLyrics {
		timings := timingRegexp.FindAllString(lyric, -1)

		for _, ts := range timings {
			lyric = strings.Replace(lyric, ts, "", 1)
		}
		lyric = sanitizeLyric(lyric)

		hasRepetitiveLyrics = hasRepetitiveLyrics || len(timings) > 1

		for _, timingStr := range timings {
			timing := parseTiming(timingStr)
			if timing == -1 {
				continue
			}
			out = append(out, playerStructs.Lyric{
				Timing: timing,
				Text:   lyric,
			})
		}
	}

	if hasRepetitiveLyrics {
		slices.SortFunc(out, func(i, j playerStructs.Lyric) int {
			return cmp.Compare(i.Timing, j.Timing)
		})
	}

	return
}

// A simple sanitize requires trimming any carriage return and space symbols
func sanitizeLyric(lyric string) string {
	return strings.TrimSpace(strings.TrimRight(lyric, "\r"))
}

// Returns the timestamp in seconds, specified in the provided timeTag
func parseTiming(timeTag string) float64 {
	// [01:23.45]
	if len(timeTag) != 10 {
		return -1
	}
	minutes, err := strconv.ParseFloat(timeTag[1:3], 64)
	if err != nil {
		return -1
	}
	seconds, err := strconv.ParseFloat(timeTag[4:9], 64)
	if err != nil {
		return -1
	}
	return minutes*60.0 + seconds
}

func removeMismatches(song playerStructs.Song, dtos []DTO) []DTO {
	if len(dtos) == 0 {
		return dtos
	}

	var matchingLyricsData []DTO = make([]DTO, 0, len(dtos))

	for _, dto := range dtos {
		if strings.EqualFold(dto.Title, song.Title) &&
			// If player doesn't provide the song's duration, ignore the duration check
			// Otherwise, do a check that prevents different versions of a song of messing up the response
			((song.Duration != 0) == (math.Abs(float64(dto.Duration)-song.Duration) <= 2)) {
			matchingLyricsData = append(matchingLyricsData, dto)
		}
	}

	return matchingLyricsData
}

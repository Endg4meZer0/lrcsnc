package lrc

import (
	"cmp"
	playerStructs "lrcsnc/internal/pkg/structs/player"
	"regexp"
	"slices"
	"strconv"
	"strings"
)

var timingRegexp = regexp.MustCompile(`(\[[0-9]{2}:[0-9]{2}.[0-9]{1,3}])+`)

// ConvertPlain is pretty much self-explanatory.
//
// It just splits the provided data string by the new line symbol,
// then constructs a Lyrics object out of that.
func ConvertPlain(data string) playerStructs.Lyrics {
	lines := strings.Split(data, "\n")
	out := make(playerStructs.Lyrics, 0, len(lines))
	for _, lyric := range lines {
		out = append(out, playerStructs.Lyric{
			Timing: 0,
			Text:   strings.TrimSpace(lyric),
		})
	}
	return out
}

// ConvertSynced is self-explanatory as well.
//
// It splits the provided data string by the new line symbol,
// then tries to use known LRC practices for synced lyrics
// to construct the Lyrics object out of that.
func ConvertSynced(data string) playerStructs.Lyrics {
	hasRepetitiveLyrics := false
	lines := strings.Split(data, "\n")
	startIndex := 0
	for _, l := range lines {
		if !timingRegexp.MatchString(l) {
			startIndex++
		} else {
			break
		}
	}

	out := make(playerStructs.Lyrics, 0, len(lines[startIndex:]))

	for _, lyric := range lines[startIndex:] {
		timings := timingRegexp.FindAllString(lyric, -1)

		for _, ts := range timings {
			lyric = strings.Replace(lyric, ts, "", 1)
		}
		lyric = strings.TrimSpace(lyric)

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

	return out
}

func parseTiming(timing string) float64 {
	// min 00:01.2, max 00:01.234
	if len(timing) < 9 || len(timing) > 11 {
		return -1
	}

	timeTag := timing[1 : len(timing)-1]
	parts := strings.Split(timeTag, ":")
	// now it should be
	// "00" and "01.23" (or .2, or .234, doesn't matter much)

	minutes, err := strconv.ParseFloat(parts[0], 64)
	if err != nil {
		return -1
	}
	seconds, err := strconv.ParseFloat(parts[1], 64)
	if err != nil {
		return -1
	}
	return minutes*60.0 + seconds
}

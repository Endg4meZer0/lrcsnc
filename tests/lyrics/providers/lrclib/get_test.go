package lrclib

import (
	"slices"
	"testing"

	lrclib "lrcsnc/internal/lyrics/providers/lrclib"
	"lrcsnc/internal/pkg/errors"
	playerStructs "lrcsnc/internal/pkg/structs/player"
	"lrcsnc/internal/pkg/types"
)

type Response struct {
	StatusCode int
	Body       string
}

// TestGetLyrics tests the ability to get different kinds of
// lyrics from LrcLib.
func TestGetLyrics(t *testing.T) {
	tests := []struct {
		name  string
		song  playerStructs.Song
		ldata playerStructs.LyricsData
	}{
		{
			name: "existing",
			song: playerStructs.Song{Title: "Earthless", Artists: []string{"Night Verses"}, Album: "From the Gallery of Sleep", Duration: 383},
			ldata: playerStructs.LyricsData{
				Lyrics: []playerStructs.Lyric{
					{Timing: 344.18, Text: "\"He is the one who gave me the horse"},
					{Timing: 346.74, Text: "So I could ride into the desert and see"},
					{Timing: 350.77, Text: "The future.\""},
					{Timing: 351.40999999999997, Text: ""},
					{Timing: 358.75, Text: "\"He is the one who gave me the horse"},
					{Timing: 361.39, Text: "So I could ride into the desert and see"},
					{Timing: 365.29, Text: "The future.\""},
					{Timing: 366.07, Text: ""},
				},
				LyricsState: types.LyricsStateSynced,
			},
		},
		{
			name: "existing-unicode",
			song: playerStructs.Song{Title: "狼之主", Artists: []string{"塞壬唱片-MSR"}, Album: "敘拉古人OST", Duration: 215},
			ldata: playerStructs.LyricsData{
				Lyrics: []playerStructs.Lyric{
					{Timing: 20.48, Text: "You're tough, but it's never been about you"},
					{Timing: 23.86, Text: "You're free, but cement your feet, a statue"},
					{Timing: 27.02, Text: "Your rules, but you'd rather make up something"},
					{Timing: 30.37, Text: "You're dead, you were never good for nothing"},
					{Timing: 33.62, Text: "Double negatives leading me in circles"},
					{Timing: 36.86, Text: "Twist infinity"},
					{Timing: 38.7, Text: "You drive me insane"},
					{Timing: 40.78, Text: "Hit it hard, a broken wall"},
					{Timing: 44.13, Text: "Hit hard, I gave it all"},
					{Timing: 47.36, Text: "Hit hard, a family tie, oh"},
					{Timing: 51.46, Text: "But you'd rather just fight"},
					{Timing: 53.41, Text: ""},
					{Timing: 56.62, Text: ""},
					{Timing: 66.11, Text: "Your dirt never washed off in an April shower"},
					{Timing: 69.49, Text: "You're crushed by the weight of those you can't devour"},
					{Timing: 72.73, Text: "You're armed, but the plan never executed"},
					{Timing: 75.82, Text: "You're shocked to let you hold a gun and never shoot it"},
					{Timing: 79.17, Text: "Regulating when the rules are simply saturated"},
					{Timing: 82.52, Text: "Is it everything, or is it just that I'm insane?"},
					{Timing: 86.43, Text: "Hit it hard, a broken wall"},
					{Timing: 89.84, Text: "Hit hard, I gave it all"},
					{Timing: 93.03, Text: "You tried and failed, all it fell, whoa"},
					{Timing: 100.18, Text: ""},
					{Timing: 167.69, Text: "All of this to say you lost it all to gain some power"},
					{Timing: 170.89, Text: "All of this, just say you plant a seed and kill the flower"},
					{Timing: 174.09, Text: "All of this to say you talk your way into your silence"},
					{Timing: 177.26, Text: "All of this is just a ploy to force your hand to violence"},
					{Timing: 180.36, Text: "It's a waste of time thinking you got tough"},
					{Timing: 183.53, Text: "When it's never really been enough"},
					{Timing: 185.24, Text: "Am I insane?"},
					{Timing: 187.34, Text: "Hit it hard, a broken wall"},
					{Timing: 190.52, Text: "Hit hard, I gave it all"},
					{Timing: 193.77, Text: "You tried and failed, all it fell, whoa"},
					{Timing: 200.68, Text: ""},
				},
				LyricsState: types.LyricsStateSynced,
			},
		},
		{
			name: "existing-plain",
			song: playerStructs.Song{Title: "All My Life, My Heart Has Yearned for a Thing I Cannot Name", Artists: []string{"Harm"}, Album: "a song you can't feel anymore.", Duration: 170},
			ldata: playerStructs.LyricsData{
				Lyrics: []playerStructs.Lyric{
					{Timing: 0, Text: "And here we stand now, intimate strangers in the end"},
					{Timing: 0, Text: "With these cold sheets we lay between"},
					{Timing: 0, Text: "We're holding onto what, makes us emptier again"},
					{Timing: 0, Text: ""},
					{Timing: 0, Text: "Your skin it lingers with my regret"},
					{Timing: 0, Text: "As I'm falling asleep and you're grabbing the wheel"},
					{Timing: 0, Text: "You say hold on tight it'll be okay"},
					{Timing: 0, Text: "Did you really think that?"},
					{Timing: 0, Text: ""},
					{Timing: 0, Text: "Pull me into the sea so we don't drown alone"},
					{Timing: 0, Text: "You've made your bed with mine so we can call this home"},
					{Timing: 0, Text: ""},
					{Timing: 0, Text: "And as the waves crash around us"},
					{Timing: 0, Text: "I watch your eyes grow with fear"},
					{Timing: 0, Text: "We have created a riptide"},
					{Timing: 0, Text: "Our fate lies here"},
					{Timing: 0, Text: ""},
					{Timing: 0, Text: "Your skin it lingers with my regret"},
					{Timing: 0, Text: "As I'm falling asleep and you're grabbing the wheel"},
					{Timing: 0, Text: "You say hold on tight it'll be okay"},
					{Timing: 0, Text: "Did you really think that?"},
					{Timing: 0, Text: ""},
					{Timing: 0, Text: "Pull me out of here"},
					{Timing: 0, Text: "Onto land"},
					{Timing: 0, Text: "As I disappear"},
					{Timing: 0, Text: "Hold my hand"},
					{Timing: 0, Text: "Lover, we're not playing fair"},
				},
				LyricsState: types.LyricsStatePlain,
			},
		},
		{
			name: "existing-instrumental",
			song: playerStructs.Song{Title: "Vice Wave", Artists: []string{"Night Verses"}, Album: "From the Gallery of Sleep", Duration: 300},
			ldata: playerStructs.LyricsData{
				Lyrics:      []playerStructs.Lyric{},
				LyricsState: types.LyricsStateInstrumental,
			},
		},
		{
			name: "non-existing",
			song: playerStructs.Song{Title: "Moonmore", Artists: []string{"Day Choruses"}, Album: "From the Gallery of Minecraft Pictures idk", Duration: 283},
			ldata: playerStructs.LyricsData{
				Lyrics:      []playerStructs.Lyric{},
				LyricsState: types.LyricsStateNotFound,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := lrclib.Provider{}.Get(tt.song)
			if err != nil && !(tt.ldata.LyricsState == types.LyricsStateNotFound && err == errors.ErrLyricsNotFound) {
				t.Errorf("[tests/lyrics/providers/lrclib/get/%v] Error: %v", tt.name, err)
				return
			}
			if !slices.Equal(got.Lyrics, tt.ldata.Lyrics) || got.LyricsState != tt.ldata.LyricsState {
				t.Errorf("[tests/lyrics/providers/lrclib/get/%v] Received %v, want %v", tt.name, got, tt.ldata)
			}
		})
	}
}

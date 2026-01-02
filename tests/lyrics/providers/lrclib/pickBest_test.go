package lrclib

import (
	"slices"
	"testing"

	lrclib "lrcsnc/internal/lyrics/providers/lrclib"
	playerStructs "lrcsnc/internal/pkg/structs/player"
	"lrcsnc/internal/pkg/types"
)

// TestPickBestLyrics tests the LrcLib's module to pick best lyrics data
// from received ones.
// This is a very primitive test because it tests only one variation.
// Maybe later will be more and this will turn into a proper test.
func TestPickBestLyrics(t *testing.T) {
	song := playerStructs.Song{Title: "Armageddon Eyes", Artists: []string{"Invent Animate", "Silent Planet"}, Album: "Bloom In Heaven", AlbumArtists: []string{"Invent Animate"}, Duration: 243}
	ldata := playerStructs.LyricsData{
		Lyrics: playerStructs.Lyrics{
			{Timing: 35.43, Text: "Stargazing, cosmic signs"},
			{Timing: 41.53, Text: "Born under moonlight"},
			{Timing: 44.35, Text: "Armageddon eyes"},
			{Timing: 47.18, Text: "I had a vision of your ghost wandering"},
			{Timing: 53.95, Text: "Return to make us whole"},
			{Timing: 59.08, Text: "The voice I felt in the distance was crying quietly"},
			{Timing: 69.24, Text: "Is it true you sing"},
			{Timing: 74.94, Text: "In the bloom you bring?"},
			{Timing: 81.09, Text: "I sense that you're following"},
			{Timing: 86.78999999999999, Text: "My ghost in the spaces, where I've never been"},
			{Timing: 94.15, Text: "Overtaking false dominions"},
			{Timing: 96.94, Text: "Poisoned meridians"},
			{Timing: 99.61, Text: "Alchemizing the oppression to equilibrium"},
			{Timing: 105.14, Text: "Immolation of the separation"},
			{Timing: 110.07, Text: "Plant the seed, Earth flourishing"},
			{Timing: 112.96000000000001, Text: "Revenants reap a garden from the ashes"},
			{Timing: 116.65, Text: "Is it true you sing"},
			{Timing: 122.21, Text: "In the bloom you bring?"},
			{Timing: 127.67, Text: "I sense that you're following"},
			{Timing: 132.7, Text: "Following my ghost in the spaces in between"},
			{Timing: 141.5, Text: "Armageddon eyes"},
			{Timing: 145.77, Text: "Paradise"},
			{Timing: 148.68, Text: "Return"},
			{Timing: 150.67000000000002, Text: "In those armageddon eyes"},
			{Timing: 164.43, Text: "Myriad"},
			{Timing: 174.32, Text: "Return me to the source"},
			{Timing: 178.57, Text: "The voice I felt in the distance"},
			{Timing: 183.86, Text: "The voice I felt"},
			{Timing: 187.94, Text: "Is it true you sing"},
			{Timing: 193.87, Text: "In the bloom you bring?"},
			{Timing: 199.81, Text: "I sense that you're following"},
			{Timing: 203.66, Text: "That you're following my ghost in the spaces in between"},
			{Timing: 213.19, Text: "Stargazing, cosmic signs"},
			{Timing: 218.87, Text: "Born under moonlight"},
			{Timing: 221.65, Text: "Armageddon eyes"},
			{Timing: 224.13, Text: "I had a vision of our ghosts becoming one"},
			{Timing: 231.05, Text: "The blooming makes us whole"},
			{Timing: 236, Text: "Reborn"},
			{Timing: 242.68, Text: ""},
		},
		LyricsState: types.LyricsStateSynced,
	}

	g, err := lrclib.Provider{}.Get(song)
	if err != nil || ldata.LyricsState != g.LyricsState || slices.CompareFunc(g.Lyrics, ldata.Lyrics, func(l1 playerStructs.Lyric, l2 playerStructs.Lyric) int {
		if l1.Timing == l2.Timing && l1.Text == l2.Text {
			return 0
		} else if l1.Timing < l2.Timing {
			return -1
		} else {
			return 1
		}
	}) != 0 {
		t.Errorf("pickBest doesn't work as intended, please debug... err = %v", err)
	}
}

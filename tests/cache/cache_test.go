package cache_test

import (
	"lrcsnc/internal/cache"
	"lrcsnc/internal/pkg/global"
	"lrcsnc/internal/pkg/structs"
	"lrcsnc/internal/pkg/types"
	"testing"
)

func TestStoreGetCycle(t *testing.T) {
	global.Config.C.Cache.StoreCondition.IfSynced = true
	global.Config.C.Cache.Dir = "$HOME/.cache/lrcsnc"
	testSong := structs.Song{
		Title:    "Is This A Test?",
		Artists:  []string{"Endg4me_"},
		Album:    "lrcsnc",
		Duration: 12.12,
		LyricsData: structs.LyricsData{
			Lyrics: []structs.Lyric{
				{Timing: 4.12, Text: "Pam-pam-pampararam"},
				{Timing: 7.54, Text: "Pam-pam-pam-param-pamparam"},
			},
			LyricsState: types.LyricsStateSynced,
		},
	}
	canStore := global.Config.C.Cache.StoreCondition.IsEnabledFor(testSong.LyricsData.LyricsState)
	if canStore == false {
		t.Error("[tests/cache/TestStoreGetCycle] ERROR: Failed to store lyrics in cache: store conditions are not working properly.")
	}
	err := cache.Store(&testSong)
	if err != nil {
		t.Errorf("[tests/cache/TestStoreGetCycle] ERROR: Failed to store lyrics in cache: %v", err)
	}
	defer cache.Remove(&testSong)

	global.Config.C.Cache.Enabled = false
	answerDisabled, cacheStateDisabled := cache.Fetch(&testSong)

	global.Config.C.Cache.Enabled = true
	answerInfLifeSpan, cacheStateInfLifeSpan := cache.Fetch(&testSong)

	if len(answerDisabled.Lyrics) != 0 || answerDisabled.LyricsState != 0 || cacheStateDisabled != cache.CacheStateDisabled {
		t.Errorf("[tests/cache/TestStoreGetCycle] ERROR: Disabling caching in config still allows fetching cached data")
	}

	if len(answerInfLifeSpan.Lyrics) != 2 || answerInfLifeSpan.LyricsState != 0 || cacheStateInfLifeSpan != cache.CacheStateActive {
		t.Errorf("[tests/cache/TestStoreGetCycle] ERROR: Received wrong cached data: expected %v, %v and %v, received %v, %v and %v",
			testSong.LyricsData.Lyrics, testSong.LyricsData.LyricsState, cache.CacheStateActive,
			answerInfLifeSpan.Lyrics, answerInfLifeSpan.LyricsState, cacheStateInfLifeSpan,
		)
	}

	if answerInfLifeSpan.Lyrics[0] != testSong.LyricsData.Lyrics[0] ||
		answerInfLifeSpan.Lyrics[1] != testSong.LyricsData.Lyrics[1] {
		t.Errorf("[tests/cache/TestStoreGetCycle] ERROR: Received wrong cached data: expected %v and %v, received %v and %v",
			testSong.LyricsData.Lyrics[0], testSong.LyricsData.Lyrics[1],
			answerInfLifeSpan.Lyrics[0], answerInfLifeSpan.Lyrics[1],
		)
	}
}

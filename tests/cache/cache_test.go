package cache_test

import (
	"lrcsnc/internal/cache"
	lyricStruct "lrcsnc/internal/lyrics/struct"
	"lrcsnc/internal/pkg/global"
	"lrcsnc/internal/pkg/types"
	playerStruct "lrcsnc/internal/player/struct"
	"os"
	"testing"
)

func TestStoreGetCycle(t *testing.T) {
	global.Config.C.Cache.Dir = "./lrcsnc-cache"
	if _, err := os.ReadDir(os.ExpandEnv(global.Config.C.Cache.Dir)); os.IsNotExist(err) {
		err = os.MkdirAll(os.ExpandEnv(global.Config.C.Cache.Dir), 0o755)
		if err != nil {
			t.Skip("[tests/cache/TestStoreGetCycle] WARN: Can't create database on this system, skipping...")
		}
	}

	cache.Init()
	defer cache.Close()

	global.Config.C.Cache.StoreCondition.IfSynced = true
	testSong := playerStruct.Song{
		Metadata: playerStruct.SongMetadata{
			Title:    "Is This A Test?",
			Artists:  []string{"Endg4me_"},
			Album:    "lrcsnc",
			Duration: 12.12,
		},
		LyricsData: lyricStruct.LyricsData{
			Lyrics: []lyricStruct.Lyric{
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
	err := cache.StorageInstance.Store(&testSong)
	if err != nil {
		t.Errorf("[tests/cache/TestStoreGetCycle] ERROR: Failed to store lyrics in cache: %v", err)
	}
	defer cache.StorageInstance.Remove(&testSong)

	// This test is now commented out since the Enabled check
	// is now only in lyrics/fetch.
	//
	// global.Config.C.Cache.Enabled = false
	// answerDisabled, cacheStateDisabled := cache.StorageInstance.Fetch(&testSong)
	// global.Config.C.Cache.Enabled = true

	answerInfLifeSpan, cacheStateInfLifeSpan := cache.StorageInstance.Fetch(&testSong)

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

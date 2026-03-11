package lyrics

import (
	"errors"
	"fmt"
	"strings"

	"lrcsnc/internal/cache"
	errs "lrcsnc/internal/lyrics/errors"
	"lrcsnc/internal/lyrics/providers"
	lyricStruct "lrcsnc/internal/lyrics/struct"
	"lrcsnc/internal/output/pkg/event"
	"lrcsnc/internal/output/server"
	"lrcsnc/internal/pkg/global"
	"lrcsnc/internal/pkg/log"
	"lrcsnc/internal/pkg/types"
)

// Fetch retrieves the lyrics data for the current song.
// It first checks if caching is enabled and attempts to retrieve the lyrics from the cache.
// If the lyrics are not found in the cache, it fetches the lyrics from the configured lyrics provider.
// If the lyrics are successfully retrieved and caching is enabled, it stores the lyrics in the cache.
func Fetch() (lyricStruct.LyricsData, error) {
	global.Player.M.Lock()
	song := global.Player.P.Song
	global.Player.M.Unlock()

	log.Debug("lyrics/fetch", fmt.Sprintf("Fetching lyrics for song %v - %v", strings.Join(song.Metadata.Artists, ", "), song.Metadata.Title))

	// yea i'm not covering this with mutexes good luck timing this out future me
	if global.Config.C.Cache.Enabled {
		cachedData, cacheState := cache.StorageInstance.Fetch(&song)
		if cacheState == cache.CacheStateActive {
			log.Debug("lyrics/fetch", "Cache hit; using cached data.")
			return cachedData, nil
		}
	}

	log.Debug("lyrics/fetch", fmt.Sprintf("Moving to the online part; using %v", global.Config.C.Lyrics.Provider))

	go server.ReceiveEvent(event.Event{
		Type: event.EventTypeLyricsStateChanged,
		Data: event.EventTypeLyricsStateChangedData{
			State: types.LyricsStateLoading,
		},
	})

	res, err := providers.Providers[global.Config.C.Lyrics.Provider].Get(song)
	if err != nil {
		if errors.Is(err, errs.NotFound) {
			log.Debug("lyrics/fetch", "The lyrics, unfortunately, were not found")
		} else {
			log.Error("lyrics/fetch", fmt.Sprintf("Could not get the lyrics: %s", err))
		}

		return res, err
	}

	log.Debug("lyrics/fetch", "Lyrics were successfully fetched from online")

	if global.Config.C.Cache.Enabled && global.Config.C.Cache.StoreCondition.IsEnabledFor(res.LyricsState) {
		song.LyricsData = res
		cache.StorageInstance.Store(&song)
	}

	return res, nil
}

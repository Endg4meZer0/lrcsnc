package lyrics

import (
	"errors"
	"fmt"
	"strings"

	"lrcsnc/internal/cache"
	"lrcsnc/internal/lyrics/providers"
	errs "lrcsnc/internal/pkg/errors"
	"lrcsnc/internal/pkg/global"
	"lrcsnc/internal/pkg/log"
	"lrcsnc/internal/pkg/structs"
)

// Fetch retrieves the lyrics data for the current song.
// It first checks if caching is enabled and attempts to retrieve the lyrics from the cache.
// If the lyrics are not found in the cache, it fetches the lyrics from the configured lyrics provider.
// If the lyrics are successfully retrieved and caching is enabled, it stores the lyrics in the cache.
func Fetch() (structs.LyricsData, error) {
	global.Player.M.Lock()
	song := global.Player.P.Song
	global.Player.M.Unlock()

	log.Debug("lyrics/fetch", fmt.Sprintf("Fetching lyrics for song %v - %v", strings.Join(song.Artists, ", "), song.Title))

	// yea i'm not covering this with mutexes good luck timing this out future me
	if global.Config.C.Cache.Enabled {
		cachedData, cacheState := cache.Fetch(&song)
		if cacheState == cache.CacheStateActive {
			return cachedData, nil
		}
	}

	log.Debug("lyrics/fetch", fmt.Sprintf("Moving to the online part; using %v", global.Config.C.Lyrics.Provider))

	res, err := providers.Providers[global.Config.C.Lyrics.Provider].Get(song)
	if err != nil {
		if errors.Is(err, errs.ErrLyricsNotFound) {
			log.Debug("lyrics/fetch", "The lyrics, unfortunately, were not found")
		} else {
			log.Error("lyrics/fetch", fmt.Sprintf("Could not get the lyrics: %s", err))
		}

		return res, err
	}

	log.Debug("lyrics/fetch", "Lyrics were successfully fetched from online")

	if global.Config.C.Cache.Enabled && global.Config.C.Cache.StoreCondition.IsEnabledFor(res.LyricsState) {
		song.LyricsData = res
		cache.Store(&song)
	}

	return res, nil
}

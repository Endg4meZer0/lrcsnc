package cache

import (
	"encoding/json"
	"fmt"
	"hash/fnv"
	"math"
	"os"
	"strings"
	"time"

	"lrcsnc/internal/pkg/errors"
	"lrcsnc/internal/pkg/global"
	"lrcsnc/internal/pkg/log"
	playerStructs "lrcsnc/internal/pkg/structs/player"
)

// Fetch retrieves the cached lyrics data for a given song.
// It first checks if the cache is enabled. If not, it returns an empty LyricsData and CacheStateDisabled.
// If the cache is enabled, it conplayerStructs the cache file path using the song ID.
// It attempts to read the cache file and unmarshal its contents into a LyricsData struct.
// If successful, it checks if the cache has expired based on the configured cache lifespan.
// It returns the cached data along with the appropriate CacheState (Active, Expired, or NonExistant).
func Fetch(song *playerStructs.Song) (playerStructs.LyricsData, CacheState) {
	global.Config.M.Lock()
	defer global.Config.M.Unlock()

	if !global.Config.C.Cache.Enabled {
		return playerStructs.LyricsData{}, CacheStateDisabled
	}
	cacheDirectory := getCacheDir()

	filename := getFilename(song)
	fullPath := cacheDirectory + "/" + filename + ".json"

	log.Debug("cache/Fetch", fmt.Sprintf("Fetching cache for song %v - %v under the name %v.json", strings.Join(song.Artists, ", "), song.Title, filename))

	if file, err := os.ReadFile(fullPath); err == nil {
		var cachedData playerStructs.LyricsData
		err = json.Unmarshal(file, &cachedData)
		if err != nil {
			log.Error("cache/Fetch", "Couldn't unmarshal the data: "+err.Error())
			return playerStructs.LyricsData{}, CacheStateNonExistant
		}

		log.Debug("cache/Fetch", "Done")

		if global.Config.C.Cache.LifeSpan != 0 {
			cacheStats, _ := os.Lstat(fullPath)
			isExpired := time.Since(cacheStats.ModTime()).Hours() >= float64(global.Config.C.Cache.LifeSpan)
			if isExpired {
				log.Debug("cache/Fetch", fmt.Sprintf("Cache has expired (%vh <= %vh)", time.Since(cacheStats.ModTime()).Hours(), float64(global.Config.C.Cache.LifeSpan)))
				return cachedData, CacheStateExpired
			} else {
				return cachedData, CacheStateActive
			}
		} else {
			return cachedData, CacheStateActive
		}
	} else {
		log.Debug("cache/Fetch", "Cache does not exist (in the end it doesn't even matter)")
		return playerStructs.LyricsData{}, CacheStateNonExistant
	}
}

// Store saves the lyrics data of a given song to a JSON file in the cache directory.
func Store(song *playerStructs.Song) error {
	global.Config.M.Lock()
	cacheDirectory := getCacheDir()
	global.Config.M.Unlock()

	if _, err := os.ReadDir(cacheDirectory); err != nil {
		err = os.Mkdir(cacheDirectory, 0o744)
		if err != nil {
			return errors.DirUnwriteable
		}
	}

	filename := getFilename(song)
	fullPath := cacheDirectory + "/" + filename + ".json"

	log.Debug("cache/Store", fmt.Sprintf("Storing cache for song %v - %v under the name %v.json", strings.Join(song.Artists, ", "), song.Title, filename))

	encodedData, err := json.Marshal(song.LyricsData)
	if err != nil {
		log.Error("cache/Store", "Failed to marshal the data: "+err.Error())
		return errors.MarshalFail
	}

	if err := os.WriteFile(fullPath, []byte(encodedData), 0o644); err != nil {
		log.Error("cache/Store", "Failed to write the cache file: "+err.Error())
		return errors.FileUnwriteable
	}
	log.Debug("cache/Store", "Done")
	return nil
}

// Remove deletes the cached data for the given song from the cache directory.
// If the cache directory is not reachable, it returns an ErrDirUnreachable error.
// If the specific cached file for the song cannot be removed,
// it logs an error message and returns an ErrFileUnreachable error.
//
// It's not really used anywhere for now, only in cache unit test,
// but it was useful in the past and might be again in the future.
func Remove(song *playerStructs.Song) error {
	global.Config.M.Lock()
	cacheDirectory := getCacheDir()
	global.Config.M.Unlock()

	if _, err := os.ReadDir(cacheDirectory); err != nil {
		return errors.DirUnreachable
	}

	filename := getFilename(song)
	fullPath := cacheDirectory + "/" + filename + ".json"
	log.Debug("cache/Remove", fmt.Sprintf("Removing cache for song %v - %v under the name %v.json", strings.Join(song.Artists, ", "), song.Title, filename))

	if err := os.Remove(fullPath); err != nil {
		log.Error("cache/Remove", fmt.Sprintf("Couldn't delete the cached data for %v. Maybe the data didn't exist in the first place?", filename))
		return errors.FileUnreachable
	}
	log.Debug("cache/Remove", "Done")
	return nil
}

func getCacheDir() string {
	return os.ExpandEnv(global.Config.C.Cache.Dir)
}

func getFilename(song *playerStructs.Song) string {
	h := fnv.New64a()
	h.Write(fmt.Appendf([]byte{}, "%v.%v.%v.%v",
		song.Title,
		strings.Join(song.Artists, ", "),
		song.Album,
		math.Round(song.Duration),
	))
	return fmt.Sprint(h.Sum64())
}

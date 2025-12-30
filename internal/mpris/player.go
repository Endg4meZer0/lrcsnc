package mpris

import (
	"slices"
	"strings"

	"lrcsnc/internal/pkg/global"
	"lrcsnc/internal/pkg/log"
	playerStructs "lrcsnc/internal/pkg/structs/player"
	"lrcsnc/internal/pkg/types"

	"github.com/Endg4meZer0/go-mpris"
	"github.com/godbus/dbus/v5"
)

// ChangePlayer is used to:
//
// 1) Check if the current player is still alive (in that case does nothing)
//
// 2) Find a new player among MPRIS clients that fits the filters
// (if it finds, changes the player and informs corresponding channels)
func ChangePlayer() error {
	log.Debug("mpris/ChangePlayer", "Started")

	// Getting the players list
	players, err := mpris.List(conn)
	if err != nil {
		log.Error("mpris/ChangePlayer", "Got an error while using mpris.List: "+err.Error())
		return err
	}
	log.Debug("mpris/ChangePlayer", "Current players available in MPRIS: "+strings.Join(players, ", "))

	// A little helper function to filter players
	playerInFilter := func(player string) bool {
		global.Config.M.Lock()
		defer global.Config.M.Unlock()
		if len(global.Config.C.Player.IncludedPlayers) != 0 {
			for _, includedPlayer := range global.Config.C.Player.IncludedPlayers {
				if strings.Contains(player, includedPlayer) {
					return true
				}
			}
		}

		if len(global.Config.C.Player.ExcludedPlayers) != 0 {
			for _, excludedPlayer := range global.Config.C.Player.ExcludedPlayers {
				if strings.Contains(player, excludedPlayer) {
					return false
				}
			}
		}

		return true
	}

	// Check if the current player is still alive and kicking
	if player != nil {
		log.Debug("mpris/ChangePlayer", "There is a player handle stored already. Checking if it's alive yet...")
		currentPlayer := player.GetName()
		if slices.Contains(players, currentPlayer) && playerInFilter(currentPlayer) {
			log.Debug("mpris/ChangePlayer", "It is alive. No extra action taken.")
			return nil
		} else {
			log.Debug("mpris/ChangePlayer", "It is dead. Starting unregistering procedure.")
			// Remove the signal handler from the current player before assigning to a new player
			err = player.UnregisterSignalReceiver(playerSignalReceiver)
			if err != nil {
				log.Error(
					"mpris/ChangePlayer",
					"An error occurred while unregistering signal receive channel. Recreating and redeploying the watcher. More: "+err.Error(),
				)
				playerSignalReceiver = make(chan *dbus.Signal)
				go playerSignalWatcher()
			} else {
				log.Debug(
					"mpris/ChangePlayer",
					"Successfully unregistered signal receiver.",
				)
			}
		}
	}

	// Find a new player that passes the filters and supports MPRIS to the extent that we need
	log.Debug("mpris/ChangePlayer", "Starting to pick a new player to watch.")
	for _, p := range players {
		if playerInFilter(p) {
			log.Info("mpris/ChangePlayer", "Found a fitting player: '"+p+"', trying to connect...")
			player = mpris.New(conn, p)

			pass := true

			pname, err := player.GetIdentity()
			if err != nil {
				log.Error("mpris/ChangePlayer", "Got an error when using player.GetIdentity: "+err.Error())
				pass = false
			}
			pps, err := GetPlaybackStatus()
			if err != nil {
				log.Error("mpris/ChangePlayer", "Got an error when using player.GetPlaybackStatus: "+err.Error())
				pass = false
			}
			ppos, err := GetPosition()
			if err != nil {
				log.Error("mpris/ChangePlayer", "Got an error when using player.GetPosition: "+err.Error())
				pass = false
			}
			prate, err := GetRate()
			if err != nil {
				log.Error("mpris/ChangePlayer", "Got an non-critical error when using player.GetRate: "+err.Error())
				prate = 1
			}
			md, err := GetMetadata()
			if err != nil {
				log.Error("mpris/ChangePlayer", "Got an error when using player.GetMetadata: "+err.Error())
				pass = false
			}

			err = CheckMetadata(md)
			if err != nil {
				log.Error("mpris/ChangePlayer", "Got an error when using ApplyMetadataOntoGlobal: "+err.Error())
				pass = false
			}

			if !pass {
				log.Info("mpris/ChangePlayer", "Failed to gather necessary data from player '"+p+"'. Skipping...")
				continue
			}

			// Lock the player mutex while we're updating the data
			global.Player.M.Lock()

			global.Player.P.Name = pname
			global.Player.P.PlaybackStatus = pps
			global.Player.P.Position = ppos
			global.Player.P.Rate = prate
			ApplyMetadataOntoGlobal(md)

			global.Player.M.Unlock()

			// Register signal receiver for the new player
			err = player.RegisterSignalReceiver(playerSignalReceiver)
			if err != nil {
				log.Fatal("mpris/ChangePlayer/RegisterSignalReceiver", "Cannot watch for player signals. More: "+err.Error())
			}

			log.Debug("mpris/ChangePlayer", "Successfully connected to player '"+p+"' and gathered necessary data.")
			log.Info("mpris/ChangePlayer", "Switched to player '"+p+"'")

			return nil
		}
	}

	// If no player is found, set the player to nil
	log.Info("mpris/ChangePlayer", "No active player found. Zzz")
	player = nil
	global.Player.M.Lock()
	global.Player.P = playerStructs.Player{
		PlaybackStatus: mpris.PlaybackStopped,
		Position:       0.0,
		Rate:           1.0,
		Song: playerStructs.Song{
			LyricsData: playerStructs.LyricsData{
				LyricsState: types.LyricsStateUnknown,
			},
		},
		Name: "",
	}
	global.Player.M.Unlock()
	return nil
}

func CheckMetadata(md mpris.Metadata) (err error) {
	_, err = md.Title()
	if err != nil {
		return err
	}
	_, err = md.Artist()
	if err != nil {
		return err
	}
	_, err = md.Album()
	if err != nil {
		return err
	}
	_, err = md.Length()
	if err != nil {
		return err
	}
	return
}

// ApplyMetadataOnGlobal edits the global player variable (specifically the Song)
// in accordance to provided metadata.
//
// It does NOT lock the mutex.
func ApplyMetadataOntoGlobal(md mpris.Metadata) (err error) {
	global.Player.P.Song.Title, err = md.Title()
	if err != nil {
		return err
	}
	global.Player.P.Song.Artists, err = md.Artist()
	if err != nil {
		return err
	}
	global.Player.P.Song.Album, err = md.Album()
	if err != nil {
		return err
	}
	dur, err := md.Length()
	if err != nil {
		return err
	}
	global.Player.P.Song.Duration = float64(dur) / 1000 / 1000
	global.Player.P.Song.LyricsData.LyricsState = types.LyricsStateLoading

	return nil
}

// GetPlaybackStatus returns current playback status
func GetPlaybackStatus() (mpris.PlaybackStatus, error) {
	if player == nil {
		return mpris.PlaybackStopped, nil
	}

	return player.GetPlaybackStatus()
}

// GetPosition returns current position
func GetPosition() (float64, error) {
	if player == nil {
		return 0, nil
	}

	val, err := player.GetPosition()
	if err != nil {
		return 0, err
	}

	return float64(val) / 1000 / 1000, nil
}

// GetRate returns current rate
func GetRate() (float64, error) {
	if player == nil {
		return 0, nil
	}

	return player.GetRate()
}

// GetMetadata returns current metadata
func GetMetadata() (mpris.Metadata, error) {
	if player == nil {
		return nil, nil
	}

	return player.GetMetadata()
}

// SetPosition sets new position for the player
func SetPosition(pos float64) error {
	if player == nil {
		return nil
	}

	return player.SetPosition(int64(pos * 1000 * 1000))
}

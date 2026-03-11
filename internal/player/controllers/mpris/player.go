//go:build linux

package mpris

import (
	"context"
	"slices"
	"strings"

	"lrcsnc/internal/pkg/global"
	"lrcsnc/internal/pkg/log"
	"lrcsnc/internal/pkg/types"
	playerStruct "lrcsnc/internal/player/struct"

	"github.com/Endg4meZer0/go-mpris"
	"github.com/godbus/dbus/v5"
)

// ChangePlayer is used to:
//
// 1) Check if the current c.player is still alive (in that case does nothing)
//
// 2) Find a new c.player among MPRIS clients that fits the filters
// (if it finds, changes the c.player and informs corresponding channels)
func (c *controller) ChangePlayer() error {
	log.Debug("mpris/ChangePlayer", "Started")

	// Getting the c.players list
	players, err := mpris.List(c.conn)
	if err != nil {
		log.Error("mpris/ChangePlayer", "Got an error while using mpris.List: "+err.Error())
		return err
	}
	log.Debug("mpris/ChangePlayer", "Current c.players available in MPRIS: "+strings.Join(players, ", "))

	// Check if the current c.player is still alive and kicking
	if c.player != nil {
		log.Debug("mpris/ChangePlayer", "There is a c.player handle stored already. Checking if it's alive yet...")
		currentPlayer := c.player.GetName()
		if slices.Contains(players, currentPlayer) && playerInFilter(currentPlayer) {
			log.Debug("mpris/ChangePlayer", "It is alive. No extra action taken.")
			return nil
		} else {
			log.Debug("mpris/ChangePlayer", "It is dead. Starting unregistering procedure.")
			// Remove the signal handler from the current c.player before assigning to a new c.player
			err = c.player.UnregisterSignalReceiver(c.playerSignalReceiver)
			if err != nil {
				log.Error(
					"mpris/ChangePlayer",
					"An error occurred while unregistering signal receive channel. Recreating and redeploying the watcher. More: "+err.Error(),
				)
				c.playerSignalWatcherCancelFunc()
				c.playerSignalReceiver = make(chan *dbus.Signal)
				ctx, cfunc := context.WithCancel(context.Background())
				c.playerSignalWatcher(ctx)
				c.playerSignalWatcherCancelFunc = cfunc
			} else {
				log.Debug(
					"mpris/ChangePlayer",
					"Successfully unregistered signal receiver.",
				)
			}
		}
	}

	// Find a new c.player that passes the filters and supports MPRIS to the extent that we need
	log.Debug("mpris/ChangePlayer", "Starting to pick a new c.player to watch.")
	for _, p := range players {
		if playerInFilter(p) {
			log.Info("mpris/ChangePlayer", "Found a fitting c.player: '"+p+"', trying to connect...")
			c.player = mpris.New(c.conn, p)

			pass := true

			pname, err := c.player.GetIdentity()
			if err != nil {
				log.Error("mpris/ChangePlayer", "Got an error when using c.player.GetIdentity: "+err.Error())
				pass = false
			}
			pps, err := c.GetPlaybackStatus()
			if err != nil {
				log.Error("mpris/ChangePlayer", "Got an error when using c.player.GetPlaybackStatus: "+err.Error())
				pass = false
			}
			ppos, err := c.GetPosition()
			if err != nil {
				log.Error("mpris/ChangePlayer", "Got an error when using c.player.GetPosition: "+err.Error())
				pass = false
			}
			prate, err := c.GetRate()
			if err != nil {
				log.Error("mpris/ChangePlayer", "Got an non-critical error when using c.player.GetRate: "+err.Error())
				prate = 1
			}
			md, err := c.GetMPRISMetadata()
			if err != nil {
				log.Error("mpris/ChangePlayer", "Got an error when using c.player.GetMetadata: "+err.Error())
				pass = false
			}

			checkMDPass, errs := checkMetadata(md)
			if len(errs) != 0 {
				log.Error("mpris/ChangePlayer", "Got errors when checking metadata.")
				for _, err := range errs {
					log.Error("mpris/ChangePlayer", "\t"+err.Error())
				}
				if checkMDPass {
					log.Warn("mpris/ChangePlayer", "These errors are not critical, but related functionality may be impacted.")
				}
				pass = checkMDPass && pass
			}

			if !pass {
				log.Info("mpris/ChangePlayer", "Failed to gather necessary data from c.player '"+p+"'. Skipping...")
				continue
			}

			// Lock the c.player mutex while we're updating the data
			global.Player.M.Lock()

			global.Player.P.Update(pname, mprisPsParse(pps), ppos, prate, metadataParse(md))

			global.Player.M.Unlock()

			// Register signal receiver for the new c.player
			err = c.player.RegisterSignalReceiver(c.playerSignalReceiver)
			if err != nil {
				log.Fatal("mpris/ChangePlayer/RegisterSignalReceiver", "Cannot watch for c.player signals. More: "+err.Error())
			}

			log.Debug("mpris/ChangePlayer", "Successfully connected to c.player '"+p+"' and gathered necessary data.")
			log.Info("mpris/ChangePlayer", "Switched to c.player '"+p+"'")

			return nil
		}
	}

	// If no c.player is found, set the c.player to nil
	log.Info("mpris/ChangePlayer", "No active c.player found. Zzz")
	c.player = nil
	global.Player.M.Lock()

	global.Player.P.Update("", types.PlaybackStatusStopped, 0.0, 1.0, playerStruct.SongMetadata{})

	global.Player.M.Unlock()
	return nil
}

func checkMetadata(md mpris.Metadata) (pass bool, errs []error) {
	pass = true
	var err error
	_, err = md.Title()
	if err != nil {
		pass = false
		errs = append(errs, err)
	}
	_, err = md.Artist()
	if err != nil {
		pass = false
		errs = append(errs, err)
	}
	_, err = md.Length()
	if err != nil {
		pass = false
		errs = append(errs, err)
	}
	_, err = md.Album()
	if err != nil {
		errs = append(errs, err)
	}
	_, err = md.AlbumArtist()
	if err != nil {
		errs = append(errs, err)
	}
	_, err = md.ArtURL()
	if err != nil {
		errs = append(errs, err)
	}
	_, err = md.URL()
	if err != nil {
		errs = append(errs, err)
	}
	return
}

func mprisPsParse(ps mpris.PlaybackStatus) types.PlaybackStatus {
	switch ps {
	case mpris.PlaybackPlaying:
		return types.PlaybackStatusPlaying
	case mpris.PlaybackPaused:
		return types.PlaybackStatusPaused
	default:
		return types.PlaybackStatusStopped
	}
}

func metadataParse(md mpris.Metadata) (smd playerStruct.SongMetadata) {
	smd.Title, _ = md.Title()
	smd.Artists, _ = md.Artist()
	smd.Album, _ = md.Album()
	smd.AlbumArtists, _ = md.AlbumArtist()
	smd.AlbumArt, _ = md.ArtURL()
	smd.Url, _ = md.URL()
	dur, _ := md.Length()
	smd.Duration = float64(dur) / 1000 / 1000

	return
}

// GetPlaybackStatus returns current playback status
func (c *controller) GetPlaybackStatus() (mpris.PlaybackStatus, error) {
	if c.player == nil {
		return mpris.PlaybackStopped, nil
	}

	return c.player.GetPlaybackStatus()
}

// GetPosition returns current position
func (c *controller) GetPosition() (float64, error) {
	if c.player == nil {
		return 0, nil
	}

	val, err := c.player.GetPosition()
	if err != nil {
		return 0, err
	}

	return float64(val) / 1000 / 1000, nil
}

// GetRate returns current rate
func (c *controller) GetRate() (float64, error) {
	if c.player == nil {
		return 0, nil
	}

	return c.player.GetRate()
}

// GetMetadata returns current metadata
func (c *controller) GetMPRISMetadata() (mpris.Metadata, error) {
	if c.player == nil {
		return nil, nil
	}

	return c.player.GetMetadata()
}

// SetPosition sets new position for the c.player
func (c *controller) SetPosition(pos float64) error {
	if c.player == nil {
		return nil
	}

	return c.player.SetPosition(int64(pos * 1000 * 1000))
}

// playerInFilter is a helper function for detecting if the c.player fits the config's filters
func playerInFilter(player string) bool {
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

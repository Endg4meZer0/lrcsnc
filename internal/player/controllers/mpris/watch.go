//go:build linux

package mpris

import (
	"context"
	"fmt"
	"lrcsnc/internal/pkg/log"
	"lrcsnc/internal/player/signals"
	"reflect"
	"strings"

	mprislib "github.com/Endg4meZer0/go-mpris"
	"github.com/godbus/dbus/v5"
)

func (c *controller) playerSignalWatcher(ctx context.Context) {
	for {
		select {
		case signal, ok := <-c.playerSignalReceiver:
			if !ok {
				log.Error("mpris/watchPlayerSignals", "Player signal channel sent invalid *dbus.Signal object")
				continue
			}
			switch mprislib.GetSignalType(signal) {
			case mprislib.SignalSeeked:
				if len(signal.Body) != 1 {
					log.Error("mpris/watchPlayerSignals", fmt.Sprintf("The Seeked signal should have only 1 value, but there are %v: %v.", len(signal.Body), signal.Body))
				}

				var v int64
				v, ok := signal.Body[0].(int64)
				if !ok {
					v2, ok := signal.Body[0].(uint64)
					if !ok {
						log.Error("mpris/watchPlayerSignals", "The Seeked signal contains a value that is not int64 or uint64.")
						continue
					}
					v = int64(v2)
				}

				signals.MessageChannel <- signals.Message{Type: signals.SignalSeeked, Data: float64(v) / 1000 / 1000}
			case mprislib.SignalPropertiesChanged:
				if len(signal.Body) != 3 {
					log.Error("mpris/watchPlayerSignals", fmt.Sprintf("The PropertiesChanged signal should have 3 values, but there are %v: %v.", len(signal.Body), signal.Body))
					continue
				}

				v, ok := signal.Body[1].(map[string]dbus.Variant)
				if !ok {
					log.Error("mpris/watchPlayerSignals", fmt.Sprintf("The PropertiesChanged signal contains not a map[string]dbus.Variant but a %v.", reflect.TypeOf(signal.Body[1]).Name()))
					continue
				}

				// Check whether the playback status has changed
				playbackStatus, ok := v["PlaybackStatus"]
				if ok {
					p, ok := playbackStatus.Value().(string)
					if !ok {
						log.Error("mpris/watchPlayerSignals", fmt.Sprintf("The PlaybackStatus value is not a string but a %v.", playbackStatus.Signature().String()))
					} else {
						signals.MessageChannel <- signals.Message{Type: signals.SignalPlaybackStatusChanged, Data: mprisPsParse(mprislib.PlaybackStatus(p))}
					}
				}

				// Check whether the rate has changed
				rate, ok := v["Rate"]
				if ok {
					r, ok := rate.Value().(float64)
					if !ok {
						log.Error("mpris/watchPlayerSignal", fmt.Sprintf("The Rate value is not a float64 but a %v.", reflect.TypeOf(rate.Value()).Name()))
					} else {
						signals.MessageChannel <- signals.Message{Type: signals.SignalRateChanged, Data: r}
					}
				}

				// Check whether the metadata has changed
				metadata, ok := v["Metadata"]
				if ok {
					m, ok := metadata.Value().(map[string]dbus.Variant)
					if !ok {
						log.Error("mpris/watchPlayerSignal", fmt.Sprintf("The Metadata value is not a map[string]dbus.Variant but a %v.", reflect.TypeOf(metadata.Value()).Name()))
					} else {
						signals.MessageChannel <- signals.Message{Type: signals.SignalMetadataChanged, Data: metadataParse(mprislib.Metadata(m))}
					}
				}
			}
		case <-ctx.Done():
			return
		}
	}
}

func (c *controller) nameOwnerChangeWatcher(ctx context.Context) {
	for {
		select {
		case s, ok := <-c.nameOwnerChangedSignalReceiver:
			// The channel should not be empty, but we must check to be sure.
			// If it is empty, then most likely it happened because of premature connection close.
			// That means we should reopen the connection.
			if !ok {
				log.Warn("mpris/nameOwnerChangesWatcher", "The channel is closed. Reconnecting to DBus...")
				go c.Connect()
				return
			}

			// Skipping any not NameOwnerChanged signals
			if signalType := mprislib.GetSignalType(s); signalType != mprislib.SignalNameOwnerChanged {
				// log.Debug("mpris/nameOwnerChangesWatcher", fmt.Sprintf("Signal is not NameOwnerChanged (%s). Skipping...", signalType))
				continue
			}

			// Skipping any signals not related to MPRIS
			// See also: https://dbus.freedesktop.org/doc/dbus-specification.html#bus-messages-name-owner-changed
			if !strings.HasPrefix(s.Body[0].(string), mprislib.BaseInterface) {
				// log.Debug("mpris/nameOwnerChangesWatcher", fmt.Sprintf("Signal is not related to MPRIS (%s). Skipping...", s.Body[0].(string)))
				continue
			}

			log.Debug("mpris/nameOwnerChangesWatcher", "Received a useful signal. Probing ChangePlayer and passing to signals.MessageChannel...")

			err := c.ChangePlayer()
			if err != nil {
				log.Error("mpris/nameOwnerChangesWatcher", err.Error())
			}

			signals.MessageChannel <- signals.Message{Type: signals.SignalPlayerChanged, Data: nil}
		case <-ctx.Done():
			return
		}
	}
}

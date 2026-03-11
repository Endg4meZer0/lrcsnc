//go:build linux

package mpris

import (
	"context"
	"lrcsnc/internal/pkg/log"
	"lrcsnc/internal/player/signals"

	"github.com/Endg4meZer0/go-mpris"
	"github.com/godbus/dbus/v5"
)

// Connect connects to D-Bus and returns any error it encounters.
// The connection is then stored privately in this module.
func (c *controller) Connect() error {
	var err error

	// Get a private connection from D-Bus
	c.conn, err = dbus.SessionBusPrivate()
	if err != nil {
		log.Error("mpris/Connect", err.Error())
		return err
	}
	log.Debug("mpris/Connect", "Got a private connection from D-Bus")

	// Do the needed procedures like auth...
	err = c.conn.Auth(nil)
	if err != nil {
		c.Disconnect()
		log.Error("mpris/Connect", err.Error())
		return err
	}
	log.Debug("mpris/Connect", "Authentificated the connection")

	// ...and hello
	err = c.conn.Hello()
	if err != nil {
		c.Disconnect()
		return err
	}
	log.Debug("mpris/Connect", "Greeted D-Bus")

	// Also register the NameOwnerChanged signal receiver channel to see
	// if there are new players or old are removed
	err = mpris.RegisterNameOwnerChanged(c.conn, c.nameOwnerChangedSignalReceiver)
	if err != nil {
		c.Disconnect()
		log.Fatal("mpris/Connect", "Cannot watch for player signals. More: "+err.Error())
		return err
	}
	log.Debug("mpris/Connect", "Registered the name owner changed signal receiver channel")

	// And now we can get the active player (if there is any)
	err = c.ChangePlayer()
	if err != nil {
		log.Error("mpris/Connect", err.Error())
	}
	log.Debug("mpris/Connect", "Got the initial player info from MPRIS")

	signals.MessageChannel <- signals.Message{Type: signals.SignalReady, Data: nil}
	log.Debug("mpris/Connect", "Sent a SignalReady to MPRISMessageChannel")

	// And deploy the watchers after cancelling past ones if there were any
	if c.playerSignalWatcherCancelFunc != nil {
		c.playerSignalWatcherCancelFunc()
	}
	if c.nameOwnerChangedSignalWatcherCancelFunc != nil {
		c.nameOwnerChangedSignalWatcherCancelFunc()
	}

	ctx1, cfunc1 := context.WithCancel(context.Background())
	go c.playerSignalWatcher(ctx1)
	c.playerSignalWatcherCancelFunc = cfunc1

	ctx2, cfunc2 := context.WithCancel(context.Background())
	go c.nameOwnerChangeWatcher(ctx2)
	c.nameOwnerChangedSignalWatcherCancelFunc = cfunc2

	log.Debug("mpris/Connect", "Deployed the MPRIS/D-Bus signal watchers")

	log.Info("mpris/Connect", "Successfully connected to D-Bus")

	return nil
}

// Disconnect disconnects from D-Bus.
func (c *controller) Disconnect() error {
	c.playerSignalWatcherCancelFunc()
	c.nameOwnerChangedSignalWatcherCancelFunc()

	err := mpris.UnregisterNameOwnerChanged(c.conn, c.nameOwnerChangedSignalReceiver)
	if err != nil {
		log.Error("mpris/Disconnect", err.Error())
		return err
	}

	err = c.conn.Close()
	if err != nil {
		log.Error("mpris/Disconnect", err.Error())
		return err
	}

	log.Info("mpris/Disconnect", "Successfully disconnected from DBus")
	return nil
}

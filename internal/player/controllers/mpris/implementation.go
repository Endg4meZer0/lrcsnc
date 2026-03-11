//go:build linux

package mpris

import (
	mprislib "github.com/Endg4meZer0/go-mpris"
	"github.com/godbus/dbus/v5"
)

type controller struct {
	conn                           *dbus.Conn
	player                         *mprislib.Player
	playerSignalReceiver           chan *dbus.Signal
	nameOwnerChangedSignalReceiver chan *dbus.Signal

	playerSignalWatcherCancelFunc           func()
	nameOwnerChangedSignalWatcherCancelFunc func()
}

var Controller *controller = &controller{
	playerSignalReceiver:           make(chan *dbus.Signal),
	nameOwnerChangedSignalReceiver: make(chan *dbus.Signal),
}

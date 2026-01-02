package server

import (
	"errors"
	"io/fs"
	"lrcsnc/internal/output/pkg/event"
	"lrcsnc/internal/pkg/global"
	"lrcsnc/internal/pkg/log"
	"net"
	"os"
)

type server struct {
	protocol    string
	listenPath  string
	netListener net.Listener
	conns       Connections

	// 0 - PlayerChanged
	// 1 - RateChanged
	// 2 - SongChanged
	// 3 - PlaybackStatusChanged
	// 4 - LyricsStateChanged
	// 5 - LyricsChanged
	// 6 - ActiveLyricChanged
	// 7 - OverwriteRequired
	lastEvents []event.Event
}

var srv server

func InitServer() {
	global.Config.M.Lock()

	// If we are not specifically server-mode...
	if !global.Config.C.Net.IsServer {
		// Check if we are in client mode.
		if global.Config.C.Net.Protocol != "" {
			global.Config.M.Unlock()
			return
		}

		// If not, we must be in standalone mode, in which case
		// the communication process will go internally.
		srv = server{
			protocol:    "",
			listenPath:  "",
			netListener: nil,
			conns:       Connections{},
			lastEvents:  nil,
		}

		global.Config.M.Unlock()
		return
	}

	protocol := global.Config.C.Net.Protocol
	listenPath := global.Config.C.Net.ListenAt
	global.Config.M.Unlock()

	// If protocol is "unix", check
	// if listenPath is a valid path that is reachable
	// and if there is a file already there,
	// give fatal error.
	if protocol == "unix" {
		if d, err := os.Lstat(listenPath); err != nil && !errors.Is(err, os.ErrNotExist) {
			log.Fatal("output/server", "Listen path is invalid: file exists or there is an issue with permissions. You may want to check it out and remove file/change perms manually. More: "+err.Error())
		} else if err == nil && d.Mode()&fs.ModeSocket == fs.ModeSocket {
			log.Fatal("output/server", "The specified listen path is busy: a socket is already open there. Specify another listen path or, if nothing is actually using this socket, you could try removing it.")
		}
	}

	l, err := net.Listen(protocol, listenPath)
	if err != nil {
		log.Fatal("output/server", "An error occured when trying to start a server: "+err.Error())
	}

	srv = server{
		protocol:    protocol,
		listenPath:  listenPath,
		netListener: l,
		conns:       Connections{Map: make(map[uint]net.Conn), FirstFreeSlot: 0},
		lastEvents:  makeDefaultLastEvents(),
	}

	go func() {
		for {
			conn, err := srv.netListener.Accept()
			log.Debug("output/server", "A client connected")
			srv.conns.M.Lock()
			if err != nil {
				log.Error("output/server", "Failed to accept a connection: "+err.Error())
			}
			f := srv.conns.FirstFreeSlot
			srv.conns.Map[f] = conn
			srv.conns.M.Unlock()
			go srv.sendLastEventsTo(f)
			go srv.conns.FindFirstAvailableSlot()
		}
	}()
}

func ReceiveEvent(e event.Event) {
	srv.handleEvent(e)
}

func (s *server) handleEvent(e event.Event) {
	if e.Type == event.EventTypeServerClosed {
		log.Info("output/server", "Start closing the server...")
		for i, c := range s.conns.Map {
			c.Close()
			delete(s.conns.Map, i)
		}
		if s.protocol == "unix" {
			os.Remove(s.listenPath)
		}
		return
	}

	if s.protocol != "" {
		switch e.Type {
		case event.EventTypePlayerChanged:
			s.lastEvents[0] = e
		case event.EventTypeRateChanged:
			s.lastEvents[1] = e
		case event.EventTypeSongChanged:
			s.lastEvents[2] = e
		case event.EventTypePlaybackStatusChanged:
			s.lastEvents[3] = e
		case event.EventTypeLyricsStateChanged:
			s.lastEvents[4] = e
		case event.EventTypeLyricsChanged:
			s.lastEvents[5] = e
		case event.EventTypeActiveLyricChanged:
			s.lastEvents[6] = e
		case event.EventTypeOverwriteRequired:
			s.lastEvents[7] = e
		}
	}

	s.sendEvent(e)
}

func CloseServer() {
	srv.handleEvent(event.Event{Type: event.EventTypeServerClosed, Data: event.EventTypeServerClosedData{}})
}

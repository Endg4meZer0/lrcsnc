package server

import (
	"errors"
	"lrcsnc/internal/output/pkg/event"
	"lrcsnc/internal/pkg/global"
	"lrcsnc/internal/pkg/log"
	"net"
	"os"
)

type Server struct {
	Protocol    string
	ListenPath  string
	NetListener net.Listener
	Conns       Connections

	// 0 - PlayerChanged
	// 1 - RateChanged
	// 2 - SongChanged
	// 3 - PlaybackStatusChanged
	// 4 - LyricsStateChanged
	// 5 - LyricsChanged
	// 6 - ActiveLyricChanged
	// 7 - OverwriteRequired
	LastEvents []event.Event
}

var server Server

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
		server = Server{
			Protocol:    "",
			ListenPath:  "",
			NetListener: nil,
			Conns:       Connections{},
			LastEvents:  nil,
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
		if _, err := os.Lstat(listenPath); err != nil && !errors.Is(err, os.ErrNotExist) {
			log.Fatal("output/server", "Listen path is invalid: file exists or there is an issue with permissions. You may want to check it out and remove file/change perms manually. More: "+err.Error())
		}
	}

	l, err := net.Listen(protocol, listenPath)
	if err != nil {
		log.Fatal("output/server", "An error occured when trying to start a server: "+err.Error())
	}

	server = Server{
		Protocol:    protocol,
		ListenPath:  listenPath,
		NetListener: l,
		Conns:       Connections{Map: make(map[uint]net.Conn), FirstFreeSlot: 0},
		LastEvents:  makeDefaultLastEvents(),
	}

	go func() {
		for {
			conn, err := server.NetListener.Accept()
			server.Conns.M.Lock()
			if err != nil {
				log.Error("output/server", "Failed to accept a connection: "+err.Error())
			}
			server.Conns.Map[server.Conns.FirstFreeSlot] = conn
			server.Conns.M.Unlock()
			go server.Conns.FindFirstAvailableSlot()
			go server.sendLastEvents()
		}
	}()
}

func ReceiveEvent(e event.Event) {
	server.handleEvent(e)
}

func (s *Server) handleEvent(e event.Event) {
	if e.Type == event.EventTypeServerClosed {
		log.Info("output/server", "Start closing the server...")
		for i, c := range s.Conns.Map {
			c.Close()
			delete(s.Conns.Map, i)
		}
		if server.Protocol == "unix" {
			os.Remove(s.ListenPath)
		}
		return
	}

	switch e.Type {
	case event.EventTypePlayerChanged:
		s.LastEvents[0] = e
	case event.EventTypeRateChanged:
		s.LastEvents[1] = e
	case event.EventTypeSongChanged:
		s.LastEvents[2] = e
	case event.EventTypePlaybackStatusChanged:
		s.LastEvents[3] = e
	case event.EventTypeLyricsStateChanged:
		s.LastEvents[4] = e
	case event.EventTypeLyricsChanged:
		s.LastEvents[5] = e
	case event.EventTypeActiveLyricChanged:
		s.LastEvents[6] = e
	case event.EventTypeOverwriteRequired:
		s.LastEvents[7] = e
	}

	s.sendEvent(e)
}

func CloseServer() {
	server.handleEvent(event.Event{Type: event.EventTypeServerClosed, Data: event.EventTypeServerClosedData{}})
}

package server

import (
	"errors"
	"lrcsnc/internal/output/pkg/event"
	"lrcsnc/internal/pkg/global"
	"lrcsnc/internal/pkg/log"
	"net"
	"os"
	"strings"
)

type Server struct {
	Protocol    string
	ListenPath  string
	NetListener net.Listener
	Conns       Connections
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
	}

	go func() {
		for {
			server.Conns.M.Lock()
			conn, err := server.NetListener.Accept()
			server.Conns.M.Lock()
			if err != nil {
				log.Error("output/server", "Failed to accept a connection: "+err.Error())
			}
			server.Conns.Map[server.Conns.FirstFreeSlot] = conn
			server.Conns.M.Unlock()
			go server.Conns.FindFirstAvailableSlot()
		}
	}()
}

func ReceiveEvent(e event.Event) {
	server.handleEvent(e)
}

func (s *Server) handleEvent(e event.Event) {
	switch e.Type {
	default:
		s.sendEvent(e)
	case event.EventTypeServerClosed:
		log.Info("output/server", "Start closing the server...")
		for i, c := range s.Conns.Map {
			c.Close()
			delete(s.Conns.Map, i)
		}
		if server.Protocol == "unix" {
			os.Remove(s.ListenPath)
		}
	}
}

func CloseServer() {
	server.handleEvent(event.Event{Type: event.EventTypeServerClosed, Data: event.EventTypeServerClosedData{}})
}

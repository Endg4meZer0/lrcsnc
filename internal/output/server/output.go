package server

import (
	"encoding/json"
	"errors"
	"io"
	"lrcsnc/internal/output/client"
	"lrcsnc/internal/output/pkg/event"
	"lrcsnc/internal/pkg/log"
	"net"
	"syscall"
	"time"
)

func (s *Server) sendEvent(e event.Event) {
	if s.Protocol == "" {
		client.ReceiveEvent(e)
		return
	}

	d, err := json.Marshal(map[event.EventType]any{e.Type: e.Data})
	if err != nil {
		log.Fatal("output/server", "SendEventAsync returned error during marshalling data. What's up? Error: "+err.Error())
	}
	for i, c := range s.Conns.Map {
		c.SetWriteDeadline(time.Now().Add(500 * time.Millisecond))
		_, err := c.Write(append(d, byte('\n')))
		if errors.Is(err, syscall.EPIPE) || errors.Is(err, net.ErrClosed) || errors.Is(err, io.EOF) {
			delete(s.Conns.Map, i)
		}
		if errors.Is(err, syscall.ETIMEDOUT) {
			c.Close()
			delete(s.Conns.Map, i)
		}
	}
}

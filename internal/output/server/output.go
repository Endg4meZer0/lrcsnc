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

func (s *server) sendEvent(e event.Event) {
	if s.protocol == "" {
		client.ReceiveEvent(e)
		return
	}

	d, err := json.Marshal(e)
	if err != nil {
		log.Fatal("output/server", "SendEventAsync returned error during marshalling data. What's up? Error: "+err.Error())
	}
	for i, c := range s.conns.Map {
		c.SetWriteDeadline(time.Now().Add(500 * time.Millisecond))
		_, err := c.Write(append(d, byte('\n')))
		if errors.Is(err, syscall.EPIPE) || errors.Is(err, net.ErrClosed) || errors.Is(err, io.EOF) {
			log.Debug("output/server", "A client disconnected.")
			delete(s.conns.Map, i)
		}
		if errors.Is(err, syscall.ETIMEDOUT) {
			log.Debug("output/server", "A client timed out. Disconnecting...")
			c.Close()
			delete(s.conns.Map, i)
		}
	}
}

func (s *server) sendLastEventsTo(index uint) {
	c := s.conns.Map[index]
	for _, e := range s.lastEvents {
		d, err := json.Marshal(e)
		if err != nil {
			log.Fatal("output/server", "SendEventAsync returned error during marshalling data. What's up? Error: "+err.Error())
		}
		c.SetWriteDeadline(time.Now().Add(1000 * time.Millisecond))
		_, err = c.Write(append(d, byte('\n')))
		if errors.Is(err, syscall.EPIPE) || errors.Is(err, net.ErrClosed) || errors.Is(err, io.EOF) {
			delete(s.conns.Map, index)
		}
		if errors.Is(err, syscall.ETIMEDOUT) {
			c.Close()
			delete(s.conns.Map, index)
		}
	}
}

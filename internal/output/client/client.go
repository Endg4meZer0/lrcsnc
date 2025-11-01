package client

import (
	"bufio"
	"encoding/json"
	"fmt"
	"lrcsnc/internal/output/pkg/event"
	"lrcsnc/internal/pkg/global"
	"lrcsnc/internal/pkg/log"
	"lrcsnc/internal/pkg/util/dynreplacer"
	"net"
	"os"
)

type client struct {
	conn net.Conn

	activeText        string
	pendingText       string
	outputDestination *os.File
	outputPath        string
	tempFile          string
	overwrite         string

	templateReplacer   *dynreplacer.DynamicReplacer
	instrumentalActive bool
}

var cl client

func InitClient() {
	global.Config.M.Lock()

	log.Info("output/client", "Initializing client...")

	if global.Config.C.Net.Protocol == "" {
		log.Info("output/client", "Using internal communication (standalone mode).")
		cl = client{conn: nil}
	} else {
		log.Info("output/client", "Trying to connect to server...")
		l, err := net.Dial(global.Config.C.Net.Protocol, global.Config.C.Net.ListenAt)
		if err != nil {
			log.Fatal("output/client", fmt.Sprintf("Failed to set up a client listening to %s using protocol %s: %s", global.Config.C.Net.ListenAt, global.Config.C.Net.Protocol, err.Error()))
		}
		log.Info("output/client", "Successfully connected to server at "+global.Config.C.Net.Protocol+"://"+global.Config.C.Net.ListenAt)

		cl = client{
			conn: l,
		}

		// Start the event reader
		go func() {
			reader := bufio.NewReader(cl.conn)
			for {
				str, err := reader.ReadString('\n')
				if err != nil {
					log.Fatal("output/client", "Reading data from server returned an error: "+err.Error())
				}
				log.Debug("output/client", "Received "+str+" from server.")

				var e event.Event
				err = json.Unmarshal([]byte(str), &e)
				if err != nil {
					log.Fatal("output/client", "Unexpected data received: "+str)
				}

				go ReceiveEvent(e)
			}
		}()
	}

	// Initial check for config's output destination
	if global.Config.C.Client.Destination != "stdout" {
		cl.changeOutput()
	} else {
		cl.outputDestination = os.Stdout
		cl.outputPath = "/dev/stdout"
	}
	global.Config.M.Unlock()

	log.Info("output/client", "Output initialized, getting ready to roll.")

	cl.setTemplateReplacer()

	log.Info("output/client", "Client is ready.")
}

func ReceiveEvent(e event.Event) {
	log.Debug("output/client", fmt.Sprintf("Received new event: %v", e))
	cl.handleEvent(e)
}

func Close() {
	cl.close()
}

func (c *client) handleEvent(e event.Event) {
	switch e.Type {
	case event.EventTypeActiveLyricChanged:
		c.handleActiveLyricChanged(e.Data.(event.EventTypeActiveLyricChangedData))
	case event.EventTypeSongChanged:
		c.handleSongChanged(e.Data.(event.EventTypeSongChangedData))
	case event.EventTypePlayerChanged:
		c.handlePlayerChanged(e.Data.(event.EventTypePlayerChangedData))
	case event.EventTypePlaybackStatusChanged:
		c.handlePlaybackStatusChanged(e.Data.(event.EventTypePlaybackStatusChangedData))
	case event.EventTypeRateChanged:
		c.handleRateChanged(e.Data.(event.EventTypeRateChangedData))
	case event.EventTypeLyricsStateChanged:
		c.handleLyricsStateChanged(e.Data.(event.EventTypeLyricsStateChangedData))
	case event.EventTypeLyricsChanged:
		c.handleLyricsChanged(e.Data.(event.EventTypeLyricsChangedData))
	case event.EventTypeOverwriteRequired:
		c.handleOverwriteRequired(e.Data.(event.EventTypeOverwriteRequiredData))
	case event.EventTypeServerClosed:
		c.handleServerClosed(e.Data.(event.EventTypeServerClosedData))
	case event.EventTypeConfigReloaded:
		global.Config.M.Lock()
		if global.Config.C.Client.Destination != "stdout" && global.Config.C.Client.Destination != cl.outputPath {
			cl.changeOutput()
		} else if global.Config.C.Client.Destination == "stdout" {
			cl.outputDestination = os.Stdout
			cl.outputPath = "/dev/stdout"
		}
		global.Config.M.Unlock()
		log.Debug("output/client", "Updated ")
	}
}

func (c *client) close() {
	if c.outputDestination != nil && !c.isOutputStd() {
		c.outputDestination.Close()
	}
}

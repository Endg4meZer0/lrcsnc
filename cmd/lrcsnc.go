package cmd

import (
	"os"
	"os/signal"
	"syscall"

	"lrcsnc/internal/config"
	"lrcsnc/internal/mpris"
	"lrcsnc/internal/output/client"
	"lrcsnc/internal/output/server"
	"lrcsnc/internal/pkg/global"
	"lrcsnc/internal/pkg/log"
	"lrcsnc/internal/setup"
	"lrcsnc/internal/sync"
)

func Start() {
	// Handle all the general setup...
	setup.Setup()
	// ...and check for dependencies
	setup.CheckDependencies()

	// Start the USR1 signal listener for config updates
	// TODO: replace with live file watcher
	usr1Sig := make(chan os.Signal, 1)
	signal.Notify(usr1Sig, syscall.SIGUSR1)

	go func() {
		for {
			<-usr1Sig
			config.Update()
		}
	}()

	// Initialize the client
	// (only if not explicitly launched in server mode)
	if !global.Config.C.Net.IsServer {
		client.InitClient()
		defer client.Close()
	}

	// Initialize the server
	// (only if not explicitly launched in client mode)
	if global.Config.C.Net.IsServer || global.Config.C.Net.Protocol == "" {
		server.InitServer()
		defer server.CloseServer()

		// Deploy the main watchers
		sync.Start()

		// Initialize the player listener session
		err := mpris.Connect()
		if err != nil {
			log.Fatal("cmd", "Error when configuring MPRIS. Check logs for more info.")
		}
		defer mpris.Disconnect()
	}

	exitSigs := make(chan os.Signal, 1)
	signal.Notify(exitSigs, syscall.SIGINT, syscall.SIGTERM)

	log.Info("cmd", "lrcsnc has started.")

	<-exitSigs
	log.Info("cmd", "Exit signal received, bye!")
}

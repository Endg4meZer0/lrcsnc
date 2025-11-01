// go::build forceposix

package setup

import (
	"fmt"
	"os"

	"lrcsnc/internal/config"
	"lrcsnc/internal/pkg/global"
	"lrcsnc/internal/pkg/log"
	"lrcsnc/internal/pkg/types"

	"github.com/jessevdk/go-flags"
)

// Version is linked through -X (check Makefile)
var version = "dev"

var opts struct {
	NoLog                bool   `long:"no-log" description:"Disables logging. --log-file and --log-level are ignored if this flag is set." env:"LRCSNC_NO_LOG"`
	LogPath              string `long:"log-file" description:"Sets the log file path to use." default:"$HOME/.local/state/lrcsnc/log" env:"LRCSNC_LOG_FILE"`
	LogLevel             string `long:"log-level" description:"Sets the log level used by logger." default:"info" choice:"debug" choice:"info" choice:"warn" choice:"error" choice:"fatal" env:"LRCSNC_LOG_LEVEL"`
	ConfigPath           string `short:"c" long:"config" description:"Sets the config file to use" env:"LRCSNC_CONFIG"`
	ConfigGeneratePath   string `long:"config-gen" description:"Generates a config from the default one and places it in the specified path, then exits" env:"LRCSNC_CONFIG_GEN"`
	IsServer             bool   `short:"s" long:"server" description:"Start lrcsnc in server mode. Requires protocol and listen path to be set." env:"LRCSNC_SERVER"`
	Protocol             string `short:"p" long:"protocol" description:"Sets the protocol to use for net connections." choice:"unix" choice:"tcp" choice:"tcp4" choice:"tcp6" env:"LRCSNC_PROTOCOL"`
	ListenAt             string `short:"l" long:"listen" description:"Sets the port/socket path to listen at. If launched as server, starts serving on this path instead. Format: <path> for Unix or <host>:<port> for TCP" env:"LRCSNC_LISTEN_AT"`
	CacheDirectory       string `short:"d" long:"cache-dir" description:"Sets the cache directory" env:"LRCSNC_CACHE"`
	ClientOutputFilePath string `short:"o" long:"output" description:"Sets an output file to use instead of standard output when in standalone or client mode" env:"LRCSNC_OUTPUT"`
	DisplayVersion       bool   `short:"v" long:"version" description:"Display the version"`
}

// Setup parses the command line flags (or their environment variable equivalents)
// and sets up the logger, config and some other settings.
func Setup() {
	_, err := flags.Parse(&opts)
	if flags.WroteHelp(err) {
		os.Exit(0)
	}
	if err != nil {
		log.Fatal("setup", err.Error())
	}

	// Generic flags: -v...
	if opts.DisplayVersion {
		fmt.Println(version)
		os.Exit(0)
	}

	if opts.ConfigGeneratePath != "" {
		if err := config.CopyDefaultTo(opts.ConfigGeneratePath); err != nil {
			log.Fatal("setup", fmt.Sprintf("Failed to generate config to '%s': %s", opts.ConfigGeneratePath, err.Error()))
		}
		os.Exit(0)
	}

	// Logger flags: --no-log, --log-file, --log-level
	if !opts.NoLog { // if the --no-log flag is NOT set, then initialize the logger
		log.SetDestination(opts.LogPath)

		l := types.LogLevelType(opts.LogLevel)
		if l.ToInt() == -1 {
			log.Fatal("setup", fmt.Sprintf("The provided log level (%v) is invalid. Possible values: debug, info, warn, error, fatal.", opts.LogLevel))
		}
		log.SetLevel(l)

		log.Init()
	}

	// Try to read config from the provided path.
	// Failing should crash the app.
	if opts.ConfigPath != "" {
		log.Info("setup", fmt.Sprintf("Trying to read config from the provided path (%v)...", opts.ConfigPath))
		err = config.Read(opts.ConfigPath)
		if err != nil {
			log.Fatal("setup", fmt.Sprintf("The provided config from path %v failed to load; the launch will not continue. Check logs for more information.", opts.ConfigPath))
		}
	} else {
		// If the config path flag is not set, try to read other configs
		log.Info("setup", "Trying to read user-wide config...")
		if err := config.ReadUserWide(); err != nil {
			log.Info("setup", "The user-wide config will be ignored.")
			log.Info("setup", "Trying to read the system-wide config...")
			if err := config.ReadSystemWide(); err != nil {
				log.Info("setup", "The system-wide config will be ignored.")
				log.Info("setup", "Will be using the default config.")
				if err := config.ReadDefault(); err != nil {
					log.Fatal("setup", "The default config is invalid; the launch will not continue.")
				}
			}
		}
	}

	// Explicitly set the protocol if the flag is not empty
	if opts.Protocol != "" {
		global.Config.C.Net.Protocol = opts.Protocol
	}

	// Explicitly set the path/host to listen to if the flag is not empty
	if opts.ListenAt != "" {
		global.Config.C.Net.ListenAt = opts.ListenAt
	}

	// Enable server mode if the flag is set
	if opts.IsServer {
		global.Config.C.Net.IsServer = opts.IsServer
	}

	// Explicitly change cache directory for this app instance if the flag is set
	if opts.CacheDirectory != "" {
		// Only ignore the directory if there is an error and it is not a "not exist" error
		// (e.g. if the path doesn't lead to an actual directory, or if there are no permissions to read it)
		if _, err := os.ReadDir(os.ExpandEnv(opts.CacheDirectory)); err != nil && !os.IsNotExist(err) {
			log.Error("setup", fmt.Sprintf("The provided cache directory (%v) is invalid and will be ignored.", opts.CacheDirectory))
		} else {
			global.Config.C.Cache.Dir = opts.CacheDirectory
		}
	}

	// If lrcsnc is not in server mode, explicitly set the output file path for this app instance if the flag is set
	if opts.ClientOutputFilePath != "" && !global.Config.C.Net.IsServer {
		// We'll try to write to or create the file on the specified path first to ensure it is valid
		if k, err := os.OpenFile(opts.ClientOutputFilePath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644); err != nil {
			log.Error("setup", fmt.Sprintf("The provided output file path (%v) is invalid and will be ignored. Error: %v", opts.ClientOutputFilePath, err))
		} else {
			k.Close()
			global.Config.C.Client.Destination = opts.ClientOutputFilePath
			// The output is not initialized yet, so no events are sent to the output controller
		}
	}
}

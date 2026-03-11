package net

type Config struct {
	// IsServer is a boolean that exists mostly for the purposes
	// of making an easier way to code the switch to server mode.
	// It can be set to true in the config file as well:
	// though it is not recommended to do that in your main config file.
	IsServer bool `toml:"is-server"`
	// Protocol represents the protocol used to connect or to start a server
	// on. Can use values available in net.Listen (so "unix", "tcp" and some other variations)
	// and empty for standalone mode.
	Protocol string `toml:"protocol"`
	// ListenAt represents the path of a server to listen to.
	// If lrcsnc fails to connect to a server at this path, then it crashes.
	ListenAt string `toml:"listen-at"`
}

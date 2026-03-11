package configStruct

import (
	"lrcsnc/internal/config/struct/cache"
	"lrcsnc/internal/config/struct/client"
	"lrcsnc/internal/config/struct/lyrics"
	"lrcsnc/internal/config/struct/net"
	"lrcsnc/internal/config/struct/player"
)

type Config struct {
	// Player config is for player related things. Currently it is used
	// for specifying included/excluded players for the watcher.
	Player player.Config `toml:"player"`
	// Net config is for client-server mode settings.
	Net net.Config `toml:"net"`
	// Lyrics config currently has stuff to do with lyrics providers,
	// time offset and romanization
	Lyrics lyrics.Config `toml:"lyrics"`
	// Cache config has an "enabled" toggle, dir path and life span
	Cache cache.Config `toml:"cache"`
	// Client config has lots of personalized settings for the
	// lrcsnc's native client that is used in standalone and client modes.
	// Stuff like what kind of output should be going from lrcsnc is described here.
	Client client.Config `toml:"client"`
}

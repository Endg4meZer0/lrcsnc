package player

type Config struct {
	IncludedPlayers []string `toml:"included-players"`
	ExcludedPlayers []string `toml:"excluded-players"`
}

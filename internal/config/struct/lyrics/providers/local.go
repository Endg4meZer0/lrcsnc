package providers

type LocalProviderConfig struct {
	TryEmbeddedFirst bool `toml:"try-embedded-first"`
	CacheInternally  bool `toml:"cache-internally"`
}

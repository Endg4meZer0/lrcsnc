package lyrics

import (
	"lrcsnc/internal/pkg/structs/config/lyrics/romanization"
	"lrcsnc/internal/pkg/types"
)

type Config struct {
	Provider     types.LyricsProviderType `toml:"provider"`
	TimingOffset float64                  `toml:"timing-offset"`
	Romanization romanization.Config      `toml:"romanization"`
}

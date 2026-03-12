package lyrics

import (
	"lrcsnc/internal/config/struct/lyrics/providers"
	"lrcsnc/internal/config/struct/lyrics/romanization"
	"lrcsnc/internal/pkg/types"
)

type Config struct {
	Providers    []types.LyricsProviderType `toml:"providers"`
	TimingOffset float64                    `toml:"timing-offset"`
	Romanization romanization.Config        `toml:"romanization"`

	LocalProviderConfig  providers.LocalProviderConfig  `toml:"providers-local"`
	LrcLibProviderConfig providers.LrcLibProviderConfig `toml:"providers-lrclib"`
}

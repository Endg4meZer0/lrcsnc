package formats

import (
	"lrcsnc/internal/lyrics/formats/lrc"
	"lrcsnc/internal/pkg/types"
)

var Formats = map[types.LyricsFormat]LyricsFormat{
	types.LyricsFormatLrcPlain:  lrc.LyricsFormatLrcPlain{},
	types.LyricsFormatLrcSynced: lrc.LyricsFormatLrcSynced{},
}

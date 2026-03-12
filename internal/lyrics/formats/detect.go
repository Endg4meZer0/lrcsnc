package formats

import (
	"lrcsnc/internal/pkg/types"
	"regexp"
)

var timingRegexp = regexp.MustCompile(`(\[[0-9]{2}:[0-9]{2}.[0-9]{1,3}])+`)

func DetectFormat(data string) types.LyricsFormat {
	if timingRegexp.MatchString(data) {
		return types.LyricsFormatLrcSynced
	} else {
		return types.LyricsFormatLrcPlain
	}
}

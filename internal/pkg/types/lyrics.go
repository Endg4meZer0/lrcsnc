package types

type LyricsState byte

const (
	LyricsStateSynced       LyricsState = 0
	LyricsStatePlain        LyricsState = 1
	LyricsStateInstrumental LyricsState = 2
	LyricsStateNotFound     LyricsState = 3
	LyricsStateLoading      LyricsState = 4
	LyricsStateUnknown      LyricsState = 5
)

var lyricsStateStrings = map[LyricsState]string{
	LyricsStateSynced:       "synced",
	LyricsStatePlain:        "plain",
	LyricsStateInstrumental: "instrumental",
	LyricsStateNotFound:     "not-found",
	LyricsStateLoading:      "loading",
	LyricsStateUnknown:      "unknown",
}

func (l LyricsState) String() string {
	return lyricsStateStrings[l]
}

func (l LyricsState) ToCacheStoreCondition() CacheStoreConditionType {
	switch l {
	case LyricsStateSynced:
		return CacheStoreConditionSynced
	case LyricsStatePlain:
		return CacheStoreConditionPlain
	case LyricsStateInstrumental:
		return CacheStoreConditionInstrumental
	default:
		return CacheStoreConditionNone
	}
}

package util

import (
	playerStructs "lrcsnc/internal/pkg/structs/player"
	"lrcsnc/internal/pkg/util"
	"testing"
)

func TestCalcMultiplier(t *testing.T) {
	test1 := []playerStructs.Lyric{
		{
			Timing: 0.11,
			Text:   "BUGAGA",
		},
		{
			Timing: 0.31,
			Text:   "BUGAGA",
		},
		{
			Timing: 0.651,
			Text:   "BUGAGA",
		},
		{
			Timing: 1.222,
			Text:   "",
		},
	}

	if n := util.CalculateMultiplier(test1, 2); n != 3 {
		t.Errorf("[tests/util/TestCalcMultiplier] ERROR: mismatch, got %v instead of 3", n)
	}
	if n := util.CalculateMultiplier(test1, 1); n != 2 {
		t.Errorf("[tests/util/TestCalcMultiplier] ERROR: mismatch, got %v instead of 2", n)
	}
	if n := util.CalculateMultiplier(test1, 3); n != 0 {
		t.Errorf("[tests/util/TestCalcMultiplier] ERROR: mismatch, got %v instead of 0", n)
	}
}

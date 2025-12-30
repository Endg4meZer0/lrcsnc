package player

import (
	playerStructs "lrcsnc/internal/pkg/structs/player"
	"testing"
)

func TestCalcMultiplier(t *testing.T) {
	test1 := playerStructs.Lyrics{
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

	if n := test1.CalculateMultiplierFor(2); n != 3 {
		t.Errorf("[tests/pkg/structs/player/TestCalcMultiplier] ERROR: mismatch, got %v instead of 3", n)
	}
	if n := test1.CalculateMultiplierFor(1); n != 2 {
		t.Errorf("[tests/pkg/structs/player/TestCalcMultiplier] ERROR: mismatch, got %v instead of 2", n)
	}
	if n := test1.CalculateMultiplierFor(3); n != 0 {
		t.Errorf("[tests/pkg/structs/player/TestCalcMultiplier] ERROR: mismatch, got %v instead of 0", n)
	}
}

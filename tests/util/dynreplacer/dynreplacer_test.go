package dynreplacer

import (
	"fmt"
	"lrcsnc/internal/pkg/global"
	"lrcsnc/internal/pkg/util/dynreplacer"
	"testing"
)

func TestDynReplacer(t *testing.T) {
	global.Player.P.Position = 123.11

	dr := dynreplacer.NewDynamicReplacer(
		map[string]func() string{
			"test":     func() string { return "TEST?!?!?!" },
			"what":     func() string { return fmt.Sprintf("huh %v", 1) },
			"position": func() string { return fmt.Sprintf("%v", global.Player.P.Position) },
		},
	)

	repl1 := dr.Replace("{test}ing {stuff} inside {what} the box {position}")
	true1 := `TEST?!?!?!ing {stuff} inside huh 1 the box 123.11`
	if repl1 != true1 {
		t.Errorf("[util/TestDynReplacer] ERROR: Received " + repl1 + " instead of " + true1)
	}

	global.Player.P.Position = 514.12
	repl2 := dr.Replace("{test}ing {what} inside {stuff} the box {position}")
	true2 := `TEST?!?!?!ing huh 1 inside {stuff} the box 514.12`
	if repl2 != true2 {
		t.Errorf("[util/TestDynReplacer] ERROR: Received " + repl2 + " instead of " + true2)
	}
}

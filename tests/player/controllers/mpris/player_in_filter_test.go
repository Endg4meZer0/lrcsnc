package mpris

import (
	configStruct "lrcsnc/internal/config/struct"
	"lrcsnc/internal/config/struct/player"
	"strings"
	"testing"
)

func TestPlayerInFilter(t *testing.T) {
	conf_incl := configStruct.Config{Player: player.Config{IncludedPlayers: []string{"fooyin", "spotify"}}}
	conf_excl := configStruct.Config{Player: player.Config{ExcludedPlayers: []string{"firefox", "zen"}}}

	player := "spotify.instance_1_123123"
	for _, includedPlayer := range conf_incl.Player.IncludedPlayers {
		if strings.Contains(player, includedPlayer) {
			t.Logf("[tests/pkg/structs/player/TestPlayerInFilter] INFO: player %v is included in the config", player)
		}
	}

	player = "firefox.instance_1_123123"
	for _, excludedPlayer := range conf_excl.Player.ExcludedPlayers {
		if strings.Contains(player, excludedPlayer) {
			t.Logf("[tests/pkg/structs/player/TestPlayerInFilter] INFO: player %v is excluded in the config", player)
		}
	}
}

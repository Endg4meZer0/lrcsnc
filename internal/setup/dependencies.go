package setup

import (
	"lrcsnc/internal/pkg/global"
)

// CheckDependencies checks if the required dependencies are installed
// and makes the appropriate adjustments to the config.
func CheckDependencies() {
	global.Config.M.Lock()
	defer global.Config.M.Unlock()

	// kakasi - for japanese romanization (no longer a thing right now, but the piece
	// of code is here for certain development purposes like I don't want to rewrite this
	// again when I need it again)

	// if _, err := exec.LookPath("kakasi"); err != nil && global.Config.C.Lyrics.Romanization.Japanese {
	// 	log.Info("setup/dependencies", "kakasi was not found; disabling Japanese romanization. If you want it back, please, install kakasi.")
	// 	global.Config.C.Lyrics.Romanization.Japanese = false
	// }
}

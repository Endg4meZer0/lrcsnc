package client

import (
	"lrcsnc/internal/pkg/global"
	"strconv"
	"strings"
)

// formatToTemplate formats the output to the Template configuration option.
func (c *client) formatToTemplate() string {
	return c.templateReplacer.Replace(global.Config.C.ClientOutput.Template)
}

// formatLyric formats the pending lyric-only string to be displayed
// in accordance with the text format configuration.
func (c *client) formatLyric(mult int) string {
	global.Config.M.Lock()
	defer global.Config.M.Unlock()
	global.Player.M.Lock()
	defer global.Player.M.Unlock()

	// The check on instrumental should happen before,
	// but here's another check just in case
	if c.pendingText == "" {
		return ""
	}

	// Calculating the multiplier value
	multiplier := ""
	if mult > 1 {
		multiplier = strings.ReplaceAll(global.Config.C.ClientOutput.Format.Multiplier, "%value%", strconv.Itoa(mult))
	}
	replacer := strings.NewReplacer(
		"%lyric%", c.pendingText,
		"%multiplier%", multiplier,
	)

	return replacer.Replace(global.Config.C.ClientOutput.Format.Lyric)
}

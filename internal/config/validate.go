package config

import (
	"fmt"
	"os"
	"path"
	"strings"

	"lrcsnc/internal/pkg/structs"
)

type ValidationError struct {
	Path    string
	Message string
	Fatal   bool
}

func (v ValidationError) Error() string {
	return v.Message
}

type ValidationErrors []ValidationError

func Validate(c *structs.Config) (errs ValidationErrors) {
	errs = make(ValidationErrors, 0)

	// Check whether protocol value is allowed
	if c.Net.Protocol != "" && c.Net.Protocol != "unix" && c.Net.Protocol != "tcp" && c.Net.Protocol != "tcp4" && c.Net.Protocol != "tcp6" {
		errs = append(errs, ValidationError{
			Path:    "net/protocol",
			Message: fmt.Sprintf("'%s' is not a valid value. Expected \"\", \"unix\", \"tcp\", \"tcp4\" or \"tcp6\".", c.Net.Protocol),
			Fatal:   true,
		})
	}

	// Check whether listen path value is allowed
	if c.Net.Protocol == "tcp" && len(strings.Split(c.Net.ListenAt, ":")) != 2 {
		errs = append(errs, ValidationError{
			Path:    "net/listen-path",
			Message: fmt.Sprintf("'%s' is not a valid value. Paired with \"tcp\" protocol, it should consist of <host>:<port>.", c.Net.ListenAt),
			Fatal:   true,
		})
	}

	// Check whether lrclib is set as the lyrics provider
	if c.Lyrics.Provider != "lrclib" {
		errs = append(errs, ValidationError{
			Path:    "lyrics/lyrics-provider",
			Message: fmt.Sprintf("'%s' is not a valid value. Allowed values are 'lrclib' (sure hope there will be more in the future)", c.Lyrics.Provider),
			Fatal:   true,
		})
	}

	// Check if client's output's destination is writeable if it's not stdout
	if !c.Net.IsServer && c.ClientOutput.Destination != "stdout" && !isPathWriteable(c.ClientOutput.Destination) {
		errs = append(errs, ValidationError{
			Path:    "client-output/destination",
			Message: fmt.Sprintf("'%s' is not a writeable path. Please make sure the path exists and is writeable", c.ClientOutput.Destination),
			Fatal:   true,
		})
	}

	// Check if the instrumental interval is set to <0.1s
	if !c.Net.IsServer && c.ClientOutput.Format.Instrumental.Interval < 0.1 {
		errs = append(errs, ValidationError{
			Path:    "client-output/format/instrumental/interval",
			Message: fmt.Sprintf("'%f' is not a valid value. Using the possible minimum instead (0.1s)", c.ClientOutput.Format.Instrumental.Interval),
			Fatal:   false,
		})
		c.ClientOutput.Format.Instrumental.Interval = 0.1
	}

	// Check if max symbols is less than 1
	if !c.Net.IsServer && c.ClientOutput.Format.Instrumental.MaxSymbols < 1 {
		errs = append(errs, ValidationError{
			Path:    "client-output/format/instrumental/max-symbols",
			Message: fmt.Sprintf("'%d' is not a valid value. Using the possible minimum instead (1)", c.ClientOutput.Format.Instrumental.MaxSymbols),
			Fatal:   false,
		})
		c.ClientOutput.Format.Instrumental.MaxSymbols = 1
	}

	return
}

func isPathWriteable(p string) bool {
	p = path.Clean(p)
	f, err := os.OpenFile(p, os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		return false
	} else {
		f.Close()
		return true
	}
}

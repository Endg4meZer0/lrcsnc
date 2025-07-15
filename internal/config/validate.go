package config

import (
	"fmt"
	"os"
	"path"

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

	// Check whether output value is allowed
	if c.Output.Type != "piped" {
		errs = append(errs, ValidationError{
			Path:    "global/output",
			Message: fmt.Sprintf("'%s' is not a valid value. Allowed values are 'piped' (will be more in the future)", c.Output.Type),
			Fatal:   true,
		})
	}

	// Check whether lrclib is set as the lyrics provider
	if c.Lyrics.Provider != "lrclib" {
		errs = append(errs, ValidationError{
			Path:    "global/lyrics-provider",
			Message: fmt.Sprintf("'%s' is not a valid value. Allowed values are 'lrclib' (sure hope there will be more in the future)", c.Lyrics.Provider),
			Fatal:   true,
		})
	}

	// Check if piped output's destination is writeable if it's not stdout
	if c.Output.Type == "piped" && c.Output.Piped.Destination != "stdout" && !isPathWriteable(c.Output.Piped.Destination) {
		errs = append(errs, ValidationError{
			Path:    "output/piped/destination",
			Message: fmt.Sprintf("'%s' is not a writeable path. Please make sure the path exists and is writeable", c.Output.Piped.Destination),
			Fatal:   true,
		})
	}

	// Check if the instrumental interval is set to <0.1s
	if c.Output.Type == "piped" && c.Output.Piped.Format.Instrumental.Interval < 0.1 {
		errs = append(errs, ValidationError{
			Path:    "output/piped/instrumental/interval",
			Message: fmt.Sprintf("'%f' is not a valid value. Using the possible minimum instead (0.1s)", c.Output.Piped.Format.Instrumental.Interval),
			Fatal:   false,
		})
		c.Output.Piped.Format.Instrumental.Interval = 0.1
	}

	// Check if max symbols is less than 1
	if c.Output.Type == "piped" && c.Output.Piped.Format.Instrumental.MaxSymbols < 1 {
		errs = append(errs, ValidationError{
			Path:    "output/piped/instrumental/max-symbols",
			Message: fmt.Sprintf("'%d' is not a valid value. Using the possible minimum instead (1)", c.Output.Piped.Format.Instrumental.MaxSymbols),
			Fatal:   false,
		})
		c.Output.Piped.Format.Instrumental.MaxSymbols = 1
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

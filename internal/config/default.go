package config

import (
	_ "embed"
	"errors"
	"os"
	"path"
	"strings"

	errs "lrcsnc/internal/pkg/errors"
	"lrcsnc/internal/pkg/global"
	"lrcsnc/internal/pkg/log"
	"lrcsnc/internal/pkg/structs"

	"github.com/pelletier/go-toml/v2"
)

//go:embed examples/config.toml
var defaultConfig []byte

func ReadDefault() error {
	var config structs.Config

	if err := toml.Unmarshal(defaultConfig, &config); err != nil {
		var decodeErr *toml.DecodeError
		if errors.As(err, &decodeErr) {
			lines := strings.Join(strings.Split(decodeErr.String(), "\n"), "\n\t")
			log.Error("config/ReadDefault", "(WHO MESSED WITH THE DEFAULT CONFIG??) Error parsing the default config file: \n\t"+lines)
			return errs.ErrConfigFileInvalid
		}
	}

	wrongs := Validate(&config)
	fatal := false
	for _, v := range wrongs {
		if v.Fatal {
			log.Error("config: "+v.Path, v.Message)
			fatal = true
		} else {
			log.Warn("config: "+v.Path, v.Message)
		}
	}

	if !fatal {
		global.Config.M.Lock()
		global.Config.C = config
		global.Config.Path = "default"
		global.Config.M.Unlock()
		postRead()
		log.Info("config/ReadDefault", "Default config loaded successfully.")
	} else {
		log.Error("config/ReadDefault", "Fatal errors in the default config were detected during validation. How? Idk.")
		return errs.ErrConfigFatalValidation
	}

	return nil
}

func CopyDefaultTo(p string) error {
	if _, err := os.ReadDir(path.Dir(p)); os.IsNotExist(err) {
		os.MkdirAll(path.Dir(p), 0755)
	} else if err != nil {
		return err
	}

	return os.WriteFile(p, defaultConfig, 0644)
}

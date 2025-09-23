package config

import (
	"errors"
	"os"
	"strings"

	"lrcsnc/internal/output"
	errs "lrcsnc/internal/pkg/errors"
	"lrcsnc/internal/pkg/global"
	"lrcsnc/internal/pkg/log"
	"lrcsnc/internal/pkg/structs"

	"github.com/pelletier/go-toml/v2"
)

func Parse(configFile []byte) error {
	var config structs.Config

	if err := toml.Unmarshal(configFile, &config); err != nil {
		var decodeErr *toml.DecodeError
		if errors.As(err, &decodeErr) {
			lines := strings.Join(strings.Split(decodeErr.String(), "\n"), "\n\t")
			log.Error("config/Parse", "Error parsing the config file: \n\t"+lines)
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

	if fatal {
		log.Error("config/Parse", "Fatal errors in the config were detected during validation.")
		return errs.ErrConfigFatalValidation
	}

	global.Config.M.Lock()
	global.Config.C = config
	global.Config.M.Unlock()
	log.Info("config/Parse", "Config file loaded successfully.")

	return nil
}

func Read(path string) error {
	if _, err := os.Stat(os.ExpandEnv(path)); os.IsNotExist(err) {
		log.Error("config/Read", "Config file does not exist or is unreachable.")
		return errs.ErrFileUnreachable
	}

	configFile, err := os.ReadFile(os.ExpandEnv(path))
	if err != nil {
		log.Error("config/Read", "Config file is reachable, but unreadable.")
		return errs.ErrFileUnreadable
	}

	err = Parse(configFile)
	if err != nil {
		return err
	}

	global.Config.M.Lock()
	global.Config.Path = path
	global.Config.M.Unlock()

	return nil
}

func ReadUserWide() error {
	userConfigDir, err := os.UserConfigDir()
	if err != nil {
		return err
	}

	return Read(userConfigDir + "/lrcsnc/config.toml")
}

func ReadSystemWide() error {
	return Read("/etc/lrcsnc/config.toml")
}

func Update() {
	if global.Config.Path == "default" {
		return
	}

	if err := Read(global.Config.Path); err != nil {
		switch {
		case errors.Is(err, errs.ErrFileUnreachable):
			log.Error("config/Update", "The config file is now unreachable. The configuration will remain the same until restart or until the config file reappears.")
		default:
			log.Error("config/Update", "Unknown error: "+err.Error())
		}
	} else {
		output.Controllers[global.Config.C.Output.Type].OnConfigUpdate()
	}
}

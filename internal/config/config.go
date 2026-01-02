package config

import (
	"errors"
	"os"
	"strings"

	errs "lrcsnc/internal/config/errors"
	"lrcsnc/internal/output/client"
	"lrcsnc/internal/output/pkg/event"
	genericErrs "lrcsnc/internal/pkg/errors"
	"lrcsnc/internal/pkg/global"
	"lrcsnc/internal/pkg/log"
	configStruct "lrcsnc/internal/pkg/structs/config"

	"github.com/pelletier/go-toml/v2"
)

func Parse(configFile []byte) error {
	var config configStruct.Config

	if err := toml.Unmarshal(configFile, &config); err != nil {
		var decodeErr *toml.DecodeError
		if errors.As(err, &decodeErr) {
			lines := strings.Join(strings.Split(decodeErr.String(), "\n"), "\n\t")
			log.Error("config/Parse", "Error parsing the config file: \n\t"+lines)
			return errs.FileInvalid
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
		return errs.FatalValidation
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
		return genericErrs.FileUnreachable
	}

	configFile, err := os.ReadFile(os.ExpandEnv(path))
	if err != nil {
		log.Error("config/Read", "Config file is reachable, but unreadable.")
		return genericErrs.FileUnreadable
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
		case errors.Is(err, genericErrs.FileUnreachable):
			log.Error("config/Update", "The config file is now unreachable. The configuration will remain the same until restart or until the config file reappears.")
		default:
			log.Error("config/Update", "Unknown error: "+err.Error())
		}
	}

	if !global.Config.C.Net.IsServer {
		client.ReceiveEvent(event.Event{
			Type: event.EventTypeConfigReloaded,
			Data: nil,
		})
	}
}

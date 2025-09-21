package config

import (
	_ "embed"
	"os"
	"path"
)

//go:embed examples/config.toml
var defaultConfig []byte

func ReadDefault() error {
	return Parse(defaultConfig)
}

func CopyDefaultTo(p string) error {
	if _, err := os.ReadDir(path.Dir(p)); os.IsNotExist(err) {
		os.MkdirAll(path.Dir(p), 0755)
	} else if err != nil {
		return err
	}

	return os.WriteFile(p, defaultConfig, 0644)
}

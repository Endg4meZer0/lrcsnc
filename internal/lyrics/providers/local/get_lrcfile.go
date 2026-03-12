package local

import (
	"os"
	"path"
	"strings"
)

func getFromLrcFile(songFilePath string) (string, error) {
	songExt := path.Ext(songFilePath)
	songFilePath, _ = strings.CutSuffix(songFilePath, songExt)

	data, err := os.ReadFile(songFilePath + ".lrc")
	if err != nil {
		return "", err
	}
	return string(data), nil
}

package local

import (
	"os"

	"github.com/dhowden/tag"
)

func getFromTags(songFilePath string) (string, error) {
	f, err := os.OpenFile(songFilePath, os.O_RDONLY, 0o644)
	if err != nil {
		return "", err
	}
	md, err := tag.ReadFrom(f)
	if err != nil {
		return "", err
	}
	return md.Lyrics(), nil
}

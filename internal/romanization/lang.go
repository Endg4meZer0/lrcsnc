package romanization

import (
	"lrcsnc/internal/romanization/langs"
)

type Romanizer interface {
	Romanize(string) string
}

var Romanizers map[Language]Romanizer = map[Language]Romanizer{
	LanguageKorean:  langs.RomanizerKr{},
	LanguageChinese: langs.RomanizerZh{},
}

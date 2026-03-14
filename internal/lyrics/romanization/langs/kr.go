package langs

import (
	kr "github.com/srevinsaju/korean-romanizer-go"
)

type RomanizerKr struct{}

func (RomanizerKr) Romanize(str string) string {
	r := kr.NewRomanizer(str)
	return r.Romanize()
}

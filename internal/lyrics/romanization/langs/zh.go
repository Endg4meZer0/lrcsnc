package langs

import (
	"strings"

	zh "github.com/mozillazg/go-pinyin"
)

type RomanizerZh struct{}

func (RomanizerZh) Romanize(p string) string {
	out := strings.Builder{}
	for _, r := range p {
		pinyin := zh.Pinyin(string(r), zh.Args{Style: zh.Tone})
		if len(pinyin) != 0 && len(pinyin[0]) != 0 {
			out.WriteString(pinyin[0][0])
		} else {
			out.WriteRune(r)
		}
	}
	return out.String()
}

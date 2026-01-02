package romanization_test

import (
	"lrcsnc/internal/pkg/global"
	playerStructs "lrcsnc/internal/pkg/structs/player"
	"lrcsnc/internal/romanization"
	"testing"
)

func TestRomanize(t *testing.T) {
	global.Config.C.Lyrics.Romanization.Chinese = true
	global.Config.C.Lyrics.Romanization.Korean = true

	krLyrics := []playerStructs.Lyric{{Text: "어? 나한테 다가오니?"}}
	zhLyrics := []playerStructs.Lyric{{Text: "哦？你在接近我吗？"}}
	romanLyrics := []playerStructs.Lyric{{Text: "france?!?"}}

	romanization.Romanize(krLyrics)
	romanization.Romanize(zhLyrics)
	romanization.Romanize(romanLyrics)

	rightAnswerKorean := []playerStructs.Lyric{{Text: "Eo? Nahante dagaoni?"}}
	rightAnswerChinese := []playerStructs.Lyric{{Text: "Ó? Nǐzàijiējìnwǒma?"}}
	rightAnswerDefault := []playerStructs.Lyric{{Text: "france?!?"}}

	if krLyrics[0] != rightAnswerKorean[0] ||
		zhLyrics[0] != rightAnswerChinese[0] ||
		romanLyrics[0] != rightAnswerDefault[0] {
		t.Errorf(
			"[tests/romanization/TestRomanize] ERROR: Expected \"%v\", \"%v\" and \"%v\"; received \"%v\", \"%v\" and \"%v\"",
			rightAnswerKorean[0], rightAnswerChinese[0], rightAnswerDefault[0],
			krLyrics[0], zhLyrics[0], romanLyrics[0])
	}
}

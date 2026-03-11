package romanization_test

import (
	lyricStruct "lrcsnc/internal/lyrics/struct"
	"lrcsnc/internal/pkg/global"
	"lrcsnc/internal/romanization"
	"testing"
)

func TestRomanize(t *testing.T) {
	global.Config.C.Lyrics.Romanization.Chinese = true
	global.Config.C.Lyrics.Romanization.Korean = true

	krLyrics := []lyricStruct.Lyric{{Text: "어? 나한테 다가오니?"}}
	zhLyrics := []lyricStruct.Lyric{{Text: "哦？你在接近我吗？"}}
	romanLyrics := []lyricStruct.Lyric{{Text: "france?!?"}}

	romanization.Romanize(krLyrics)
	romanization.Romanize(zhLyrics)
	romanization.Romanize(romanLyrics)

	rightAnswerKorean := []lyricStruct.Lyric{{Text: "Eo? Nahante dagaoni?"}}
	rightAnswerChinese := []lyricStruct.Lyric{{Text: "Ó? Nǐzàijiējìnwǒma?"}}
	rightAnswerDefault := []lyricStruct.Lyric{{Text: "france?!?"}}

	if krLyrics[0] != rightAnswerKorean[0] ||
		zhLyrics[0] != rightAnswerChinese[0] ||
		romanLyrics[0] != rightAnswerDefault[0] {
		t.Errorf(
			"[tests/romanization/TestRomanize] ERROR: Expected \"%v\", \"%v\" and \"%v\"; received \"%v\", \"%v\" and \"%v\"",
			rightAnswerKorean[0], rightAnswerChinese[0], rightAnswerDefault[0],
			krLyrics[0], zhLyrics[0], romanLyrics[0])
	}
}

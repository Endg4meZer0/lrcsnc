package romanization_test

import (
	"lrcsnc/internal/pkg/global"
	playerStructs "lrcsnc/internal/pkg/structs/player"
	"lrcsnc/internal/romanization"
	"os/exec"
	"testing"
)

func TestRomanize(t *testing.T) {
	global.Config.C.Lyrics.Romanization.Japanese = true
	global.Config.C.Lyrics.Romanization.Chinese = true
	global.Config.C.Lyrics.Romanization.Korean = true

	noKakasiJpLyric := "ああ? 私に近づいてるの?"

	jpLyrics := []playerStructs.Lyric{{Text: "ああ？私に近づいてるの？"}}
	krLyrics := []playerStructs.Lyric{{Text: "어? 나한테 다가오니?"}}
	zhLyrics := []playerStructs.Lyric{{Text: "哦？你在接近我吗？"}}
	romanLyrics := []playerStructs.Lyric{{Text: "france?!?"}}
	romanization.Romanize(jpLyrics)
	romanization.Romanize(krLyrics)
	romanization.Romanize(zhLyrics)
	romanization.Romanize(romanLyrics)

	rightAnswerJapanese := []playerStructs.Lyric{{Text: "Aa? Watashi ni chikazu iteruno?"}}
	rightAnswerKorean := []playerStructs.Lyric{{Text: "Eo? Nahante dagaoni?"}}
	rightAnswerChinese := []playerStructs.Lyric{{Text: "Ó? Nǐzàijiējìnwǒma?"}}
	rightAnswerDefault := []playerStructs.Lyric{{Text: "france?!?"}}

	if _, err := exec.LookPath("kakasi"); (err == nil && jpLyrics[0] != rightAnswerJapanese[0]) ||
		(err != nil && jpLyrics[0].Text != noKakasiJpLyric) ||
		krLyrics[0] != rightAnswerKorean[0] ||
		zhLyrics[0] != rightAnswerChinese[0] ||
		romanLyrics[0] != rightAnswerDefault[0] {
		t.Errorf(
			"[tests/romanization/TestRomanize] ERROR: Expected \"%v\", \"%v\", \"%v\" and \"%v\"; received \"%v\", \"%v\", \"%v\" and \"%v\"",
			rightAnswerJapanese[0], rightAnswerKorean[0], rightAnswerChinese[0], rightAnswerDefault[0],
			jpLyrics[0], krLyrics[0], zhLyrics[0], romanLyrics[0])
	}
}

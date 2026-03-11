package lyricStruct

import (
	"encoding/json"
	"fmt"
	"lrcsnc/internal/pkg/log"
)

type Lyric struct {
	Timing float64
	Text   string
}

type Lyrics []Lyric

func (lyrics Lyrics) CalculateMultiplierFor(ind int) (value int) {
	if ind == -1 {
		return 0
	}
	if lyrics[ind].Text == "" {
		return 0
	}

	for i := ind - 1; i >= 0 && lyrics[ind].Text == lyrics[i].Text; i-- {
		value++
	}
	return value + 1
}

func (ls Lyrics) ToJSON() string {
	d, err := json.Marshal(ls)
	if err != nil {
		log.Fatal("pkg/structs/player", fmt.Sprintf("Failed to marshal lyrics; what's wrong? Error: %v", err))
	}
	return string(d)
}

package player

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

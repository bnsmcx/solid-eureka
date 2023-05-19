package bot

import (
	"sync"
	"time"
)

type Summary struct {
	Cash float64
}

type Bot struct {
	Name     string
	Cash     float64
	LongWin  int
	ShortWin int
	Mu       *sync.Mutex
	SB       map[string]Summary
}

func (b Bot) Trade() {
	for i := 0; i < 5; i++ {
		b.Cash++
		b.UpdateScoreboard()
		time.Sleep(time.Second)
	}
}

func (b Bot) UpdateScoreboard() {
	b.Mu.Lock()
	b.SB[b.Name] = Summary{b.Cash}
	b.Mu.Unlock()
}

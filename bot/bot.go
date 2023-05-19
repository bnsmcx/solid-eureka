package bot

import (
	"log"
	"solid-eureka/binance"
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
	for {
		longAvg, shortAvg, err := binance.GetAverages(b.LongWin, b.ShortWin)
		if err != nil {
			log.Println("Bot.Trade(): ", err)
		}
		log.Println(longAvg, shortAvg)
		//if short_avg > long_avg:
		//self.sell(daily_price)
		//else:
		//self.buy(daily_price)
		b.UpdateScoreboard()
		time.Sleep(time.Second)
	}
}

func (b Bot) UpdateScoreboard() {
	b.Mu.Lock()
	b.SB[b.Name] = Summary{b.Cash}
	b.Mu.Unlock()
}

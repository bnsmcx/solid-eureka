package bot

import (
	"log"
	"solid-eureka/binance"
	"solid-eureka/yahoo"
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
		time.Sleep(time.Second)
		err := binance.Ping()
		if err != nil {
			log.Println(err)
			continue
		}

		//longAvg, shortAvg, err := yahoo.GetAverages(b.LongWin, b.ShortWin)
		_, _, err = yahoo.GetAverages(b.LongWin, b.ShortWin)
		if err != nil {
			log.Println("Bot.Trade(): ", err)
		}
		//if short_avg > long_avg:
		//self.sell(daily_price)
		//else:
		//self.buy(daily_price)
		b.UpdateScoreboard()
	}
}

func (b Bot) UpdateScoreboard() {
	b.Mu.Lock()
	b.SB[b.Name] = Summary{b.Cash}
	b.Mu.Unlock()
}

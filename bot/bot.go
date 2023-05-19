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
	Shares   float64
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

		longAvg, shortAvg, err := yahoo.GetAverages(b.LongWin, b.ShortWin)
		if err != nil {
			log.Println("Bot.Trade(): ", err)
		}
		currentPrice, err := binance.GetPrice()
		if err != nil {
			log.Println("Bot.Trade(): ", err)
		}

		if shortAvg > longAvg && currentPrice > longAvg {
			b.Sell(currentPrice)
		} else if shortAvg < longAvg && currentPrice < longAvg {
			b.Buy(currentPrice)
		}
		b.UpdateScoreboard()
	}
}

func (b Bot) UpdateScoreboard() {
	b.Mu.Lock()
	b.SB[b.Name] = Summary{b.Cash}
	b.Mu.Unlock()
}

func (b Bot) Sell(price float64) {
	if b.Shares == 0 {
		return
	}
	b.Cash += b.Shares * price
	b.Shares = 0
	log.Printf("SELL: %s sold %.2f shares at $%.2f", b.Name, b.Shares, price)
}

func (b Bot) Buy(price float64) {
	if b.Cash < 1 {
		return
	}
	b.Shares += price / b.Cash
	b.Cash = 0
	log.Printf("BUY: %s bought %.2f shares at $%.2f", b.Name, b.Shares, price)
}

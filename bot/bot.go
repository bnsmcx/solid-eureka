package bot

import (
	"log"
	"math/rand"
	"solid-eureka/binance"
	"solid-eureka/yahoo"
	"sync"
	"time"
)

type Summary struct {
	Cash     float64
	Shares   float64
	TotalVal float64
}

type Bot struct {
	Name     string
	Cash     float64
	Shares   float64
	Basis    float64
	TotalVal float64
	LongWin  int
	ShortWin int
	Mu       *sync.Mutex
	SB       map[string]Summary
}

func (b Bot) Trade() {
	for {
		b.UpdateScoreboard()
		time.Sleep(time.Minute * time.Duration(rand.Intn(5)))
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
	}
}

func (b *Bot) UpdateScoreboard() {
	b.Mu.Lock()
	b.TotalVal = b.Cash + (b.Shares * b.Basis)
	b.SB[b.Name] = Summary{b.Cash, b.Shares, b.TotalVal}
	b.Mu.Unlock()
}

func (b *Bot) Sell(price float64) {
	if b.Shares == 0 {
		return
	}
	b.Cash += b.Shares * price
	b.Shares = 0
	log.Printf("SELL: %s sold %f shares at $%.2f", b.Name, b.Shares, price)
}

func (b *Bot) Buy(price float64) {
	if b.Cash < 1 {
		return
	}
	b.Shares += b.Cash / price
	b.Cash = 0
	b.Basis = price
	log.Printf("BUY: %s bought %f shares at $%.2f", b.Name, b.Shares, price)
}

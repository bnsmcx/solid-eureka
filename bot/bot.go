package bot

import (
	"log"
	"solid-eureka/test"
	"sync"
)

type Summary struct {
	Cash     float64
	Shares   float64
	TotalVal float64
}

type Bot struct {
	Name     string
	Cash     float64
	LongWin  int
	ShortWin int
	Mu       *sync.Mutex
	SB       map[string]Summary
	Shares   float64
	Basis    float64
	TotalVal float64
}

func (b *Bot) Trade() {
	day := 0
	for {
		b.UpdateScoreboard()
		if day > 365 {
			return
		}

		//longAvg, shortAvg, err := yahoo.GetAverages(b.LongWin, b.ShortWin)
		shortAvg, longAvg, longMAD, currentPrice, err := test.GetAverages(b.LongWin, b.ShortWin, day)
		_ = longMAD
		if err != nil {
			//log.Println("Bot.Trade(): ", err)
			return
		}
		b.Basis = currentPrice
		//currentPrice, err := binance.GetPrice()
		if err != nil {
			//log.Println("Bot.Trade(): ", err)
			return
		}

		if shortAvg > longAvg+longMAD && b.Shares > 0.0 {
			b.Sell(currentPrice)
		} else if shortAvg < longAvg-longMAD && b.Cash > 0.0 {
			b.Buy(currentPrice)
		}
		day++
	}
}

func (b *Bot) UpdateScoreboard() {
	b.Mu.Lock()
	b.TotalVal = b.Cash + (b.Shares * b.Basis)
	summary := Summary{b.Cash, b.Shares, b.TotalVal}
	b.SB[b.Name] = summary
	b.Mu.Unlock()
}

func (b *Bot) Sell(price float64) {
	b.Cash += b.Shares * price
	log.Printf("SELL: %s sold %f shares at $%.2f", b.Name, b.Shares, price)
	b.Shares = 0
}

func (b *Bot) Buy(price float64) {
	b.Shares += b.Cash / price
	b.Cash = 0
	b.Basis = price
	log.Printf("BUY: %s bought %f shares at $%.2f", b.Name, b.Shares, price)
}

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
	Name          string
	Cash          float64
	LongWin       int
	ShortWin      int
	Mu            *sync.Mutex
	SB            map[string]Summary
	Shares        float64
	Basis         float64
	TotalVal      float64
	EnableLogging bool
	MADMultiplier float64
}

func (b *Bot) Trade() float64 {
	tick := 0
	for {
		b.UpdateScoreboard()

		//longAvg, shortAvg, err := yahoo.GetAverages(b.LongWin, b.ShortWin)
		shortAvg, longAvg, longMAD, currentPrice, err := test.GetAverages(b.LongWin, b.ShortWin, tick)
		if err != nil {
			if b.EnableLogging {
				log.Println("Bot.Trade(): ", err)
			}
			return b.TotalVal
		}
		b.Basis = currentPrice
		//currentPrice, err := binance.GetPrice()
		if err != nil {
			if b.EnableLogging {
				log.Println("Bot.Trade(): ", err)
			}
			return b.TotalVal
		}

		if shortAvg > longAvg+(longMAD*b.MADMultiplier) && b.Shares > 0.0 {
			b.Sell(currentPrice)
		} else if shortAvg < longAvg-(longMAD*b.MADMultiplier) && b.Cash > 0.0 {
			b.Buy(currentPrice)
		}
		tick++
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
	if b.EnableLogging {
		log.Printf("SELL: %s sold %f shares at $%.2f", b.Name, b.Shares, price)
	}
	b.Shares = 0
}

func (b *Bot) Buy(price float64) {
	b.Shares += b.Cash / price
	b.Cash = 0
	b.Basis = price
	if b.EnableLogging {
		log.Printf("BUY: %s bought %f shares at $%.2f", b.Name, b.Shares, price)
	}
}

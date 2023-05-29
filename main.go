package main

import (
	"fmt"
	"log"
	"net/http"
	"solid-eureka/binance"
	"solid-eureka/bot"
	"solid-eureka/coincap"
	"solid-eureka/test"
	"strings"
	"sync"
	"time"
)

var scoreboard = make(map[string]bot.Summary)
var mu sync.Mutex

type botSetting struct {
	LongWin       int
	ShortWin      int
	MADMultiplier float64
}

var botSettings = []botSetting{
	{31, 28, 0.3}, {16, 2, 1.5}, {60, 4, 2.4},
	{31, 18, 0.3}, {6, 2, 5.5}, {20, 4, 2.4},
	{9, 2, 1.5}, {30, 4, 2.4},
}

func main() {
	//for i := 0; i < 8; i++ {
	//	startDate := time.Now().AddDate(0, 0, -(30 * (i + 1)))
	//	endDate := time.Now().AddDate(0, 0, -(30 * i))
	//	fmt.Println("\n### Testing ", startDate.String(), " through ", endDate.String())
	//	data, err := coincap.GetDataForRange(startDate.UnixMilli(), endDate.UnixMilli())
	//	if err != nil {
	//		log.Fatalln(err)
	//	}
	//	findOptimalBotSettings(data)
	//}

	runSimulation()
}

func runSimulation() {
	startDate := time.Now().AddDate(0, 0, -30)
	endDate := time.Now()
	fmt.Println("\n### Testing ", startDate.String(), " through ", endDate.String())
	data, err := coincap.GetDataForRange(startDate.UnixMilli(), endDate.UnixMilli())
	if err != nil {
		log.Fatalln(err)
	}
	test.ActiveDataSet = data
	server()
}

func findOptimalBotSettings(data []float64) {

	wg := sync.WaitGroup{}
	start := time.Now()
	var topPerformer float64
	var settings bot.Bot
	test.ActiveDataSet = data
	for i := 120; i > 0; i-- {
		for j := 45; j > 0; j-- {
			for k := 0.0; k < 6.0; k += .1 {
				wg.Add(1)
				go func(i, j int, k float64) {
					b := bot.Bot{
						Cash:     100,
						LongWin:  i,
						ShortWin: j,
						Mu:       &mu,
						SB:       scoreboard,
					}
					b.EnableLogging = false
					b.MADMultiplier = k
					b.Trade()
					mu.Lock()
					if b.TotalVal > topPerformer {
						topPerformer = b.TotalVal
						settings = b
					}
					mu.Unlock()
					wg.Done()
				}(i, j, k)
			}
		}
	}
	wg.Wait()
	fmt.Printf("\n\nTop Performer: $%.2f  (%.2f%%)\n",
		topPerformer, ((topPerformer-100)/100.0)*100)
	fmt.Printf("\tLong Window: %d\n", settings.LongWin)
	fmt.Printf("\tShort Window: %d\n", settings.ShortWin)
	fmt.Printf("\tMAD Multiplier: %.2f\n", settings.MADMultiplier)
	fmt.Println("\nExecution time: ", time.Since(start).String())

}

func server() {
	var bots = buildBots(botSettings)
	for _, b := range bots {
		b := b
		b.EnableLogging = true
		go b.Trade()
	}

	http.HandleFunc("/", handleScoreboard)
	err := http.ListenAndServe(":8000", nil)
	if err != nil {
		log.Fatalln(err)
	}
}

func buildBots(settings []botSetting) []bot.Bot {
	cash, err := binance.GetAccountBalance()
	if err != nil {
		log.Fatalln(err)
	}

	var bots []bot.Bot
	for _, s := range settings {
		b := bot.Bot{
			Name:          fmt.Sprintf("%d-%d-%.1f", s.LongWin, s.ShortWin, s.MADMultiplier),
			Cash:          cash / float64(len(settings)),
			LongWin:       s.LongWin,
			ShortWin:      s.ShortWin,
			MADMultiplier: s.MADMultiplier,
			Mu:            &mu,
			SB:            scoreboard,
		}
		bots = append(bots, b)
	}

	return bots
}

func handleScoreboard(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte(getScoreboard()))
}

func getScoreboard() string {
	var sb strings.Builder

	headerFormat := "| %*s|%*s |%*s |%*s |\n"
	entryFormat := "| %-*s|%*.2f |%*f |%*.2f |\n"
	border := " -----------------------------------------------\n"

	sb.WriteString("                   SCOREBOARD                   \n")
	sb.WriteString(border)
	sb.WriteString(fmt.Sprintf(headerFormat, 10, " ", 10, "Cash", 10, "Shares", 10, "Total"))
	sb.WriteString(border)

	portfolioValue := 0.0

	mu.Lock()
	for k, v := range scoreboard {
		portfolioValue += v.TotalVal
		entry := fmt.Sprintf(entryFormat,
			10, k,
			10, v.Cash,
			10, v.Shares,
			10, v.TotalVal)
		sb.WriteString(entry)
	}
	mu.Unlock()

	sb.WriteString(border)
	sb.WriteString(fmt.Sprintf("      Portfolio Value:    %f\n", portfolioValue))

	return sb.String()
}

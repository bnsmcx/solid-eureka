package main

import (
	"fmt"
	"log"
	"net/http"
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

var CONTROLbotSettings = []botSetting{
	{3, 1, 1.5},
	{117, 42, 1.3},
	{25, 22, 0.2},
	{39, 38, 0.1},
	{31, 29, 0.1},
	{32, 31, 0.1},
	{13, 11, 0.5},
}

var TESTbotSettings = []botSetting{
	{3, 1, 1.5},
	{117, 42, 1.3},
	{25, 22, 0.2},
	{39, 38, 0.1},
	{31, 29, 0.1},
	{32, 31, 0.1},
	{13, 11, 0.5},
}

func main() {
	//var winners []bot.Bot
	//for i := 0; i < 12; i++ {
	//	startDate := time.Now().AddDate(0, 0, -(30 * (i + 1)))
	//	endDate := time.Now().AddDate(0, 0, -(30 * i))
	//	fmt.Println("\n### Testing ", startDate.String(), " through ", endDate.String())
	//	data, err := coincap.GetDataForRange(startDate.UnixMilli(), endDate.UnixMilli())
	//	if err != nil {
	//		log.Fatalln(err)
	//	}
	//	winners = append(winners, findOptimalBotSettings(data))
	//}
	//botSettings = parseSettingsFromBots(winners)
	runSimulationFrom(6)
}

func parseSettingsFromBots(bots []bot.Bot) []botSetting {
	var settings []botSetting
	for _, b := range bots {
		s := botSetting{
			LongWin:       b.LongWin,
			ShortWin:      b.ShortWin,
			MADMultiplier: b.MADMultiplier,
		}
		settings = append(settings, s)
	}
	return settings
}

func runSimulationFrom(months int) {
	var avgReturnSum float64
	for month := months; month >= 0; month-- {
		startDate := time.Now().AddDate(0, 0, -30*(month+1))
		endDate := time.Now().AddDate(0, 0, -30*month)
		data, err := coincap.GetDataForRange(startDate.UnixMilli(), endDate.UnixMilli())
		if err != nil {
			log.Fatalln(err)
		}
		test.ActiveDataSet = data
		var bots = buildBots(CONTROLbotSettings)
		//var bots = buildBots(TESTbotSettings)

		var total float64
		var wg sync.WaitGroup
		for _, b := range bots {
			b := b
			b.EnableLogging = false
			wg.Add(1)
			go func(b bot.Bot) {
				total += b.Trade()
				wg.Done()
			}(b)
		}
		wg.Wait()
		startValue := float64(len(bots) * 100)
		returnRate := (total - startValue) / startValue
		fmt.Println(month, fmt.Sprintf("%.2f%%", returnRate))
	}
	fmt.Println("\nAverage return: ", fmt.Sprintf("%.4f%%", avgReturnSum/float64(months)))
}

func findOptimalBotSettings(data []float64) bot.Bot {
	wg := sync.WaitGroup{}
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
	fmt.Printf("%+v", settings)
	return settings
}

func server() {
	var bots = buildBots(CONTROLbotSettings)
	//var bots = buildBots(TESTbotSettings)
	for _, b := range bots {
		b := b
		b.EnableLogging = false
		go b.Trade()
	}

	http.HandleFunc("/", handleScoreboard)
	err := http.ListenAndServe(":8000", nil)
	if err != nil {
		log.Fatalln(err)
	}
}

func buildBots(settings []botSetting) []bot.Bot {
	//cash, err := binance.GetAccountBalance()
	//if err != nil {
	//	log.Fatalln(err)
	//}

	var bots []bot.Bot
	for i, s := range settings {
		b := bot.Bot{
			Name:          fmt.Sprintf("%d-%d-%.1f-%d", s.LongWin, s.ShortWin, s.MADMultiplier, i),
			Cash:          100,
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

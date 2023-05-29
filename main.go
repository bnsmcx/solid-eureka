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

func main() {
	for i := 0; i < 8; i++ {
		startDate := time.Now().AddDate(0, 0, -(180 * (i + 1)))
		endDate := time.Now().AddDate(0, 0, -(180 * i))
		fmt.Println("\n### Testing ", startDate.String(), " through ", endDate.String())
		data, err := coincap.GetDataForRange(startDate.UnixMilli(), endDate.UnixMilli())
		if err != nil {
			log.Fatalln(err)
		}
		findOptimalBotSettings(data)
	}

	//runSimulation()
}

func runSimulation() {
	startDate := time.Now().AddDate(0, 0, -365)
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
	var bots = []bot.Bot{
		{Name: "A", Cash: 100, LongWin: 5, ShortWin: 1, MADMultiplier: 1.8},
		{Name: "B", Cash: 100, LongWin: 13, ShortWin: 12, MADMultiplier: 0.3},
		{Name: "C", Cash: 100, LongWin: 25, ShortWin: 2, MADMultiplier: 3.8},
		{Name: "D", Cash: 100, LongWin: 26, ShortWin: 10, MADMultiplier: 1.3},
		{Name: "E", Cash: 100, LongWin: 34, ShortWin: 32, MADMultiplier: 0.2},
		{Name: "F", Cash: 100, LongWin: 46, ShortWin: 25, MADMultiplier: 0.8},
		{Name: "G", Cash: 100, LongWin: 90, ShortWin: 3, MADMultiplier: 4.0},
		{Name: "H", Cash: 100, LongWin: 94, ShortWin: 1, MADMultiplier: 5.6},
	}
	for _, b := range bots {
		b := b
		b.Mu = &mu
		b.SB = scoreboard
		b.EnableLogging = true
		go b.Trade()
	}

	http.HandleFunc("/", handleScoreboard)
	err := http.ListenAndServe(":8000", nil)
	if err != nil {
		log.Fatalln(err)
	}
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

	order := []string{"A", "B", "C", "D", "E", "F", "G", "H"}

	portfolioValue := 0.0

	mu.Lock()
	for _, name := range order {
		portfolioValue += scoreboard[name].TotalVal
		entry := fmt.Sprintf(entryFormat,
			10, name,
			10, scoreboard[name].Cash,
			10, scoreboard[name].Shares,
			10, scoreboard[name].TotalVal)
		sb.WriteString(entry)
	}
	mu.Unlock()

	sb.WriteString(border)
	sb.WriteString(fmt.Sprintf("      Portfolio Value:    %f\n", portfolioValue))

	return sb.String()
}

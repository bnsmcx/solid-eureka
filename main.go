package main

import (
	"fmt"
	"log"
	"net/http"
	"solid-eureka/bot"
	"solid-eureka/test"
	"strings"
	"sync"
	"time"
)

var scoreboard = make(map[string]bot.Summary)
var mu sync.Mutex

func main() {
	server()
	//findOptimalBotSettings()
}

func findOptimalBotSettings() {
	dataSets := []string{
		"2015", "2016", "2017", "2018", "2019", "2020", "2021", "2022",
	}

	wg := sync.WaitGroup{}
	start := time.Now()
	fmt.Println("" +
		"\nThe following are the optimum settings for the trading " +
		"\nalgorithm for each respective year.  " +
		"\n\nEach test started with $100.")
	for _, set := range dataSets {
		var topPerformer float64
		var settings bot.Bot
		test.ActiveDataSet = "/home/ben/repos/solid-eureka/test/test_data/" +
			set + ".csv"
		for i := 120; i > 0; i-- {
			for j := 45; j > 0; j-- {
				for k := 0.0; k < 3.5; k += .5 {
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
		fmt.Printf("\n\nTop Performer for %s: $%.2f  (%.2f%%)\n",
			set, topPerformer, (topPerformer/100.0)*100)
		fmt.Printf("\tLong Window: %d\n", settings.LongWin)
		fmt.Printf("\tShort Window: %d\n", settings.ShortWin)
		fmt.Printf("\tMAD Multiplier: %.2f\n", settings.MADMultiplier)
	}
	fmt.Println("\nExecution time: ", time.Since(start).String())

}

func server() {
	var bots = []bot.Bot{
		{Name: "2015", Cash: 100, LongWin: 31, ShortWin: 3, MADMultiplier: 3.0},
		{Name: "2016", Cash: 100, LongWin: 15, ShortWin: 3, MADMultiplier: 2.5},
		{Name: "2017", Cash: 100, LongWin: 28, ShortWin: 24, MADMultiplier: 0.5},
		{Name: "2018", Cash: 100, LongWin: 30, ShortWin: 13, MADMultiplier: 1.0},
		{Name: "2019", Cash: 100, LongWin: 12, ShortWin: 3, MADMultiplier: 2.0},
		{Name: "2020", Cash: 100, LongWin: 37, ShortWin: 9, MADMultiplier: 2.0},
		{Name: "2021", Cash: 100, LongWin: 24, ShortWin: 16, MADMultiplier: 0.0},
		{Name: "2022", Cash: 100, LongWin: 25, ShortWin: 16, MADMultiplier: 0.5},
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

	order := []string{"2015", "2016", "2017", "2018", "2019", "2020", "2021", "2022"}

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

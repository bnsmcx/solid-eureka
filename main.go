package main

import (
	"fmt"
	"log"
	"net/http"
	"solid-eureka/bot"
	"strings"
	"sync"
)

var scoreboard = make(map[string]bot.Summary)
var mu sync.Mutex

func main() {
	server()
}

func server() {
	var bots = []bot.Bot{
		{Name: "Alpha", Cash: 100, LongWin: 138, ShortWin: 14, Mu: &mu, SB: scoreboard},
		{Name: "Bravo", Cash: 100, LongWin: 59, ShortWin: 12, Mu: &mu, SB: scoreboard},
		{Name: "Charlie", Cash: 100, LongWin: 109, ShortWin: 2, Mu: &mu, SB: scoreboard},
		{Name: "Delta", Cash: 100, LongWin: 71, ShortWin: 2, Mu: &mu, SB: scoreboard},
		{Name: "Echo", Cash: 100, LongWin: 71, ShortWin: 15, Mu: &mu, SB: scoreboard},
		{Name: "Foxtrot", Cash: 100, LongWin: 109, ShortWin: 15, Mu: &mu, SB: scoreboard},
		{Name: "Golf", Cash: 100, LongWin: 61, ShortWin: 15, Mu: &mu, SB: scoreboard},
		{Name: "Average", Cash: 100, LongWin: 95, ShortWin: 11, Mu: &mu, SB: scoreboard},
	}
	for _, b := range bots {
		b := b
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

	order := []string{"Average", "Alpha", "Bravo", "Charlie", "Delta", "Echo", "Foxtrot", "Golf"}

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

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

var bots = []bot.Bot{
	{"Alpha", 100, 0, 0, 0, 138, 14, &mu, scoreboard},
	{"Bravo", 100, 0, 0, 0, 59, 12, &mu, scoreboard},
	{"Charlie", 100, 0, 0, 0, 109, 2, &mu, scoreboard},
	{"Delta", 100, 0, 0, 0, 71, 2, &mu, scoreboard},
	{"Echo", 100, 0, 0, 0, 71, 15, &mu, scoreboard},
	{"Foxtrot", 100, 0, 0, 0, 109, 15, &mu, scoreboard},
	{"Golf", 100, 0, 0, 0, 61, 15, &mu, scoreboard},
	{"Average", 100, 0, 0, 0, 95, 11, &mu, scoreboard},
}

func main() {

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

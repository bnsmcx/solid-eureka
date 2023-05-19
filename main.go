package main

import (
	"fmt"
	"log"
	"net/http"
	"solid-eureka/bot"
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
		go b.Trade()
	}

	http.HandleFunc("/", handleScoreboard)
	err := http.ListenAndServe(":8000", nil)
	if err != nil {
		log.Fatalln(err)
	}
}

func handleScoreboard(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("                   SCOREBOARD                   \n"))
	w.Write([]byte(" -----------------------------------------------\n"))
	entry := fmt.Sprintf("| %*s|%*s |%*s |%*s |\n",
		10, " ", 10, "Cash", 10, "Shares", 10, "Total")
	w.Write([]byte(entry))
	w.Write([]byte(" -----------------------------------------------\n"))
	mu.Lock()
	for k, v := range scoreboard {
		entry := fmt.Sprintf("| %-*s|%*.2f |%*f |%*.2f |\n",
			10, k, 10, v.Cash, 10, v.Shares, 10, v.TotalVal)
		w.Write([]byte(entry))
	}
	w.Write([]byte(" -----------------------------------------------\n"))
	mu.Unlock()
}

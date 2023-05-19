package main

import (
	"log"
	"solid-eureka/bot"
	"sync"
	"time"
)

var scoreboard = make(map[string]bot.Summary)
var mu sync.Mutex

var bots = []bot.Bot{
	{"Alpha", 100, 138, 14, &mu, scoreboard},
	{"Bravo", 100, 59, 12, &mu, scoreboard},
	{"Charlie", 100, 109, 2, &mu, scoreboard},
	{"Delta", 100, 71, 2, &mu, scoreboard},
	{"Echo", 100, 71, 15, &mu, scoreboard},
	{"Foxtrot", 100, 109, 15, &mu, scoreboard},
	{"Golf", 100, 61, 15, &mu, scoreboard},
	{"Average", 100, 95, 11, &mu, scoreboard},
}

func main() {
	for _, b := range bots {
		go b.Trade()
	}
	time.Sleep(time.Second * time.Duration(7))
	mu.Lock()
	for k, v := range scoreboard {
		log.Println(k, v)
	}
	mu.Unlock()
}

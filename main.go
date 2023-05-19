package main

import (
	"log"
	"solid-eureka/bot"
)

var bots = []bot.Bot{
	{"Alpha", 100, 138, 14},
	{"Bravo", 100, 59, 12},
	{"Charlie", 100, 109, 2},
	{"Delta", 100, 71, 2},
	{"Echo", 100, 71, 15},
	{"Foxtrot", 100, 109, 15},
	{"Golf", 100, 61, 15},
	{"Average", 100, 95, 11},
}

func main() {
	for _, b := range bots {
		log.Println(b)
	}
}

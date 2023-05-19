package main

type botSetting struct {
	name        string
	cash        int
	longWindow  int
	shortWindow int
}

var botSettings = []botSetting{
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
	for _, s := range botSettings {
		go bot.Launch(s.name, s.cash, s.longWindow, s.shortWindow)
	}
}

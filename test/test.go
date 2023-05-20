package test

import (
	"encoding/csv"
	"log"
	"os"
	"strconv"
)

func GetAverages(longWin, shortWin, tick int) (float64, float64, float64, error) {
	f, err := os.Open("/home/ben/repos/solid-eureka/test/test_data/2022.csv")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	r := csv.NewReader(f)

	records, err := r.ReadAll()
	if err != nil {
		log.Fatal(err)
	}

	data := make([]float64, longWin)
	for i := tick; i < longWin+tick; i++ {
		if longWin+tick > len(records) {
			break
		}
		num, err := strconv.ParseFloat(records[i][0], 64)
		if err != nil {
			log.Fatal(err)
		}
		data[i-tick] = num
	}

	longSum := 0.0
	for _, v := range data {
		longSum += v
	}

	data = data[longWin-shortWin:]
	shortSum := 0.0
	for _, v := range data {
		shortSum += v
	}

	return longSum / float64(longWin), shortSum / float64(shortWin), data[len(data)-1], nil
}

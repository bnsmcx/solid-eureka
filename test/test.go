package test

import (
	"encoding/csv"
	"errors"
	"log"
	"math"
	"os"
	"strconv"
)

var ActiveDataSet []float64

func GetDataFromFile(filename string) []float64 {
	f, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	r := csv.NewReader(f)

	records, err := r.ReadAll()
	if err != nil {
		log.Fatal(err)
	}
	return parseFloats(records)
}

func GetAverages(longWin, shortWin, day int) (float64, float64, float64, float64, error) {

	if (longWin + day) > len(ActiveDataSet) {
		return 0, 0, 0, 0, errors.New("hit end of longData set")
	}

	// Grab the relevant slice of records and convert to float
	longData := ActiveDataSet[day : day+longWin]

	// Make sure we haven't hit the end of the data and get short window slice
	if len(longData)-shortWin < 1 {
		return 0, 0, 0, 0, errors.New("hit end of longData set")
	}
	shortData := longData[len(longData)-shortWin:]

	// get the averages
	longAvg := calculateAvg(longData)
	shortAvg := calculateAvg(shortData)

	// get the Mean Absolute Deviation for the long window
	longMAD := calculateMAD(longData, longAvg)

	return shortAvg, longAvg, longMAD, longData[len(longData)-1], nil
}

func calculateMAD(data []float64, avg float64) float64 {
	var sumOfDeviations float64
	for _, v := range data {
		deviation := v - avg
		sumOfDeviations += math.Abs(deviation)
	}
	return sumOfDeviations / float64(len(data))
}

func calculateAvg(data []float64) float64 {
	var sum float64
	for _, v := range data {
		sum += v
	}
	return sum / float64(len(data))
}

func parseFloats(records [][]string) []float64 {
	var data []float64
	for _, v := range records {
		f, err := strconv.ParseFloat(v[0], 64)
		if err != nil {
			continue
		}
		data = append(data, f)
	}
	return data
}

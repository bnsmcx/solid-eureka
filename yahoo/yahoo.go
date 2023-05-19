package yahoo

import (
	"encoding/csv"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"
)

var baseURL = "https://query1.finance.yahoo.com/v7/finance/download/BTC-USD"

func getPriceIntervalURI(r int) string {
	start := time.Now()
	stop := start.AddDate(0, 0, r*-1)
	intervalInfo := fmt.Sprintf("?period1=%s&period2=%s",
		strconv.FormatInt(stop.Unix(), 10),
		strconv.FormatInt(start.Unix(), 10))
	boilerplate := "&interval=1d&events=history&includeAdjustedClose=true"
	return fmt.Sprintf("%s%s", intervalInfo, boilerplate)
}

func GetAverages(longWin int, shortWin int) (float64, float64, error) {
	longAvg, err := getAverage(longWin)
	if err != nil {
		e := fmt.Errorf("yahoo.GetAverages(): %s", err)
		return 0, 0, e
	}
	shortAvg, err := getAverage(shortWin)
	if err != nil {
		e := fmt.Errorf("yahoo.GetAverages(): %s", err)
		return 0, 0, e
	}
	return longAvg, shortAvg, nil
}

func getAverage(window int) (float64, error) {
	resp, err := http.Get(baseURL + getPriceIntervalURI(window))
	if err != nil || resp.StatusCode != 200 {
		return 0, fmt.Errorf("%s, %s", err, resp.Status)
	}
	defer resp.Body.Close()

	reader := csv.NewReader(resp.Body)
	data := make([][]string, 0, 150)
	for {
		line, err := reader.Read()
		if err != nil {
			if err == io.EOF {
				break
			} else {
				return 0, err
			}
		}
		data = append(data, line)
	}

	sum := 0.0
	length := 0.0
	for _, line := range data[1:] {
		num, err := strconv.ParseFloat(line[4], 64)
		if err != nil {
			return 0, err
		}
		sum += num
		length++
	}
	return sum / length, nil
}

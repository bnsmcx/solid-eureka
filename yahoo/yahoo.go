package yahoo

import (
	"fmt"
	"time"
)

var baseURL = "https://query1.finance.yahoo.com/v7/finance/download/BTC-USD\""

func getPriceIntervalURI(r int) string {
	start := time.Now()
	stop := start.AddDate(0, 0, r*-1)
	intervalInfo := fmt.Sprintf("?period1\\=%s\\&period2\\=%s", start, stop)
	boilerplate := "\\&interval\\=1d\\&events\\=history\\&includeAdjustedClose\\=true"
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

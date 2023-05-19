package binance

import (
	"fmt"
	"net/http"
)

var baseURL = "https://api.binance.us"

func GetAverages(longWin int, shortWin int) (float64, float64, error) {
	return 0, 0, nil
}

func Ping() error {
	endpoint := "/api/v3/ping"
	response, err := http.Get(baseURL + endpoint)
	if err != nil {
		return fmt.Errorf("binance.Ping(): %s", err)
	}
	response.Body.Close()
	if response.StatusCode != 200 {
		return fmt.Errorf("binance.Ping(): %s", response.Status)
	}
	return nil
}

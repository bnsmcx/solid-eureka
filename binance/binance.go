package binance

import (
	"fmt"
	"net/http"
)

var baseURL = "https://api.binance.us"
var apiKey = "<your_api_key>"
var secretKey = "<your_secret_key>"

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

package binance

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
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

func GetPrice() (float64, error) {
	endpoint := "/api/v3/ticker/price?symbol=BTCUSD"

	type quote struct {
		Symbol string `json:"symbol"`
		Price  string `json:"price"`
	}

	response, err := http.Get(baseURL + endpoint)
	if err != nil || response.StatusCode != 200 {
		return 0, fmt.Errorf("binance.GetPrice(): %s, %s", err, response.Status)
	}
	defer response.Body.Close()

	var q quote
	err = json.NewDecoder(response.Body).Decode(&q)
	if err != nil {
		return 0, fmt.Errorf("binance.GetPrice(): %s", err)
	}

	price, err := strconv.ParseFloat(q.Price, 64)
	if err != nil {
		return 0, fmt.Errorf("binance.GetPrice(): %s", err)
	}
	return price, nil
}

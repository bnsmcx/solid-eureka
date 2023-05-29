package coincap

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"
)

type PriceData struct {
	Data []struct {
		PriceUsd string    `json:"priceUsd"`
		Time     int64     `json:"time"`
		Date     time.Time `json:"date"`
	} `json:"data"`
	Timestamp int64 `json:"timestamp"`
}

func GetDataForRange(start, end int64) ([]float64, error) {
	baseURL := "https://api.coincap.io/v2/assets/bitcoin/history?interval=h1"
	intervalSettings := fmt.Sprintf("&start=%d&end=%d", start, end)

	resp, err := http.Get(baseURL + intervalSettings)
	if err != nil || resp.StatusCode != 200 {
		return nil, fmt.Errorf("getting data from coincap.io: status: %s, error: %s",
			resp.Status, err)
	}
	defer resp.Body.Close()

	var data PriceData
	err = json.NewDecoder(resp.Body).Decode(&data)
	if err != nil {
		return nil, fmt.Errorf("decoding coincap.io response: %s", err)
	}

	var prices []float64
	var lastTime int64
	for _, item := range data.Data {
		if item.Time > lastTime {
			lastTime = item.Time
			price, err := strconv.ParseFloat(item.PriceUsd, 64)
			if err != nil {
				return nil, fmt.Errorf("expected float, got : %s", item.PriceUsd)
			}
			prices = append(prices, price)
		} else {
			return nil, errors.New("data out of order")
		}
	}
	return prices, nil
}

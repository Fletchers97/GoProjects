package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type PriceResponse struct {
	Symbol string `json:"symbol"`
	Price  string `json:"price"`
}

func fetchPrice(symbol string, stream chan string) {
	url := fmt.Sprintf("https://api.binance.com/api/v3/ticker/price?symbol=%s", symbol)

	for {
		resp, err := http.Get(url)
		if err != nil {
			stream <- fmt.Sprintf("ERROR [%s]: %v", symbol, err)
			continue
		}

		var result PriceResponse
		json.NewDecoder(resp.Body).Decode(&result)
		resp.Body.Close()

		stream <- fmt.Sprintf(" %s: $%s", result.Symbol, result.Price)

		time.Sleep(2 * time.Second)
	}
}

func main() {

	dataChannel := make(chan string)

	go fetchPrice("BTCUSDT", dataChannel)
	go fetchPrice("ETHUSDT", dataChannel)
	go fetchPrice("SOLUSDT", dataChannel)

	for {
		message := <-dataChannel
		currentTime := time.Now().Format("15:04:05")
		fmt.Printf("[%s] %s\n", currentTime, message)
	}
}

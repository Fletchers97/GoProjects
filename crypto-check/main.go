package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"
)

type PriceResponse struct {
	Symbol string `json:"symbol"`
	Price  string `json:"price"`
}

func fetchPrice(symbol string, stream chan string) {
	url := fmt.Sprintf("https://api.binance.com/api/v3/ticker/price?symbol=%s", symbol)
	var lastPrice float64

	for {
		resp, err := http.Get(url)
		if err != nil {
			stream <- fmt.Sprintf("ERROR [%s]: %v", symbol, err)
			time.Sleep(5 * time.Second)
			continue
		}

		var result PriceResponse
		json.NewDecoder(resp.Body).Decode(&result)
		resp.Body.Close()

		currentPrice, err := strconv.ParseFloat(result.Price, 64)
		if err != nil {
			stream <- fmt.Sprintf("[%s] Error parsing price for %s: %v", time.Now().Format("15:04:05"), result.Symbol, err)
			time.Sleep(5 * time.Second)
			continue
		}

		var status string
		if lastPrice != 0 {
			diff := currentPrice - lastPrice
			if currentPrice > lastPrice {
				status = fmt.Sprintf("UP by $%.2f", diff)
			} else if currentPrice < lastPrice {
				status = fmt.Sprintf("DOWN by $%.2f", -diff)
			} else {
				status = "UNCHANGED"
			}
		} else {
			status = "INITIAL PRICE"
		}
		stream <- fmt.Sprintf("[%s] %-9s | Price: $%10.2f | %s",
			time.Now().Format("15:04:05"),
			result.Symbol,
			currentPrice,
			status)

		lastPrice = currentPrice
		time.Sleep(3 * time.Second)
	}
}

func main() {
	fmt.Println("Crypto Monitor: Goroutines version with Text Alerts")

	dataChannel := make(chan string)

	symbols := []string{"BTCUSDT", "ETHUSDT", "SOLUSDT"}

	for _, s := range symbols {
		go fetchPrice(s, dataChannel)
	}

	for massage := range dataChannel {
		fmt.Println(massage)
	}
}

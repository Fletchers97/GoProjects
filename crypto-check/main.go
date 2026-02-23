package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
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
			log.Printf("[ERROR] [%s] Network issue: %v", symbol, err)
			time.Sleep(5 * time.Second)
			continue
		}

		var result PriceResponse
		err = json.NewDecoder(resp.Body).Decode(&result)
		resp.Body.Close()

		if err != nil {
			log.Printf("[WARNING] [%s] Could not decode JSON: %v", symbol, err)
			continue
		}

		currentPrice, err := strconv.ParseFloat(result.Price, 64)
		if err != nil {
			log.Printf("[ERROR] [%s] Price conversion failed: %v", symbol, err)
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
		msg := fmt.Sprintf("%-9s | $%10.2f | %s", result.Symbol, currentPrice, status)
		log.Printf("[INFO] %s", msg) // Write to log with timestamp
		stream <- msg

		lastPrice = currentPrice
		time.Sleep(3 * time.Second)
	}
}

func main() {
	// Create or open the log file for appending
	// os.O_APPEND — add information to the end, os.O_CREATE — create if it doesn't exist, os.O_WRONLY — write only
	file, err := os.OpenFile("app.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println("Could not open log file:", err)
		return
	}
	defer file.Close()

	// Setting up the logger to write to both a file and  console
	log.SetOutput(file)
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)

	fmt.Println("Monitor started. Logs are being saved to app.log")

	dataChannel := make(chan string)

	symbols := []string{"BTCUSDT", "ETHUSDT", "SOLUSDT"}

	for _, s := range symbols {
		go fetchPrice(s, dataChannel)
	}

	for message := range dataChannel {
		fmt.Printf("[%s] %s\n", time.Now().Format("15:04:05"), message)
	}
}

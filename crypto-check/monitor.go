package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"sync"
	"time"
)

func fetchPrice(ctx context.Context, wg *sync.WaitGroup, db *sql.DB, apiUrl string, symbol string, interval int, alertThreshold float64, stream chan string) {
	defer wg.Done() // Ensure we signal when this goroutine is done
	url := apiUrl + symbol
	var lastPrice float64

	for {
		select {
		case <-ctx.Done():
			log.Printf("[INFO] [%s] Stopping price fetcher", symbol)
			return
		default:

			log.Printf("[DEBUG] [%s] Sending request to: %s", symbol, url)

			resp, err := http.Get(url)
			if err != nil {
				log.Printf("[ERROR] [%s] Connection error: %v", symbol, err)
				time.Sleep(time.Duration(interval) * time.Second)
				continue
			}

			var result PriceResponse
			err = json.NewDecoder(resp.Body).Decode(&result)
			resp.Body.Close()

			if err != nil {
				log.Printf("[ERROR] [%s] JSON Decode error: %v", symbol, err)
				time.Sleep(time.Duration(interval) * time.Second)
				continue
			}

			currentPrice, err := strconv.ParseFloat(result.Price, 64)
			if err != nil {
				log.Printf("[ERROR] [%s] Price conversion error ('%s'): %v", symbol, result.Price, err)
				time.Sleep(time.Duration(interval) * time.Second)
				continue
			}

			// Save price to database
			_, err = db.Exec("INSERT INTO price_history (symbol, price, timestamp) VALUES(?, ?, ?)",
				symbol, currentPrice, time.Now())
			if err != nil {
				log.Printf("[ERROR] [%s] Database insert error: %v", symbol, err)
			}

			status := "INITIAL"
			if lastPrice != 0 {
				diff := currentPrice - lastPrice
				absDiff := diff
				if absDiff < 0 {
					absDiff = -absDiff
				}
				if absDiff >= alertThreshold {
					log.Printf("[WARNING] [%s] VOLATILITY ALERT:Price changed by $%.2f (Threshold: $%.2f)", symbol, diff, alertThreshold)
				}
				if currentPrice > lastPrice {
					status = fmt.Sprintf("UP (+$%.2f)", diff)
				} else if currentPrice < lastPrice {
					status = fmt.Sprintf("DOWN (-$%.2f)", -diff)
				} else {
					status = "STABLE"
				}
			}

			msg := fmt.Sprintf("%-9s | $%10.2f | %s", result.Symbol, currentPrice, status)
			log.Printf("[INFO] %s", msg)
			stream <- msg

			lastPrice = currentPrice
			time.Sleep(time.Duration(interval) * time.Second)
		}
	}
}

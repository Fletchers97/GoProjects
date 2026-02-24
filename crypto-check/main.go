package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	_ "github.com/glebarez/go-sqlite"
)

type Config struct {
	ApiUrl         string   `json:"api_url"`
	Symbols        []string `json:"symbols"`
	UpdateInterval int      `json:"update_interval"`
	AlertThreshold float64  `json:"alert_threshold"`
}

type PriceResponse struct {
	Symbol string `json:"symbol"`
	Price  string `json:"price"`
}

func loadConfig(fileName string) (Config, error) {
	var config Config
	configFile, err := os.Open(fileName)
	if err != nil {
		return config, fmt.Errorf("failed to open config: %w", err)
	}
	defer configFile.Close()

	err = json.NewDecoder(configFile).Decode(&config)
	if err != nil {
		return config, fmt.Errorf("failed to decode config: %w", err)
	}
	return config, nil
}

func initDB(filepath string) (*sql.DB, error) {
	db, err := sql.Open("sqlite", filepath)
	if err != nil {
		return nil, err
	}

	// Создаем таблицу, если её еще нет
	query := `
	CREATE TABLE IF NOT EXISTS price_history (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		symbol TEXT,
		price REAL,
		timestamp DATETIME
	);`

	_, err = db.Exec(query)
	if err != nil {
		return nil, err
	}

	return db, nil
}

func fetchPrice(db *sql.DB, apiUrl string, symbol string, interval int, alertThreshold float64, stream chan string) {
	url := apiUrl + symbol
	var lastPrice float64

	for {
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

func main() {
	// Setting up logs BEFORE loading the config so that config errors are also logged
	file, err := os.OpenFile("app.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Printf("Fatal: could not open log file: %v\n", err)
		return
	}
	defer file.Close()
	log.SetOutput(file)
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)

	db, err := initDB("crypto.db")
	if err != nil {
		log.Fatalf("[FATAL] Database initialization failed: %v", err)
		return
	}
	defer db.Close()

	config, err := loadConfig("config.json")
	if err != nil {
		fmt.Printf("[FATAL] %v\n", err)
		log.Fatalf("[FATAL] Configuration failed: %v", err)
		return
	}

	fmt.Printf("Monitor started. Symbols: %v. Interval: %ds\n", config.Symbols, config.UpdateInterval)

	dataChannel := make(chan string)

	for _, s := range config.Symbols {
		go fetchPrice(db, config.ApiUrl, s, config.UpdateInterval, config.AlertThreshold, dataChannel)
	}

	for message := range dataChannel {
		fmt.Printf("[%s] %s\n", time.Now().Format("15:04:05"), message)
	}
}

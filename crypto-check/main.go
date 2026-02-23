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

type Config struct {
	ApiUrl         string   `json:"api_url"`
	Symbols        []string `json:"symbols"`
	UpdateInterval int      `json:"update_interval"`
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

func fetchPrice(apiUrl string, symbol string, interval int, stream chan string) {
	url := apiUrl + symbol
	var lastPrice float64

	for {
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

		status := "INITIAL"
		if lastPrice != 0 {
			diff := currentPrice - lastPrice
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

	config, err := loadConfig("config.json")
	if err != nil {
		log.Printf("[FATAL] %v", err)
		fmt.Printf("[FATAL] %v\n", err)
		return
	}

	fmt.Printf("Monitor started. Symbols: %v. Interval: %ds\n", config.Symbols, config.UpdateInterval)

	dataChannel := make(chan string)

	for _, s := range config.Symbols {
		go fetchPrice(config.ApiUrl, s, config.UpdateInterval, dataChannel)
	}

	for message := range dataChannel {
		fmt.Printf("[%s] %s\n", time.Now().Format("15:04:05"), message)
	}
}

package main

import (
	"encoding/json"
	"fmt"
	"math"
	"os"
)

func loadConfig(fileName string) (Config, error) {
	var config Config
	configFile, err := os.Open(fileName)
	if err != nil {
		return config, err
	}
	defer configFile.Close()
	err = json.NewDecoder(configFile).Decode(&config)
	return config, err
}

func ValidateConfig(symbols []string, interval int) bool {
	if len(symbols) == 0 {
		return false // List of symbols cannot be empty
	}
	if interval <= 0 {
		return false // Interval must be a positive integer
	}
	return true
}

func FormatDisplayPrice(price float64) string {
	if price < 1.0 {
		return fmt.Sprintf("%.8f", price) // For small prices like doge, shiba, etc.
	}
	return fmt.Sprintf("%.2f", price) // For larger prices like BTC, ETH, etc.
}

func CalculatePercentageDiff(current, average float64) float64 {
	if average == 0 {
		return 0
	}
	diff := ((current - average) / average) * 100
	return math.Round(diff*100) / 100
}

func GetPriceStatus(diff float64) string {
	if diff >= 5.0 {
		return "ROCKET"
	} else if diff <= -5.0 {
		return "CRASH"
	}
	return "STABLE"
}

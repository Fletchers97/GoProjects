package main

import (
	"encoding/json"
	"testing"
)

func TestValidateConfig(t *testing.T) {
	tests := []struct {
		name     string
		symbols  []string
		interval int
		want     bool
	}{
		{"Valid config", []string{"BTCUSDT", "ETHUSDT"}, 5, true},
		{"Empty symbols", []string{}, 5, false},
		{"Negative interval", []string{"BTCUSDT"}, -1, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ValidateConfig(tt.symbols, tt.interval); got != tt.want {
				t.Errorf("ValidateConfig() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFormatDisplayPrice(t *testing.T) {
	tests := []struct {
		price float64
		want  string
	}{
		{65000.50, "65000.50"},
		{0.00001234, "0.00001234"},
		{1.0, "1.00"},
	}

	for _, tt := range tests {
		got := FormatDisplayPrice(tt.price)
		if got != tt.want {
			t.Errorf("FormatDisplayPrice(%f) = %s; want %s", tt.price, got, tt.want)
		}
	}
}

func TestCalculatePercentageDiff(t *testing.T) {
	tests := []struct {
		name    string
		current float64
		average float64
		want    float64
	}{
		{"Price went up", 110.0, 100.0, 10.0},
		{"Price went down", 90.0, 100.0, -10.0},
		{"No change", 100.0, 100.0, 0.0},
		{"Zero average", 150.0, 0.0, 0.0}, // Check division by zero case
		{"Small change", 100.05, 100.0, 0.05},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := CalculatePercentageDiff(tt.current, tt.average)
			if got != tt.want {
				t.Errorf("%s: got %f, want %f", tt.name, got, tt.want)
			}
		})
	}
}

func TestGetPriceStatus(t *testing.T) {
	tests := []struct {
		name string
		diff float64
		want string
	}{
		{"High growth", 7.5, "ROCKET"},
		{"Small growth", 2.1, "STABLE"},
		{"Big drop", -10.2, "CRASH"},
		{"Small drop", -1.5, "STABLE"},
		{"Exactly five", 5.0, "ROCKET"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetPriceStatus(tt.diff); got != tt.want {
				t.Errorf("GetPriceStatus(%f) = %v, want %v", tt.diff, got, tt.want)
			}
		})
	}
}

func TestSymbolJSONParsing(t *testing.T) {
	// Imitating the JSON response from Binance API for a symbol price
	jsonData := `{"symbol":"BTCUSDT","price":"65000.00"}`

	var res struct {
		Symbol string `json:"symbol"`
		Price  string `json:"price"`
	}

	err := json.Unmarshal([]byte(jsonData), &res)

	if err != nil {
		t.Fatalf("JSON Unmarshal failed: %v", err)
	}

	if res.Symbol != "BTCUSDT" {
		t.Errorf("Expected BTCUSDT, got %s", res.Symbol)
	}

	if res.Price != "65000.00" {
		t.Errorf("Expected 65000.00, got %s", res.Price)
	}
}

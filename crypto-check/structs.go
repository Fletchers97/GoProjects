package main

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

type CoinStats struct {
	Symbol   string  `json:"symbol"`
	Price    float64 `json:"current_price"`
	AvgPrice float64 `json:"avg_price_1h"`
}

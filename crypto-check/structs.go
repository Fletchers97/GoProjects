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

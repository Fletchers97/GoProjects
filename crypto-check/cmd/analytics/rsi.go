package main

func CalculateRSI(prices []float64) float64 {
	if len(prices) < 2 {
		return 50.0 // Недостаточно данных
	}

	var gains, losses float64
	for i := 1; i < len(prices); i++ {
		change := prices[i] - prices[i-1]
		if change > 0 {
			gains += change
		} else {
			losses -= change
		}
	}

	if losses == 0 {
		return 100.0
	}

	rs := gains / losses
	return 100.0 - (100.0 / (1 + rs))
}

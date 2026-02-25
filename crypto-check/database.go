package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/glebarez/go-sqlite"
)

func initDB(filepath string) (*sql.DB, error) {
	db, err := sql.Open("sqlite", filepath)
	if err != nil {
		return nil, err
	}

	query := `
	CREATE TABLE IF NOT EXISTS price_history (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		symbol TEXT,
		price REAL,
		timestamp DATETIME
	);`

	_, err = db.Exec(query)
	return db, err
}

func getAveragePrice(db *sql.DB, symbol string, minutes int) (float64, error) {
	var avgPrice sql.NullFloat64 // Use NullFloat64 to handle cases where there might be no data in the database for the given period

	// SQL query: calculate the average (AVG) for the period from current time backwards
	query := `
		SELECT AVG(price) 
		FROM price_history 
		WHERE symbol = ? AND timestamp > datetime('now', ?)`

	interval := fmt.Sprintf("-%d minutes", minutes)

	err := db.QueryRow(query, symbol, interval).Scan(&avgPrice)
	if err != nil {
		return 0, err
	}

	if !avgPrice.Valid {
		return 0, nil // No data available for the given symbol and time period, return 0 as average price
	}

	return avgPrice.Float64, nil
}

func getLatestStats(db *sql.DB) ([]CoinStats, error) {

	query := `
        SELECT 
            t1.symbol, 
            t1.price,
            (SELECT AVG(price) FROM price_history 
             WHERE symbol = t1.symbol AND timestamp > datetime('now', '-100 hours')) as avg_price
        FROM price_history t1
        WHERE t1.id IN (SELECT MAX(id) FROM price_history GROUP BY symbol)`

	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var stats []CoinStats
	for rows.Next() {
		var s CoinStats
		var avg sql.NullFloat64
		if err := rows.Scan(&s.Symbol, &s.Price, &avg); err != nil {
			log.Printf("[ERROR] Scan error: %v", err)
			continue
		}
		if avg.Valid {
			s.AvgPrice = avg.Float64
		}
		stats = append(stats, s)
	}
	return stats, nil
}

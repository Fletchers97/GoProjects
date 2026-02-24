package main

import (
	"database/sql"

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

package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"text/template"
)

// StartServer runs the web server on the specified port and sets up the API endpoint for stats
func StartServer(db *sql.DB, port string) {
	// Register the handler function for the /stats endpoint
	http.HandleFunc("/api/stats", getStatsHandler(db))
	http.HandleFunc("/", getIndexHandler(db))

	log.Printf("[INFO] Web server starting on http://localhost%s/stats", port)

	if err := http.ListenAndServe(port, nil); err != nil {
		log.Fatalf("[FATAL] Server failed to start: %v", err)
	}
}

func getIndexHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Collect the latest stats from the database
		stats, err := getLatestStats(db)
		if err != nil {
			http.Error(w, "Database error", http.StatusInternalServerError)
			return
		}

		// Parse and execute the HTML template, passing the stats data
		tmpl, err := template.ParseFiles("templates/index.html")
		if err != nil {
			log.Printf("[ERROR] Template error: %v", err)
			http.Error(w, "Template not found", http.StatusInternalServerError)
			return
		}

		tmpl.Execute(w, stats)
	}
}

func getStatsHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Json response header
		w.Header().Set("Content-Type", "application/json")

		// Collect the latest stats from the database
		stats, err := getLatestStats(db)
		if err != nil {
			log.Printf("[ERROR] API Stats error: %v", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		// Encode the stats as JSON and send the response
		encoder := json.NewEncoder(w)
		encoder.SetIndent("", "  ")
		if err := encoder.Encode(stats); err != nil {
			log.Printf("[ERROR] JSON encoding error: %v", err)
		}
	}
}

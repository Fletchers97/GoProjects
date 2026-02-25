package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	_ "github.com/glebarez/go-sqlite"
)

func main() {
	// Create a context that can be cancelled on interrupt signal
	ctx, cancel := context.WithCancel(context.Background())

	// Listen for interrupt signals to gracefully shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

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

	go StartServer(db, ":8080")

	dataChannel := make(chan string)
	var wg sync.WaitGroup
	for _, s := range config.Symbols {
		wg.Add(1) // Increment WaitGroup counter for each goroutine
		go fetchPrice(ctx, &wg, db, config.ApiUrl, s, config.UpdateInterval, config.AlertThreshold, dataChannel)
	}

	go func() {
		for message := range dataChannel {
			fmt.Printf("[%s] %s\n", time.Now().Format("15:04:05"), message)
		}
	}()

	sig := <-sigChan
	log.Printf("[INFO] Received signal: %v. Shutting down...", sig)
	fmt.Printf("[INFO] Received signal: %v. Shutting down...\n", sig)
	cancel()
	wg.Wait() // Wait for all fetchPrice goroutines to finish
	log.Println("[INFO] Shutdown complete.")
	close(dataChannel)

	fmt.Println("Program terminated gracefully. All data was saved.")
}

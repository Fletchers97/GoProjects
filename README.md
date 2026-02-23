# Go Learning & FinTech Pet Projects

This repository contains my practice projects in Golang, focused on building robust, concurrent financial data processing systems.

## Projects

### Crypto-Check (Binance Monitor)
A real-time cryptocurrency price monitor designed with standard practices.

**Key Features:**
* **High Concurrency:** Uses goroutines and channels to track multiple symbols simultaneously.
* **Smart Logging:** Implements a full logging hierarchy (DEBUG, INFO, WARNING, ERROR, FATAL) with file output.
* **Volatility Alerts:** Configurable price change thresholds to detect market swings.
* **External Configuration:** Fully driven by a `config.json` file (no hardcoded settings).
* **Error Handling:** Robust protection against network issues, JSON decoding errors, and configuration failures.

## Tech Stack
* **Language:** Golang 
* **API:** Binance REST API
* **Architecture:** Data-driven design via JSON configuration
* **Ops:** Professional Git workflow and repository hygiene

## Roadmap
- [x] Concurrent price fetching
- [x] JSON Configuration system
- [x] Multi-level logging & Alerts
- [ ] SQLite database integration (In progress)
- [ ] Telegram Bot notifications
- [ ] Docker containerization
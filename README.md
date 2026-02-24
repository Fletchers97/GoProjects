Go Learning & FinTech Pet Projects
This repository contains my practice projects in Golang, focused on building robust, concurrent financial data processing systems.

Projects
Crypto-Check (Binance Monitor)
A real time cryptocurrency price monitor designed with standard practices. This project evolved from a simple scraper into a modular microservice with persistent storage and analytics.

Key Features:

High Concurrency: Uses goroutines and channels to track multiple symbols simultaneously.

Modular Architecture: Cleanly separated into monitor, database, models, and utils for better maintainability.

Persistent Storage: Integrated SQLite to keep track of price history.

Real-time Analytics: Calculates Moving Averages and price deviations directly via SQL queries.

Smart Logging: Implements a full logging hierarchy (DEBUG, INFO, WARNING, ERROR, FATAL) with file output.

Volatility Alerts: Configurable price change thresholds and trend deviation alerts to detect market swings.

External Configuration: Fully driven by a config.json file (no hardcoded settings).

Error Handling: Robust protection against network issues, JSON decoding errors, and configuration failures.

Tech Stack:

Language: Golang (Concurrency, Context, Sync)

Database: SQLite (SQL, Time-series data)

API: Binance REST API

Architecture: Modular Data-driven design

Roadmap
[x] Concurrent price fetching

[x] JSON Configuration system

[x] Multi-level logging & Alerts

[x] SQLite database integration & Analytics

[ ] REST API for statistics (In progress)

[ ] Telegram Bot notifications

[ ] Docker containerization
# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

Go-based stock monitoring system for Chinese A-stock market. Fetches real-time quotes from Sina Finance, evaluates configurable rules, and sends alerts via multiple notification channels (Feishu, Server酱, DingTalk). Includes an embedded web management UI.

## Build & Run Commands

```bash
# Build
go build -o monitor ./cmd/monitor

# Run (starts web UI on :8080 + monitoring every 5 min)
./monitor

# Run with custom port
./monitor -addr :9090

# Single check mode (no web server, exits after one check)
./monitor -once

# Custom data file
./monitor -data /path/to/config.json

# Docker
docker-compose up -d
```

## Testing & Quality

No test suite or linter configuration exists yet. Standard Go tooling applies:

```bash
go vet ./...
go build ./...
```

## Architecture

The application is a single Go binary (`cmd/monitor/main.go`) using only the standard library HTTP server. Go module name: `stock-monitor`, Go 1.24.0.

### Core Components (all under `internal/`)

- **api/** — REST API server (`server.go`) and embedded HTML UI (`html.go`). Routes: `/api/stocks`, `/api/rules`, `/api/rule-types`, `/api/notifiers`, `/` (web UI).
- **monitor/** — Main loop: runs every 5 minutes, fetches quotes, evaluates rules, sends alerts. Handles graceful shutdown via SIGINT/SIGTERM.
- **datasource/** — `DataSource` interface with Sina Finance implementation. Fetches real-time quotes (`hq.sinajs.cn`) and K-line data. Handles GBK→UTF-8 encoding. Stock codes auto-prefixed with exchange (sh/sz based on code prefix).
- **rule/** — Plugin-based rule engine using registry pattern. `rule.go` defines the `Rule` interface, `registry.go` manages factory functions, `engine.go` evaluates rules. Concrete rules live in `rule/rules/`.
- **notifier/** — `Notifier` interface with implementations for Feishu (card messages), Server酱 (WeChat). `manager.go` dispatches alerts to all enabled notifiers concurrently.
- **storage/** — Thread-safe JSON file persistence (`data/config.json`) using `sync.RWMutex`. Stores stocks, rules, and notifier configs.
- **indicator/** — Technical indicators. Currently only Moving Average (`ma.go`).
- **model/** — Data models: `Stock`, `KLine`, `Alert` (with levels: info/warning/critical).

### Adding a New Rule

1. Create a file in `internal/rule/rules/`
2. Implement the `Rule` interface: `Name()`, `Description()`, `Validate()`, `Evaluate()`
3. For K-line based rules, also implement `KLineRule` interface
4. Register via `init()`: `rule.GlobalRegistry.Register("type_name", factoryFunc, "Display Name")`

### Adding a New Notifier

1. Create a file in `internal/notifier/`
2. Implement the `Notifier` interface: `Name()`, `Send(ctx, alert)`
3. Add config fields to `storage.NotifierConfig` in `storage/model.go`
4. Wire it up in `notifier/manager.go`

### K-Line Type Mapping (Sina API scale values)

5min→5, 15min→15, 30min→30, 60min→60, daily→240, weekly→1200, monthly→7200

## Key Dependencies

- `github.com/google/uuid` — Rule ID generation
- `golang.org/x/text` — GBK to UTF-8 conversion for Sina API responses
- `gopkg.in/yaml.v3` — YAML config parsing

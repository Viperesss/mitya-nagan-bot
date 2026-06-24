# Telegram Read Later Bot

A Telegram bot written in Go that allows users to save web pages and read them later. Users can send links to the bot, store them in a personal collection, and retrieve a random unread page at any time.

## Features

* Save links sent to the bot
* Store pages separately for each user
* Get a random saved page
* Remove pages after reading
* SQLite-based storage
* Duplicate link detection
* Concurrent event processing with a worker pool
* Automatic retry mechanism with exponential backoff

## Commands

| Command  | Description             |
| -------- | ----------------------- |
| `/start` | Show help message       |
| `/help`  | Show help message       |
| `/random`   | Get a random saved page |

## Technologies

* Go
* Telegram Bot API
* SQLite
* go-sqlite3
* Goroutines & Channels
* Worker Pool Pattern

## Project Structure

```text
clients/         - Telegram API clients
consumer/        - Event processing
storage/         - Storage interfaces
storage/files/   - file-based storage
storage/sqlite/  - SQLite implementation
lib/             - Shared utilities
cmd/bot/         - Application entry point
data/sqlite/     - Database files
```

## Run

```bash
go run ./cmd/bot -token '<TOKEN>'
```

## Future Improvements

* Unit and integration tests
* Docker support
* LLM integration
* Web interface for saved pages
* Graceful shutdown with context cancellation
# Telegram Read Later Bot

A Telegram bot written in Go that combines a **read-it-later** service with an AI-powered chat assistant.

Users can save web pages for later reading, retrieve a random unread page, or simply chat with the bot through an LLM powered by the OpenRouter API.

## Features

* Save links sent to the bot
* Store pages separately for each user
* Get a random saved page
* Remove pages after reading
* SQLite-based storage
* Duplicate link detection
* AI-powered conversations via OpenRouter
* Concurrent event processing with a worker pool
* Automatic retry mechanism with exponential backoff

## Commands

| Command  | Description             |
| -------- | ----------------------- |
| `/start` | Show help message       |
| `/help`  | Show help message       |
| `/random`   | Get a random saved page |

Any non-command text message is forwarded to the AI assistant.

## Technologies

* Go
* Telegram Bot API
* OpenRouter API
* SQLite
* go-sqlite3
* godotenv
* Goroutines & Channels
* Worker Pool Pattern

## Project Structure

```text
clients/
    events/       - Event processing
    llm/          - OpenRouter client
    telegram/     - Telegram Bot API client

consumer/         - Event consumer with worker pool
storage/
    files/        - File-based storage
    sqlite/       - SQLite implementation

lib/              - Shared utilities
cmd/bot/          - Application entry point
data/sqlite/      - Database files
```

## Configuration

Create a `.env` file in the project root:

```env
OPENROUTER_API_KEY=<your_api_key>
```

## Run

```bash
go run ./cmd/bot -token "<TOKEN>"
```

## Future Improvements

* Unit and integration tests
* Docker support
* Web interface for saved pages
* Graceful shutdown with context cancellation
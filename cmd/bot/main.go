package main

import (
	"context"
	"flag"
	"log"
	"os"
	tgProcessor "the-mitya-nagan-bot/clients/events/telegram"
	"the-mitya-nagan-bot/clients/llm"
	tgClient "the-mitya-nagan-bot/clients/telegram"
	event_consumer "the-mitya-nagan-bot/consumer/event-consumer"
	"the-mitya-nagan-bot/storage/sqlite"

	"github.com/joho/godotenv"
)

const (
	sqliteStoragePath = "data/sqlite/storage.db"
	batchSize         = 100
	hostAPI           = "api.telegram.org"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("failed Load .env: ", err)
	}

	apiKey := os.Getenv("OPENROUTER_API_KEY")
	if apiKey == "" {
		log.Fatal("api key cannot was empty")
	}

	// s := files.New(storagePath) // файловое хранение
	s, err := sqlite.New(sqliteStoragePath)
	if err != nil {
		log.Fatal("cannot connect to storage: ", err)
	}

	// TODO context.WithTimeOut 5 sec
	if err := s.Init(context.TODO()); err != nil {
		log.Fatal("cannot init storage: ", err)
	}

	tgClient, err := tgClient.New(hostAPI, mustToken())
	if err != nil {
		log.Fatal("cannot init tgClient: ", err)
	}

	llmClient := llm.New(apiKey)

	eventsProcessor := tgProcessor.New(tgClient, llmClient, s)

	log.Print("service started")

	consumer := event_consumer.New(eventsProcessor, eventsProcessor, batchSize)

	if err := consumer.Start(); err != nil {
		// произойдет только при аварийной остановке
		log.Fatal("service is stopped", err)
	}
	// consumer.Start(fetcher, processor) // потребитель
}

func mustToken() string {
	token := flag.String(
		"token",
		"",
		"token for access to telegram bot")

	flag.Parse()

	// Во время Parse() пакет flag записывает значение прямо по адресу переменной, поэтому успользуется pointer
	if *token == "" {
		log.Fatal("-token - is not specified")
	}

	return *token
}

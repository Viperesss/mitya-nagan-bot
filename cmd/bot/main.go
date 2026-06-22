package main

import (
	"context"
	"flag"
	"log"
	tgProcessor "the-mitya-nagan-bot/clients/events/telegram"
	tgClient "the-mitya-nagan-bot/clients/telegram"
	event_consumer "the-mitya-nagan-bot/consumer/event-consumer"
	"the-mitya-nagan-bot/storage/sqlite"
)

const (
	sqliteStoragePath = "data/sqlite/storage.db"
	batchSize         = 100
	hostAPI           = "api.telegram.org"
)

func main() {
	// s := files.New(storagePath) // файловое хранение
	s, err := sqlite.New(sqliteStoragePath)
	if err != nil {
		log.Fatal("cannot connect to storage: ", err)
	}

	// TODO context.WithTimeOut 5 sec
	if err := s.Init(context.TODO()); err != nil {
		log.Fatal("cannot init storage: ", err)
	}

	tgClient := tgClient.New(hostAPI, mustToken())
	eventsProcessor := tgProcessor.New(tgClient, s)

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

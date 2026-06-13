package main

import (
	"flag"
	"log"
	tgProcessor "the-mitya-nagan-bot/clients/events/telegram"
	tgClient "the-mitya-nagan-bot/clients/telegram"
	event_consumer "the-mitya-nagan-bot/consumer/event-consumer"
	"the-mitya-nagan-bot/storage/files"
)

const (
	storagePath = "storage"
	batchSize   = 100
)

var (
	token = flag.String(
		"token",
		"",
		"token for access to telegram bot")

	host = flag.String(
		"host",
		"",
		"the telegram API host")
)

// api.telegram.org
func main() {
	flag.Parse()

	tgClient := tgClient.New(mustHost(), mustToken())
	eventsProcessor := tgProcessor.New(tgClient, files.New(storagePath))

	log.Print("service started")

	consumer := event_consumer.New(eventsProcessor, eventsProcessor, batchSize)

	if err := consumer.Start(); err != nil {
		// произойдет только при аварийной остановке
		log.Fatal("service is stopped", err)
	}
	// consumer.Start(fetcher, processor) // потребитель
}

func mustToken() string {
	// Во время Parse() пакет flag записывает значение прямо по адресу переменной, поэтому успользуется pointer
	if *token == "" {
		log.Fatal("-token - is not specified")
	}

	return *token
}

func mustHost() string {
	if *host == "" {
		log.Fatal("-host - is not specified")
	}
	return *host
}

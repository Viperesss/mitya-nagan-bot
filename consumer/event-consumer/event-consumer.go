package event_consumer

import (
	"log"
	"the-mitya-nagan-bot/clients/events"
	"time"
)

type Consumer struct {
	fetcher   events.Fetcher
	processor events.Processor
	batchSize int // сколько событий будем обрабатывать за раз
}

func New(fetcher events.Fetcher, processor events.Processor, batchSize int) *Consumer {
	return &Consumer{
		fetcher:   fetcher,
		processor: processor,
		batchSize: batchSize,
	}
}

func (c *Consumer) Start() error {
	// вечный цикл который ждет новые события и обрабатывает их
	for {
		gotEvents, err := c.fetcher.Fetch(c.batchSize)
		if err != nil {
			log.Printf("[ERR] consumer: %s", err.Error())
			continue
		}

		if len(gotEvents) == 0 {
			time.Sleep(500 * time.Millisecond)
			// добавить механизм retry'а
			continue
		}

		if err := c.handleEvents(gotEvents); err != nil {
			log.Print(err)
			continue
		}
	}
}

// добавить ассинхронность
// https://youtu.be/HTjNyoQumJk?si=zEBxiTxU2mqQE6Ml&t=500
func (c *Consumer) handleEvents(events []events.Event) error {
	for _, event := range events {
		log.Printf("got new event: %s", event.Text)

		if err := c.processor.Process(event); err != nil {
			log.Printf("can't handle event: %s", err.Error())
			// добавить механизм retry'а
			continue
		}
	}
	return nil
}

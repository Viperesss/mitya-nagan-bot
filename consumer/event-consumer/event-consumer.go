// Package event_consumer continuosly fetches and processes incoming events.
package event_consumer

import (
	"fmt"
	"log"
	"sync"
	"the-mitya-nagan-bot/clients/events"
	"time"
)

const maxRetries = 5

// Consumer fetches events and passes them to a processor.
type Consumer struct {
	fetcher   events.Fetcher
	processor events.Processor
	batchSize int // сколько событий будем обрабатывать за раз
}

// New creates a new event consumer.
func New(fetcher events.Fetcher, processor events.Processor, batchSize int) *Consumer {
	return &Consumer{
		fetcher:   fetcher,
		processor: processor,
		batchSize: batchSize,
	}
}

// Start continuosly fetches and processes incoming events.
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

const workersCount = 3

// handleEvents processes a batch of events.
func (c *Consumer) handleEvents(eventsList []events.Event) error {
	var wg sync.WaitGroup

	eventsChan := make(chan events.Event)

	for i := 0; i < workersCount; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()

			for event := range eventsChan {
				log.Printf("got new event: %s", event.Text)

				if err := c.retry(event); err != nil {
					log.Print("Timeout exceeded:", err)
				}
			}
		}()
	}

	for _, event := range eventsList {
		eventsChan <- event
	}
	close(eventsChan)

	wg.Wait()

	return nil
}

func (c *Consumer) retry(event events.Event) error {
	delay := time.Millisecond * 500
	var err error

	for i := 1; i <= maxRetries; i++ {
		err = c.processor.Process(event)
		if err == nil {
			return nil
		}

		log.Printf("can't handle event: %v, attempt: %d/%d", err, i, maxRetries)

		if i < maxRetries {
			time.Sleep(delay)
			delay *= 2
		}
	}

	return fmt.Errorf("event processing failed after %d attempts: %w", maxRetries, err)
}

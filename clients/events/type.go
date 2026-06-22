// Package events defines messanger-agnostic event abstractions.
package events

// Fetcher retrieves events from an external source.
type Fetcher interface {
	Fetch(limit int) ([]Event, error) // offset перенесли внутрь, так как он там сам разберется сколько нужно отправить
}

// Processor processes incoming events.
type Processor interface {
	Process(e Event) error
}

// Type represents an event type.
type Type int

const (
	// Unknown represents an unsupported or unknown event type.
	Unknown Type = iota

	// Message represents a text message event.
	Message
)

// Event contains normalized event independent of the source platform.
type Event struct {
	Type Type
	Text string
	Meta interface{}
}

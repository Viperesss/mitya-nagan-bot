// Абстрактный пакет events, не привязан к конктретному мессенджеру
package events

// податель
type Fetcher interface {
	Fetch(limit int) ([]Event, error) // offset перенесли внутрь, так как он там сам разберется сколько нужно отправить
}

type Processor interface {
	Process(e Event) error
}

type Type int

const (
	Unknown Type = iota
	Message
)

// абстрактное событие
type Event struct {
	Type Type
	Text string
	Meta interface{}
}

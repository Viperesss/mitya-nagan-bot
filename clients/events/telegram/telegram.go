// Package telegram implements Telegram-based event processing.
package telegram

import (
	"errors"
	"the-mitya-nagan-bot/clients/events"
	"the-mitya-nagan-bot/clients/llm"
	"the-mitya-nagan-bot/clients/telegram"
	"the-mitya-nagan-bot/lib/e"
	"the-mitya-nagan-bot/storage"
)

// это просто набор зависимостей и состояния, которые нужны всем методам структуры
// отвечает за получение и обработку сообщений
// Processor fetches events from Telegram and processes them.
type Processor struct {
	tg      *telegram.Client
	llm     *llm.Client
	offset  int // значение offset должно переживать много вызовов метода. Поэтому оно хранится внутри объекта
	storage storage.Storage
}

// Meta contains Telegram-specific event metadata.
type Meta struct {
	ChatID   int
	Username string
}

// Клиент передаем по указателю:
// - не копируем объект Client
// - работаем с одним экземпляром клиента
//
// Возвращаем указатель:
// - не копируем Processor
// - методы смогут изменять его состояние (offset и др.)
// New creates a new Telegram event processor.
func New(client *telegram.Client, llm *llm.Client, storage storage.Storage) *Processor {
	return &Processor{
		tg:      client,
		llm:     llm,
		storage: storage,
	}
}

// Fetch получает сырые данные от Telegram и превращает их в удобный для приложения формат.
// Update - сущность тг, event - общая сущность
// Fetch retrieves updates from Telegram and converts them into application events.
func (p *Processor) Fetch(limit int) ([]events.Event, error) { // []events.Event — это просто список новых сообщений
	// Возьми Telegram-клиент, который хранится внутри Processor, и запроси у него обновления, начиная с текущего offset, в количестве не больше limit
	updates, err := p.tg.Updates(p.offset, limit)
	if err != nil {
		return nil, e.Wrap("processor Fetch fail, cannot get events:", err)
	}

	// Если список обновлений оказался пустым - то сразу заканчиваем работу функции
	if len(updates) == 0 {
		return nil, nil
	}

	// Заранее выделяем память, т.к. знаем сколько будет значений
	res := make([]events.Event, 0, len(updates)) // len, cap

	// перебираем апдейты и преобразуем их в тип ивент
	for _, u := range updates {
		res = append(res, event(u))
	}

	// Обновляем оффсет чтобы в следующий раз получить новую пачку сообщений
	p.offset = updates[len(updates)-1].UpdateId + 1

	return res, nil
}

// Process routes an event to the appropriate handler.
func (p *Processor) Process(event events.Event) error {
	switch event.Type {
	case events.Message:
		return p.processMessage(event)
	default:
		return e.Wrap("cannot process message:", errors.New("unknown event type"))
	}
}

// processMessage handles message events.
func (p *Processor) processMessage(event events.Event) error {
	meta, err := getMeta(event)
	if err != nil {
		return e.Wrap("cannot process message, getMeta:", err)
	}

	if err := p.doCmd(event.Text, meta.ChatID, meta.Username); err != nil {
		return e.Wrap("cannot process message, doCmd:", err)
	}
	return nil
}

// getMeta extracts metadata from an event.
func getMeta(event events.Event) (Meta, error) {
	// Попробуй взять значение из event.Meta и привести его к типу Meta
	res, ok := event.Meta.(Meta)
	if !ok {
		return Meta{}, e.Wrap("cannot get meta:", errors.New("unknown meta type"))
	}
	return res, nil
}

// events converts a Telegram update into an application event.
func event(upd telegram.Update) events.Event {
	updType := fetchType(upd)

	res := events.Event{
		Type: updType,
		Text: fetchText(upd),
	}

	if updType == events.Message {
		res.Meta = Meta{
			ChatID:   upd.Message.Chat.ID,
			Username: upd.Message.From.Username,
		}
	}
	return res
}

// fetchText extracts message text from a Telegram update.
func fetchText(upd telegram.Update) string {
	if upd.Message == nil {
		return ""
	}
	return upd.Message.Text
}

// fetchType determines the application event type.
func fetchType(upd telegram.Update) events.Type {
	if upd.Message == nil {
		return events.Unknown
	}
	return events.Message
}

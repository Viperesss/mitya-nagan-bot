package telegram

import (
	"errors"
	"the-mitya-nagan-bot/clients/events"
	"the-mitya-nagan-bot/clients/telegram"
	"the-mitya-nagan-bot/lib/e"
	"the-mitya-nagan-bot/storage"
)

// это просто набор зависимостей и состояния, которые нужны всем методам структуры
// отвечает за получение и обработку сообщений
type Processor struct {
	tg      *telegram.Client
	offset  int // значение offset должно переживать много вызовов метода. Поэтому оно хранится внутри объекта
	storage storage.Storage
}

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
func New(client *telegram.Client, storage storage.Storage) *Processor {
	return &Processor{
		tg:      client,
		storage: storage,
	}
}

// Fetch получает сырые данные от Telegram и превращает их в удобный для приложения формат.
// Update - сущность тг, event - общая сущность
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
	res := make([]events.Event, len(updates))

	// перебираем апдейты и преобразуем их в тип ивент
	for _, u := range updates {
		res = append(res, event(u))
	}

	// Обновляем оффсет чтобы в следующий раз получить новую пачку сообщений
	p.offset = updates[len(updates)-1].UpdateId + 1

	return res, nil
}

func (p *Processor) Process(event events.Event) error {
	switch event.Type {
	case events.Message:
		return p.processMesage(event)
	default:
		return e.Wrap("cannot process message:", errors.New("unknown event type"))
	}
}

func (p *Processor) processMesage(event events.Event) error {
	meta, err := getMeta(event)
	if err != nil {
		return e.Wrap("cannot process message, getMeta:", err)
	}

	if err := p.doCmd(event.Text, meta.ChatID, meta.Username); err != nil {
		return e.Wrap("cannot process messaage, doCmd:", err)
	}
	return nil
}

func getMeta(event events.Event) (Meta, error) {
	// Попробуй взять значение из event.Meta и привести его к типу Meta
	res, ok := event.Meta.(Meta)
	if !ok {
		return Meta{}, e.Wrap("cannot get meta:", errors.New("unknown meta type"))
	}
	return res, nil
}

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

func fetchText(upd telegram.Update) string {
	if upd.Message == nil {
		return ""
	}
	return upd.Message.Text
}

func fetchType(upd telegram.Update) events.Type {
	if upd.Message == nil {
		return events.Unknown
	}
	return events.Message
}

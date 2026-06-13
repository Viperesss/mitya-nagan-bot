package telegram

import (
	"errors"
	"log"
	"net/url"
	"strings"
	"the-mitya-nagan-bot/lib/e"
	"the-mitya-nagan-bot/storage"
)

const (
	RandomCmd = "/random"
	HelpCmd   = "/help"
	StartCmd  = "/start"
)

// типо API роутер, думает что делать с полученным сообщением
func (p *Processor) doCmd(text string, chatID int, userName string) error {
	text = strings.TrimSpace(text)

	log.Printf("got new command '%s' from '%s'", text, userName)

	if isAddCmd(text) {
		return p.savePage(text, chatID, userName)
	}

	switch text {
	case RandomCmd:
		return p.sendRandom(chatID, userName)
	case HelpCmd:
		return p.sendHelp(chatID)
	case StartCmd:
		return p.sendHello(chatID)
	default:
		return p.tg.SendMessage(chatID, msgUnknownCommand)
	}
}

func (p *Processor) savePage(text string, chatID int, userName string) error {
	// страница которую собираемся сохранить
	page := &storage.Page{
		URL:      text,
		UserName: userName,
	}

	isExists, err := p.storage.IsExists(page)
	if err != nil {
		return e.Wrap("events commands savePage fail, cannot use IsExists:", err)
	}
	if isExists {
		return p.tg.SendMessage(chatID, msgAlreadyExists)
	}

	if err := p.storage.Save(page); err != nil {
		return e.Wrap("events commands savePage fail, cannot save page in the storage:", err)
	}

	if err := p.tg.SendMessage(chatID, msgSaved); err != nil {
		return e.Wrap("events commands savePage fail, cannot send message:", err)
	}
	return nil
}

func (p *Processor) sendRandom(chatID int, userName string) error {
	page, err := p.storage.PickRandom(userName)
	if err != nil && !errors.Is(err, storage.ErrNoSavedPages) {
		return e.Wrap("events commands SendRandom fail, cannot PickRandom:", err)
	}

	if errors.Is(err, storage.ErrNoSavedPages) {
		return p.tg.SendMessage(chatID, msgNoSavedPages)
	}

	if err := p.tg.SendMessage(chatID, page.URL); err != nil {
		return e.Wrap("events commands savePage fail, cannot send message:", err)
	}

	return p.storage.Remove(page)
}

func (p *Processor) sendHello(chatID int) error {
	if err := p.tg.SendMessage(chatID, msgHello); err != nil {
		return e.Wrap("events commands SendHello fail, cannot send message:", err)
	}
	return nil
}

func (p *Processor) sendHelp(chatID int) error {
	if err := p.tg.SendMessage(chatID, msgHelp); err != nil {
		return e.Wrap("events commands SendHelp fail, cannot send message:", err)
	}
	return nil
}

func isAddCmd(text string) bool {
	link, err := url.Parse(text)
	if err == nil && link.Host != "" {
		return true
	}
	return false
}

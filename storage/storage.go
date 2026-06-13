// Абстракция
package storage

import (
	"crypto/sha1"
	"errors"
	"fmt"
	"io"
	"the-mitya-nagan-bot/lib/e"
)

// Нужно для того чтобы бот мог сообщить пользователю о там что пользователь пока ничего не сохранил
var ErrNoSavedPages = errors.New("no saved page")

// сущности передаём по указателю чтобы не создавать копии
type Storage interface {
	Save(p *Page) error // encode
	Remove(p *Page) error
	PickRandom(userName string) (*Page, error)
	IsExists(p *Page) (bool, error)
}

// сущность которую и будем хранить
// страница, ссылку на которую мы скинули боту через клиент
type Page struct {
	URL      string
	UserName string
}

func (p Page) Hash() (string, error) {
	h := sha1.New()
	if _, err := io.WriteString(h, p.URL); err != nil {
		return "", e.Wrap("storage Hash fail:", err)
	}
	if _, err := io.WriteString(h, p.UserName); err != nil {
		return "", e.Wrap("storage Hash fail:", err)
	}
	return fmt.Sprintf("%x", h.Sum(nil)), nil
}

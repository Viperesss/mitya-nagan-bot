// Абстракция
// Package storage defines page storage interfaces and entities.
package storage

import (
	"context"
	"crypto/sha1"
	"errors"
	"fmt"
	"io"
	"the-mitya-nagan-bot/lib/e"
)

// Нужно для того чтобы бот мог сообщить пользователю о там что пользователь пока ничего не сохранил
// ErrNoSavedPages is returned when a user has no saved pages.
var ErrNoSavedPages = errors.New("no saved page")

// сущности передаём по указателю чтобы не создавать копии
// Storage defines operations for storing and retrieving pages.
type Storage interface {
	Save(ctx context.Context, p *Page) error // encode
	Remove(ctx context.Context, p *Page) error
	PickRandom(ctx context.Context, userName string) (*Page, error)
	IsExists(ctx context.Context, p *Page) (bool, error)
}

// сущность которую и будем хранить
// страница, ссылку на которую мы скинули боту через клиент
// Page contains information about a page saved by a user.
type Page struct {
	URL      string
	UserName string
}

// Hash returns a hash derived from the page URL and user name.
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

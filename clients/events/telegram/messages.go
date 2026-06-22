// Package telegram implements Telegram-based event processing.
package telegram

const (
	msgHelp = `
	Я тут ссылки на хранение принимаю.
	Кинешь ссылку - спрячу в схрон.

	Команда /random выдаст случайную маляву из запасов.
	После выдачи бумага уходит в расход.`

	msgHello = "Здарова тяу.\n\n" + msgHelp

	msgUnknownCommand = "Не понял базара."
	msgNoSavedPages   = "Пусто в схроне. Ни одной малявы не завалялось."
	msgSaved          = "Принял. Спрятал в общак."
	msgAlreadyExists  = "Эта малява уже в деле."
)

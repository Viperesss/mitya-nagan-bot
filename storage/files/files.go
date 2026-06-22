// Package files implements file-based page storage.
package files

import (
	"encoding/gob"
	"errors"
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"the-mitya-nagan-bot/lib/e"
	"the-mitya-nagan-bot/storage"
)

// конкретная реализация интерфейса
// Storage stores pages on the local filesystem.
type Storage struct {
	basePath string
}

const defaultPerm = 0774

// New create a new file storage using the provided base path
func New(path string) Storage {
	return Storage{basePath: path}
}

// метод интерфейса
// Save persists the page in the filestorage.
func (s Storage) Save(page *storage.Page) (err error) {
	filePath := filepath.Join(s.basePath, page.UserName)

	// Формируем путь до директории, куда будет сохраняться файл
	if err := os.MkdirAll(filePath, defaultPerm); err != nil { // создать все дериктории кооторые входят в переданный путь
		return e.Wrap("storage files Save fail, cannot create dirs:", err)
	}

	// Формируем имя файла
	fileName, err := fileName(page)
	if err != nil {
		return e.Wrap("storage files Save fail, cannot use hash:", err)
	}

	// Дописываем имя файла к пути
	filePath = filepath.Join(filePath, fileName)

	// Создаем файл
	file, err := os.Create(filePath)
	if err != nil {
		return e.Wrap("storage files Save fail, cannot create file:", err)
	}

	defer file.Close() // игнор ошибки

	// Записываем в файл страницу в нужном формате
	if err := gob.NewEncoder(file).Encode(page); err != nil { // файл в который будет записываться результат, преобразовывем страницу
		return e.Wrap("storage files Save fail, cannot encode page:", err)
	}

	return nil
}

// Remove deletes the page from storage.
func (s Storage) Remove(p *storage.Page) error {
	fileName, err := fileName(p)
	if err != nil {
		return e.Wrap("storage files Remove fail, cannot use hash:", err)
	}

	filePath := filepath.Join(s.basePath, p.UserName, fileName)

	err = os.Remove(filePath)
	if err != nil {
		return fmt.Errorf("storage files Remove fail, cannot remove file: %s, error: %w", filePath, err)
	}

	return nil
}

// PickRandom returns a random page saved by the user.
func (s Storage) PickRandom(userName string) (page *storage.Page, err error) {
	path := filepath.Join(s.basePath, userName)

	files, err := os.ReadDir(path)
	if err != nil {
		return nil, e.Wrap("storage files PickRandom fail, cannot read dir:", err)
	}

	if len(files) == 0 {
		return nil, storage.ErrNoSavedPages
	}

	n := rand.Intn(len(files))
	file := files[n]

	// декодировать и вернуть
	return s.decodePage(filepath.Join(path, file.Name()))
}

// IsExists reports whether the page is already stored.
func (s Storage) IsExists(p *storage.Page) (bool, error) {
	fileName, err := fileName(p)
	if err != nil {
		return false, e.Wrap("storage files IsExists fail, cannot use hash:", err)
	}

	filePath := filepath.Join(s.basePath, p.UserName, fileName)

	switch _, err = os.Stat(filePath); {
	case errors.Is(err, os.ErrNotExist):
		return false, nil
	case err != nil:
		return false, fmt.Errorf("storage files IsExists fail, cannot check if file %s exists: %w", filePath, err)
	}

	return true, nil
}

// decodePage reads and decode a pages from a file.
func (s Storage) decodePage(filePath string) (*storage.Page, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, e.Wrap("storage files PickRandom fail, cannot open file:", err)
	}
	defer file.Close()

	// Туда файл будет декодирован
	var p storage.Page

	if err := gob.NewDecoder(file).Decode(&p); err != nil { // декодируется не копия
		return nil, e.Wrap("storage PickRandom fail, cannot decode file:", err)
	}

	return &p, nil
}

// Если в будующем захотим поменять способ именования
// fileName returns the filename for the page.
func fileName(p *storage.Page) (string, error) {
	return p.Hash()
}

// Для получения случайной страницы хэш не нужен - берём любой файл и читаем.
// Для поиска конкретной страницы хэш нужен - он позволяет сразу вычислить имя файла и обратиться к нему напрямую без перебора.

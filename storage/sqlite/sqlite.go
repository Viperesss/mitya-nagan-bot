package sqlite

import (
	"context"
	"database/sql" // общий интерфейс взаимодействия со всеми реляционными БД
	"the-mitya-nagan-bot/lib/e"
	"the-mitya-nagan-bot/storage"

	_ "github.com/mattn/go-sqlite3"
)

type Storage struct {
	db *sql.DB
}

func New(path string) (*Storage, error) {
	db, err := sql.Open("sqlite3", path)
	if err != nil {
		return nil, e.Wrap("storage sqlite New fail, cannot create new Storage:", err)
	}

	// проверка соединения с файлом
	if err := db.Ping(); err != nil {
		return nil, e.Wrap("storage sqlite New fail, cannot Ping DB:", err)
	}

	return &Storage{db: db}, nil
}

// Save(p *Page) error // encode
// 	Remove(p *Page) error
// 	PickRandom(userName string) (*Page, error)
// 	IsExists(p *Page) (bool, error)

func (s *Storage) Save(ctx context.Context, p *storage.Page) error {
	q := `INSERT INTO pages (url, user_name) VALUES (?, ?)` // ?, ? - аргументы запроса
	_, err := s.db.ExecContext(ctx, q, p.URL, p.UserName)   // противодействие SQL-инъекциям
	if err != nil {
		return e.Wrap("storage sqlite Save fail, cannot execute query:", err)
	}
	return nil
}

func (s *Storage) PickRandom(ctx context.Context, userName string) (*storage.Page, error) {
	q := `SELECT url FROM pages WHERE user_name = ? ORDER BY RANDOM() LIMIT 1`

	var url string

	err := s.db.QueryRowContext(ctx, q, userName).Scan(&url)
	// Если пусто
	if err == sql.ErrNoRows {
		return nil, storage.ErrNoSavedPages
	}

	if err != nil {
		return nil, e.Wrap("storage sqlite PickRandom fail, cannot execute query:", err)
	}

	// url получили из запроса
	return &storage.Page{URL: url, UserName: userName}, nil
}

func (s *Storage) Remove(ctx context.Context, p *storage.Page) error {
	q := `DELETE FROM pages WHERE url = ? AND user_name = ?`
	_, err := s.db.ExecContext(ctx, q, p.URL, p.UserName)
	if err != nil {
		return e.Wrap("storage sqlite Remove fail, cannot execute query:", err)
	}

	return nil
}

// IsExists checks if page exists in storage.
func (s *Storage) IsExists(ctx context.Context, p *storage.Page) (bool, error) {
	q := `SELECT COUNT(*) FROM pages WHERE url = ? AND user_name = ?`

	var count int

	if err := s.db.QueryRowContext(ctx, q, p.URL, p.UserName).Scan(&count); err != nil {
		return false, e.Wrap("storage sqlite isExists fail, cannot execute query:", err)
	}

	return count > 0, nil
}

func (s *Storage) Init(ctx context.Context) error {
	q := `CREATE TABLE IF NOT EXISTS pages (url TEXT, user_name TEXT)`

	_, err := s.db.ExecContext(ctx, q)
	if err != nil {
		return e.Wrap("storage sqlite Init fail, cannot execute query:", err)
	}
	return nil
}

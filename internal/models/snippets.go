package models

import (
	"database/sql"
	"time"
)

type Snippet struct {
	ID      int
	Title   string
	Content string
	Created time.Time
	Expires time.Time
}

type SnippetModel struct {
	DB *sql.DB
}

func (m *SnippetModel) Insert(title string, content string, expires int) (int, error) {
	stmt := `
	INSERT INTO snippets(title, content, created, expires)
	VALUES ($1, $2, NOW(), NOW() + INTERVAL '1 day' * $3) RETURNING id
	`

	id := 0
	if err := m.DB.QueryRow(stmt, title, content, expires).Scan(&id); err != nil {
		return 0, err
	}

	return id, nil
}

func (m *SnippetModel) Get(id int) (*Snippet, error) {
	return nil, nil
}

func (m *SnippetModel) Latest() ([]*Snippet, error) {
	return nil, nil
}

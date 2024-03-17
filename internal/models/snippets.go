package models

import (
	"database/sql"
	"errors"
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
	stmt := `
	SELECT id, title, content, created, expires
	FROM snippets
	WHERE expires > NOW() AND id = $1
	`

	var snippet Snippet
	err := m.DB.QueryRow(stmt, id).Scan(
		&snippet.ID,
		&snippet.Title,
		&snippet.Content,
		&snippet.Created,
		&snippet.Expires,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNoRecord
		}

		return nil, err
	}

	return &snippet, nil
}

func (m *SnippetModel) Latest() ([]*Snippet, error) {
	return nil, nil
}

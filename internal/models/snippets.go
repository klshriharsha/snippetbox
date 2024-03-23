package models

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Snippet struct {
	ID      int
	Title   string
	Content string
	Created time.Time
	Expires time.Time
}

type SnippetModel struct {
	DB *pgxpool.Pool
}

func (m *SnippetModel) Insert(title string, content string, expires int) (int, error) {
	stmt := `
	INSERT INTO snippets(title, content, created, expires)
	VALUES ($1, $2, NOW(), NOW() + INTERVAL '1 day' * $3) RETURNING id
	`

	id := 0
	if err := m.DB.QueryRow(context.Background(), stmt, title, content, expires).Scan(&id); err != nil {
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
	err := m.DB.QueryRow(context.Background(), stmt, id).Scan(
		&snippet.ID,
		&snippet.Title,
		&snippet.Content,
		&snippet.Created,
		&snippet.Expires,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNoRecord
		}

		return nil, err
	}

	return &snippet, nil
}

func (m *SnippetModel) Latest() ([]*Snippet, error) {
	stmt := `
	SELECT id, title, content, created, expires
	FROM snippets
	WHERE expires > NOW()
	ORDER BY id DESC
	LIMIT 10
	`
	rows, err := m.DB.Query(context.Background(), stmt)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var snippets []*Snippet

	for rows.Next() {
		var snippet Snippet

		err := rows.Scan(&snippet.ID, &snippet.Title, &snippet.Content, &snippet.Created, &snippet.Expires)
		if err != nil {
			return nil, err
		}

		snippets = append(snippets, &snippet)
	}

	// if there was an error during iteration, it is returned by `rows.Err()`
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return snippets, nil
}

package models

import (
	"database/sql"
	"errors"
	"time"

	"github.com/google/uuid"
)

type Snippet struct {
	ID      string
	Title   string
	Content string
	Created time.Time
	Expires time.Time
}

type SnippetModel struct {
	DB *sql.DB
}

type SnippetModelInterface interface {
	Insert(title string, content string, expires int) (string, error)
	Get(id string) (Snippet, error)
	Latest() ([]Snippet, error)
}

func (m *SnippetModel) Insert(title string, content string, expires int) (string, error) {
	query := `INSERT INTO snippets (id, title, content, created, expires)
            VALUES(?, ?, ?, DATETIME('now', 'utc'), DATETIME('now', 'utc', '+' || ? || ' days'))`
	id := uuid.New().String()

	_, err := m.DB.Exec(query, id, title, content, expires)
	if err != nil {
		return "", err
	}

	return id, nil
}

func (m *SnippetModel) Get(id string) (Snippet, error) {
	var s Snippet
	query := `SELECT id, title, content, created, expires
            FROM snippets
            WHERE expires > DATETIME('now', 'utc') AND id = ?`

	err := m.DB.QueryRow(query, id).Scan(&s.ID, &s.Title, &s.Content, &s.Created, &s.Expires)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return Snippet{}, ErrNoRecord
		} else {
			return Snippet{}, err
		}
	}

	return s, nil
}

func (m *SnippetModel) Latest() ([]Snippet, error) {
	query := `SELECT id, title, content, created, expires
            FROM snippets
            WHERE expires > DATETIME('now', 'utc') ORDER BY id DESC LIMIT 10`

	rows, err := m.DB.Query(query)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var snippets []Snippet

	for rows.Next() {
		var s Snippet

		err = rows.Scan(&s.ID, &s.Title, &s.Content, &s.Created, &s.Expires)
		if err != nil {
			return nil, err
		}

		snippets = append(snippets, s)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return snippets, nil
}

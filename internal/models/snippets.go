package models

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/jeremydwayne/snippets/internal/sqlc"
)

type Snippet struct {
	ID      string    `json:"id"`
	Title   string    `json:"title"`
	Content string    `json:"content"`
	Created time.Time `json:"created"`
	Expires time.Time `json:"expires"`
}

type SnippetModel struct {
	DB *sqlc.Queries
}

type SnippetModelInterface interface {
	Insert(title string, content string, expires int) (string, error)
	Get(id string) (Snippet, error)
	Latest() ([]Snippet, error)
}

func (m *SnippetModel) Insert(title string, content string, expires int) (string, error) {
	id := uuid.New().String()

	expire := time.Now().UTC().AddDate(0, 0, expires)
	params := sqlc.CreateSnippetParams{
		ID:      id,
		Title:   title,
		Content: content,
		Expires: expire,
	}

	err := m.DB.CreateSnippet(context.Background(), params)
	if err != nil {
		return "", err
	}

	return id, nil
}

func (m *SnippetModel) Get(id string) (Snippet, error) {
	s, err := m.DB.GetSnippet(context.Background(), id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return Snippet{}, ErrNoRecord
		} else {
			return Snippet{}, err
		}
	}
	snippet := Snippet{
		ID:      s.ID,
		Title:   s.Title,
		Content: s.Content,
		Created: s.Created,
		Expires: s.Expires,
	}

	return snippet, nil
}

func (m *SnippetModel) Latest() ([]Snippet, error) {
	rows, err := m.DB.LatestSnippets(context.Background())
	if err != nil {
		return nil, err
	}

	var snippets []Snippet

	for _, row := range rows {
		snippets = append(snippets, convertSnippet(row))
	}

	return snippets, nil
}

func convertSnippet(row sqlc.Snippets) Snippet {
	return Snippet{
		ID:      row.ID,
		Title:   row.Title,
		Content: row.Content,
		Created: row.Created,
		Expires: row.Expires,
	}
}

package mocks

import (
	"time"

	"github.com/jeremydwayne/snippets/internal/models"
)

var mockSnippet = models.Snippet{
	ID:      "44bfa272-d0b4-4b09-966f-eae71ddaf304",
	Title:   "An old silent pond",
	Content: "An old silent pond...",
	Created: time.Now(),
	Expires: time.Now(),
}

type SnippetModel struct{}

func (m *SnippetModel) Insert(title string, content string, expires int) (string, error) {
	return "44bfa272-d0b4-4b09-966f-eae71ddaf304", nil
}

func (m *SnippetModel) Get(id string) (models.Snippet, error) {
	switch id {
	case "44bfa272-d0b4-4b09-966f-eae71ddaf304":
		return mockSnippet, nil
	default:
		return models.Snippet{}, models.ErrNoRecord
	}
}

func (m *SnippetModel) Latest() ([]models.Snippet, error) {
	return []models.Snippet{mockSnippet}, nil
}

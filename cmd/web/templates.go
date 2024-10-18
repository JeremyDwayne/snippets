package main

import "github.com/jeremydwayne/snippets/internal/models"

type templateData struct {
	Snippet  models.Snippet
	Snippets []models.Snippet
}

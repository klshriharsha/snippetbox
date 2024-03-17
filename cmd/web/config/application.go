package config

import (
	"html/template"
	"log"

	"github.com/klshriharsha/snippetbox/internal/models"
)

// application is used for dependency injection throughout the `web` application
type Application struct {
	InfoLog  *log.Logger
	ErrorLog *log.Logger

	Snippets *models.SnippetModel

	TemplateCache map[string]*template.Template
}

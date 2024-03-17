package config

import (
	"html/template"
	"log"

	"github.com/klshriharsha/snippetbox/internal/models"
)

// application is used for dependency injection throughout the `web` application
type Application struct {
	// config
	InfoLog  *log.Logger
	ErrorLog *log.Logger

	// models
	Snippets *models.SnippetModel

	// other
	TemplateCache map[string]*template.Template
}

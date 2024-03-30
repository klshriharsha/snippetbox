package config

import (
	"html/template"
	"log"

	"github.com/alexedwards/scs/v2"
	"github.com/go-playground/form"
	"github.com/klshriharsha/snippetbox/internal/models"
)

// application is used for dependency injection throughout the `web` application
type Application struct {
	// config
	InfoLog   *log.Logger
	ErrorLog  *log.Logger
	DebugMode bool

	// models
	// Snippets exposes database operations related to snippets
	Snippets models.SnippetModelInterface
	Users    models.UserModelInterface

	// other
	// TemplateCache holds all the parsed templates in memory
	TemplateCache  map[string]*template.Template
	FormDecoder    *form.Decoder
	SessionManager *scs.SessionManager
}

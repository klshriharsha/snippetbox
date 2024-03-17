package config

import (
	"net/http"
	"time"

	"github.com/klshriharsha/snippetbox/internal/models"
)

// TemplateData holds all the data passed to Go templates
type TemplateData struct {
	Snippet     *models.Snippet
	Snippets    []*models.Snippet
	CurrentYear int
}

// NewTemplateData creates a new `TemplateData` with `CurrentYear` initialized
func NewTemplateData(r *http.Request) *TemplateData {
	return &TemplateData{CurrentYear: time.Now().Year()}
}

package config

import (
	"net/http"
	"time"

	"github.com/klshriharsha/snippetbox/internal/models"
)

type TemplateData struct {
	Snippet     *models.Snippet
	Snippets    []*models.Snippet
	CurrentYear int
}

func NewTemplateData(r *http.Request) *TemplateData {
	return &TemplateData{CurrentYear: time.Now().Year()}
}

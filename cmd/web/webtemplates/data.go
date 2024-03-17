package webtemplates

import "github.com/klshriharsha/snippetbox/internal/models"

type TemplateData struct {
	Snippet  *models.Snippet
	Snippets []*models.Snippet
}

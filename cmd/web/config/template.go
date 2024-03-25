package config

import (
	"html/template"
	"net/http"
	"path/filepath"
	"time"

	"github.com/justinas/nosurf"
	"github.com/klshriharsha/snippetbox/internal/models"
)

// TemplateData holds all the data passed to Go templates
type TemplateData struct {
	Snippet         *models.Snippet
	Snippets        []*models.Snippet
	CurrentYear     int
	Form            any
	Flash           string
	IsAuthenticated bool
	CSRFToken       string
}

// NewTemplateData creates a new `TemplateData` with `CurrentYear` initialized
func (app *Application) NewTemplateData(r *http.Request) *TemplateData {
	return &TemplateData{
		CurrentYear:     time.Now().Year(),
		Flash:           app.SessionManager.PopString(r.Context(), "flash"),
		IsAuthenticated: app.isAuthenticated(r),
		CSRFToken:       nosurf.Token(r),
	}
}

// NewTemplateCache initializes the template cache by parsing all page and partial templates and
// holding them in memory to avoid disk access at runtime
func NewTemplateCache() (map[string]*template.Template, error) {
	pages, err := filepath.Glob("./ui/html/pages/*.tmpl")
	if err != nil {
		return nil, err
	}

	cache := map[string]*template.Template{}
	for _, page := range pages {
		filename := filepath.Base(page)

		// parse the base template file
		// before parsing any templates register the custome template functions
		ts, err := template.New(filename).Funcs(functions).ParseFiles("./ui/html/base.go.tmpl")
		if err != nil {
			return nil, err
		}
		// parse all partials into the same template set
		ts, err = ts.ParseGlob("./ui/html/partials/*.tmpl")
		if err != nil {
			return nil, err
		}
		// parse the main template into the same template set
		ts, err = ts.ParseFiles(page)
		if err != nil {
			return nil, err
		}

		cache[filename] = ts
	}

	return cache, nil
}

// standardDate formats date into a human-readable format
func standardDate(t time.Time) string {
	return t.Format("02 Jan 2006 at 15:04")
}

// map of custom template functions to register before parsing any template
var functions = template.FuncMap{
	"standardDate": standardDate,
}

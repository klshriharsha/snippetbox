package config

import (
	"html/template"
	"path/filepath"
)

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

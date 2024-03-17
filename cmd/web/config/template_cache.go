package config

import (
	"html/template"
	"path/filepath"
)

func NewTemplateCache() (map[string]*template.Template, error) {
	pages, err := filepath.Glob("./ui/html/pages/*.tmpl")
	if err != nil {
		return nil, err
	}

	cache := map[string]*template.Template{}
	for _, page := range pages {
		filename := filepath.Base(page)

		files := []string{
			"./ui/html/base.go.tmpl",
			"./ui/html/partials/nav.go.tmpl",
			page,
		}
		// parses the Go template files into a set of templates
		ts, err := template.ParseFiles(files...)
		if err != nil {
			return nil, err
		}

		cache[filename] = ts
	}

	return cache, nil
}

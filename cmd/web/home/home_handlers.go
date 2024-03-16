package home

import (
	"html/template"
	"net/http"

	"github.com/klshriharsha/snippetbox/cmd/web/config"
)

// home defines the homepage route.
// writes the necessary html to response body
func HomeHandler(app *config.Application) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			app.NotFoundError(w)
			return
		}

		files := []string{
			"./ui/html/base.go.tmpl",
			"./ui/html/partials/nav.go.tmpl",
			"./ui/html/pages/home.go.tmpl",
		}

		// parses the Go template files into a set of templates
		ts, err := template.ParseFiles(files...)
		if err != nil {
			app.ServerError(w, err)
			return
		}

		// executes the `base` template which invokes other templates
		err = ts.ExecuteTemplate(w, "base", nil)
		if err != nil {
			app.ServerError(w, err)
		}
	}
}

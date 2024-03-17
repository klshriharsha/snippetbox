package home

import (
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

		snippets, err := app.Snippets.Latest()
		if err != nil {
			app.ServerError(w, err)
			return
		}

		app.RenderPage(w, http.StatusOK, "home.go.tmpl", &config.TemplateData{Snippets: snippets})
	}
}

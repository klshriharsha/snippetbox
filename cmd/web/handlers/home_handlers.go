package handlers

import (
	"net/http"

	"github.com/klshriharsha/snippetbox/cmd/web/config"
)

// HomeHandler display a list of the latest 10 un-expired snippets
func HomeHandler(app *config.Application) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		snippets, err := app.Snippets.Latest()
		if err != nil {
			app.ServerError(w, err)
			return
		}

		data := app.NewTemplateData(r)
		data.Snippets = snippets

		app.RenderPage(w, http.StatusOK, "home.go.tmpl", data)
	}
}

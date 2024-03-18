package handlers

import (
	"net/http"

	"github.com/klshriharsha/snippetbox/cmd/web/config"
	"github.com/klshriharsha/snippetbox/cmd/web/render"
)

// HomeHandler display a list of the latest 10 un-expired snippets
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

		data := render.NewTemplateData(r)
		data.Snippets = snippets

		app.RenderPage(w, http.StatusOK, "home.go.tmpl", data)
	}
}

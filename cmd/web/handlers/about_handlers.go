package handlers

import (
	"net/http"

	"github.com/klshriharsha/snippetbox/cmd/web/config"
)

// aboutHandler displays an about page for snippetbox
func aboutHandler(app *config.Application) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		data := app.NewTemplateData(r)
		app.RenderPage(w, http.StatusOK, "about.go.tmpl", data)
	})
}

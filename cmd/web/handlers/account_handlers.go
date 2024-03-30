package handlers

import (
	"net/http"

	"github.com/klshriharsha/snippetbox/cmd/web/config"
)

// AccountViewHandler displays the information of the currently logged in user
func accountViewHandler(app *config.Application) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := app.SessionManager.GetInt(r.Context(), "authenticatedUserID")
		user, _ := app.Users.GetByID(id)

		data := app.NewTemplateData(r)
		data.User = user
		app.RenderPage(w, http.StatusOK, "account-view.go.tmpl", data)
	})
}

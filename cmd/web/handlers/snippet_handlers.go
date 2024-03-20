package handlers

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
	"github.com/klshriharsha/snippetbox/cmd/web/config"
	"github.com/klshriharsha/snippetbox/cmd/web/render"
	"github.com/klshriharsha/snippetbox/internal/models"
)

// SnippetViewHandler displays the snippet corresponding to the `id` in the query parameters
func SnippetViewHandler(app *config.Application) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// get the snippet id from named parameters in the request context
		params := httprouter.ParamsFromContext(r.Context())
		id, err := strconv.Atoi(params.ByName("id"))
		if err != nil || id < 1 {
			app.NotFoundError(w)
			return
		}

		snippet, err := app.Snippets.Get(id)
		if err != nil {
			if errors.Is(err, models.ErrNoRecord) {
				app.NotFoundError(w)
				return
			}

			app.ServerError(w, err)
			return
		}

		data := render.NewTemplateData(r)
		data.Snippet = snippet

		app.RenderPage(w, http.StatusOK, "view.go.tmpl", data)
	}
}

func SnippetCreateHandler(app *config.Application) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("display a form to create snippet"))
	}
}

// SnippetCreatePostHandler creates a new snippet in the database and sends a redirect response
// to view the created snippet
func SnippetCreatePostHandler(app *config.Application) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		title := "O snail"
		content := "O snail\nClimb Mount Fuji,\nBut slowly, slowly!\n\n-Kobayashi Issa"
		expires := 7

		id, err := app.Snippets.Insert(title, content, expires)
		if err != nil {
			app.ServerError(w, err)
			return
		}

		http.Redirect(w, r, fmt.Sprintf("/snippet/view?id=%d", id), http.StatusSeeOther)
	}
}

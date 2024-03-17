package snippet

import (
	"errors"
	"fmt"
	"html/template"
	"net/http"
	"strconv"

	"github.com/klshriharsha/snippetbox/cmd/web/config"
	"github.com/klshriharsha/snippetbox/cmd/web/webtemplates"
	"github.com/klshriharsha/snippetbox/internal/models"
)

func SnippetViewHandler(app *config.Application) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.Atoi(r.URL.Query().Get("id"))
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

		files := []string{
			"./ui/html/base.go.tmpl",
			"./ui/html/partials/nav.go.tmpl",
			"./ui/html/pages/view.go.tmpl",
		}

		ts, err := template.ParseFiles(files...)
		if err != nil {
			app.ServerError(w, err)
			return
		}

		if err := ts.ExecuteTemplate(w, "base", webtemplates.TemplateData{Snippet: snippet}); err != nil {
			app.ServerError(w, err)
		}
	}
}

func SnippetCreateHandler(app *config.Application) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			w.Header().Set("Allow", http.MethodPost)
			app.ClientError(w, http.StatusMethodNotAllowed)
			return
		}

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

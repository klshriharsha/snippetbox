package handlers

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
	"github.com/klshriharsha/snippetbox/cmd/web/config"
	"github.com/klshriharsha/snippetbox/internal/models"
	"github.com/klshriharsha/snippetbox/internal/validator"
)

// SnippetViewHandler displays the snippet corresponding to the `id` in the query parameters
func snippetViewHandler(app *config.Application) http.HandlerFunc {
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

		data := app.NewTemplateData(r)
		data.Snippet = snippet

		app.RenderPage(w, http.StatusOK, "view.go.tmpl", data)
	}
}

type snippetCreateFrom struct {
	Title               string `form:"title"`
	Content             string `form:"content"`
	Expires             int    `form:"expires"`
	validator.Validator `form:"-"`
}

func snippetCreateHandler(app *config.Application) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		data := app.NewTemplateData(r)
		data.Form = snippetCreateFrom{Expires: 365}

		app.RenderPage(w, http.StatusOK, "create.go.tmpl", data)
	}
}

// SnippetCreatePostHandler creates a new snippet in the database and sends a redirect response
// to view the created snippet
func snippetCreatePostHandler(app *config.Application) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		form := snippetCreateFrom{Validator: validator.Validator{FieldErrors: make(map[string]string)}}
		if err := app.DecodePostForm(r, &form); err != nil {
			app.ClientError(w, http.StatusBadRequest)
			return
		}

		form.CheckField(validator.NotBlank(form.Title), "title", "Title cannot be empty")
		form.CheckField(validator.MaxChars(form.Content, 100), "content", "Content cannot be empty")
		form.CheckField(validator.ValidValue(form.Expires, 1, 7, 365), "expires", "Expires can only be 1, 7 or 365")
		if !form.Valid() {
			// if there are validation errors, render the same template with original field
			// values and field errors
			data := app.NewTemplateData(r)
			data.Form = form
			app.RenderPage(w, http.StatusUnprocessableEntity, "create.go.tmpl", data)
			return
		}

		id, err := app.Snippets.Insert(form.Title, form.Content, form.Expires)
		if err != nil {
			app.ServerError(w, err)
			return
		}
		app.SessionManager.Put(r.Context(), "flash", "Snippet created successfully!")

		http.Redirect(w, r, fmt.Sprintf("/snippet/view/%d", id), http.StatusSeeOther)
	}
}

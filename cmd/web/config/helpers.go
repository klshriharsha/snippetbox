package config

import (
	"fmt"
	"net/http"
	"runtime/debug"
)

func (app *Application) ServerError(w http.ResponseWriter, err error) {
	trace := fmt.Sprintf("%s\n%s", err.Error(), debug.Stack())
	app.ErrorLog.Output(2, trace)

	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

func (app *Application) ClientError(w http.ResponseWriter, status int) {
	http.Error(w, http.StatusText(status), status)
}

func (app *Application) NotFoundError(w http.ResponseWriter) {
	app.ClientError(w, http.StatusNotFound)
}

func (app *Application) RenderPage(w http.ResponseWriter, status int, page string, data *TemplateData) {
	ts, ok := app.TemplateCache[page]
	if !ok {
		err := fmt.Errorf("template %s does not exist", page)
		app.ServerError(w, err)
		return
	}

	w.WriteHeader(status)
	// executes the `base` template which invokes other templates
	if err := ts.ExecuteTemplate(w, "base", data); err != nil {
		app.ServerError(w, err)
	}
}

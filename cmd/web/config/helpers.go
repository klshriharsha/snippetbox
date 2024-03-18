package config

import (
	"bytes"
	"fmt"
	"net/http"
	"runtime/debug"

	"github.com/klshriharsha/snippetbox/cmd/web/render"
)

// ServerError responds with an HTTP 500 error
func (app *Application) ServerError(w http.ResponseWriter, err error) {
	trace := fmt.Sprintf("%s\n%s", err.Error(), debug.Stack())
	app.ErrorLog.Output(2, trace)

	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

// ClientError responds with an HTTP error of gives status code
func (app *Application) ClientError(w http.ResponseWriter, status int) {
	http.Error(w, http.StatusText(status), status)
}

// NotFoundError responds with an HTTP 404 error
func (app *Application) NotFoundError(w http.ResponseWriter) {
	app.ClientError(w, http.StatusNotFound)
}

// RenderPage finds the template corresponding to `page` in cache and renders it
func (app *Application) RenderPage(w http.ResponseWriter, status int, page string, data *render.TemplateData) {
	ts, ok := app.TemplateCache[page]
	if !ok {
		err := fmt.Errorf("template %s does not exist", page)
		app.ServerError(w, err)
		return
	}

	buf := new(bytes.Buffer)
	// execute the `base` template which invokes other templates
	// attempt to write to the buffer first and return an error if there is one
	if err := ts.ExecuteTemplate(buf, "base", data); err != nil {
		app.ServerError(w, err)
		return
	}

	// if the previous write to buffer was successful, write a 200 OK status with the right response
	w.WriteHeader(status)
	buf.WriteTo(w)
}

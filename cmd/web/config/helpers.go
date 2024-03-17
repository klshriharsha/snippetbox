package config

import (
	"bytes"
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

package main

import (
	"net/http"

	"github.com/klshriharsha/snippetbox/cmd/web/config"
	"github.com/klshriharsha/snippetbox/cmd/web/handlers"
)

// routes register allt he routes and middleware and returns a final handler
func routes(app *config.Application) http.Handler {
	mux := http.NewServeMux()

	fs := http.FileServer(staticFileSystem{http.Dir("./ui/static/")})
	// file server looks for the file under `./ui/static/`
	// so strip the `/static` prefix from request URL
	mux.Handle("/static/", http.StripPrefix("/static", fs))

	mux.HandleFunc("/", handlers.HomeHandler(app))
	mux.HandleFunc("/snippet/view", handlers.SnippetViewHandler(app))
	mux.HandleFunc("/snippet/create", handlers.SnippetCreateHandler(app))

	// LogRequestMiddleware logs information about every request
	// secureHeaders middleware runs before any request hits the mux so that all the important
	// headers are set in every response
	return app.RecoverPanic(app.LogRequestMiddleware(secureHeaders(mux)))
}

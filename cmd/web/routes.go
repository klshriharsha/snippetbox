package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/justinas/alice"
	"github.com/klshriharsha/snippetbox/cmd/web/config"
	"github.com/klshriharsha/snippetbox/cmd/web/handlers"
)

// routes register allt he routes and middleware and returns a final handler
func routes(app *config.Application) http.Handler {
	router := httprouter.New()

	// setup 404 handler with httprouter so that error handling is consistent
	router.NotFound = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		app.NotFoundError(w)
	})

	dynamic := alice.New(app.SessionManager.LoadAndSave)

	router.Handler(http.MethodGet, "/", dynamic.Then(handlers.HomeHandler(app)))
	router.Handler(http.MethodGet, "/snippet/view/:id", dynamic.Then(handlers.SnippetViewHandler(app)))
	router.Handler(http.MethodGet, "/snippet/create", dynamic.Then(handlers.SnippetCreateHandler(app)))
	router.Handler(http.MethodPost, "/snippet/create", dynamic.Then(handlers.SnippetCreatePostHandler(app)))

	// file server looks for the file under `./ui/static/`
	// so strip the `/static` prefix from request URL
	fs := http.FileServer(staticFileSystem{http.Dir("./ui/static/")})
	router.Handler(http.MethodGet, "/static/*filepath", http.StripPrefix("/static", fs))

	// alice.New simplifies the process of chaining and composing middleware
	// LogRequestMiddleware logs information about every request
	// secureHeaders middleware runs before any request hits the mux so that all the important
	// headers are set in every response
	standard := alice.New(app.RecoverPanic, app.LogRequestMiddleware, secureHeaders)

	return standard.Then(router)
}

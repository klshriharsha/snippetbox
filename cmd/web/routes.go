package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/justinas/alice"
	"github.com/klshriharsha/snippetbox/cmd/web/config"
	"github.com/klshriharsha/snippetbox/cmd/web/handlers"
	"github.com/klshriharsha/snippetbox/ui"
)

// routes register allt he routes and middleware and returns a final handler
func routes(app *config.Application) http.Handler {
	router := httprouter.New()

	// setup 404 handler with httprouter so that error handling is consistent
	router.NotFound = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		app.NotFoundError(w)
	})

	// LoadAndSave middleware initializes the session manager from the request context
	// noSurf middleware handles CSRF tokens on all pages (logout appears on all pages)
	dynamic := alice.New(app.SessionManager.LoadAndSave, noSurf, app.Authenticate)

	router.Handler(http.MethodGet, "/", dynamic.Then(handlers.HomeHandler(app)))

	router.Handler(http.MethodGet, "/user/signup", dynamic.Then(handlers.SignupHandler(app)))
	router.Handler(http.MethodPost, "/user/signup", dynamic.Then(handlers.SignupPostHandler(app)))
	router.Handler(http.MethodGet, "/user/login", dynamic.Then(handlers.LoginHandler(app)))
	router.Handler(http.MethodPost, "/user/login", dynamic.Then(handlers.LoginPostHandler(app)))

	router.Handler(http.MethodGet, "/snippet/view/:id", dynamic.Then(handlers.SnippetViewHandler(app)))

	// protects relevant routes with an authorization middleware
	protected := dynamic.Append(app.RequireAuth)

	router.Handler(http.MethodGet, "/snippet/create", protected.Then(handlers.SnippetCreateHandler(app)))
	router.Handler(http.MethodPost, "/snippet/create", protected.Then(handlers.SnippetCreatePostHandler(app)))

	router.Handler(http.MethodPost, "/user/logout", protected.Then(handlers.LogoutPostHandler(app)))

	fs := http.FileServer(http.FS(ui.Files))
	router.Handler(http.MethodGet, "/static/*filepath", fs)

	// alice.New simplifies the process of chaining and composing middleware
	// LogRequestMiddleware logs information about every request
	// secureHeaders middleware runs before any request hits the mux so that all the important
	// headers are set in every response
	standard := alice.New(app.RecoverPanic, app.LogRequestMiddleware, secureHeaders)

	return standard.Then(router)
}

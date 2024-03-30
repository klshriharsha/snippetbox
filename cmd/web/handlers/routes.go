package handlers

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/justinas/alice"
	"github.com/klshriharsha/snippetbox/cmd/web/config"
	"github.com/klshriharsha/snippetbox/cmd/web/middleware"
	"github.com/klshriharsha/snippetbox/ui"
)

// routes register allt he routes and middleware and returns a final handler
func Routes(app *config.Application) http.Handler {
	router := httprouter.New()

	// setup 404 handler with httprouter so that error handling is consistent
	router.NotFound = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		app.NotFoundError(w)
	})
	// simple ping endpoint
	router.Handler(http.MethodGet, "/ping", pingHandler(app))

	// LoadAndSave middleware initializes the session manager from the request context
	// noSurf middleware handles CSRF tokens on all pages (logout appears on all pages)
	dynamic := alice.New(app.SessionManager.LoadAndSave, middleware.NoSurf, app.Authenticate)

	router.Handler(http.MethodGet, "/", dynamic.Then(homeHandler(app)))
	// about page
	router.Handler(http.MethodGet, "/about", dynamic.Then(aboutHandler(app)))

	router.Handler(http.MethodGet, "/user/signup", dynamic.Then(signupHandler(app)))
	router.Handler(http.MethodPost, "/user/signup", dynamic.Then(signupPostHandler(app)))
	router.Handler(http.MethodGet, "/user/login", dynamic.Then(loginHandler(app)))
	router.Handler(http.MethodPost, "/user/login", dynamic.Then(loginPostHandler(app)))

	router.Handler(http.MethodGet, "/snippet/view/:id", dynamic.Then(snippetViewHandler(app)))

	// protects relevant routes with an authorization middleware
	protected := dynamic.Append(app.RequireAuth)

	router.Handler(http.MethodGet, "/snippet/create", protected.Then(snippetCreateHandler(app)))
	router.Handler(http.MethodPost, "/snippet/create", protected.Then(snippetCreatePostHandler(app)))

	router.Handler(http.MethodPost, "/user/logout", protected.Then(logoutPostHandler(app)))

	fs := http.FileServer(http.FS(ui.Files))
	router.Handler(http.MethodGet, "/static/*filepath", fs)

	// alice.New simplifies the process of chaining and composing middleware
	// LogRequestMiddleware logs information about every request
	// secureHeaders middleware runs before any request hits the mux so that all the important
	// headers are set in every response
	standard := alice.New(app.RecoverPanic, app.LogRequestMiddleware, middleware.SecureHeaders)

	return standard.Then(router)
}

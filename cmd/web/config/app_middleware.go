package config

import (
	"context"
	"fmt"
	"net/http"
)

// ##########################################################
// Middleware defined here are defined as methods on the
// Application struct as opposed to regular functions in
// cmd/web/middleware.go because these middleware need access
// to methods and properties on the Application struct
// ##########################################################

// LogRequestMiddleware runs before every request hits the mux and logs information about it
func (app *Application) LogRequestMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		app.InfoLog.Printf("%s - %s %s %s", r.RemoteAddr, r.Proto, r.Method, r.URL.RequestURI())
		next.ServeHTTP(w, r)
	})
}

func (app *Application) RecoverPanic(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// deferred functions are called when there is a panic and Go unwinds the call stack
		defer func() {
			// use the built-in recover function to check if there was a panic and respond with
			// a sensible 500 server error
			if err := recover(); err != nil {
				w.Header().Set("Connection", "close")
				app.ServerError(w, fmt.Errorf("%s", err))
			}
		}()

		next.ServeHTTP(w, r)
	})
}

// RequireAuth middleware redirects unauthenticated requests to the login route
func (app *Application) RequireAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !app.isAuthenticated(r) {
			http.Redirect(w, r, "/user/login", http.StatusSeeOther)
			return
		}

		w.Header().Add("Cache-Control", "no-store")
		next.ServeHTTP(w, r)
	})
}

// Authenticate middleware checks if the request context contains a valid user ID
func (app *Application) Authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := app.SessionManager.GetInt(r.Context(), "authenticatedUserID")
		if id == 0 {
			next.ServeHTTP(w, r)
			return
		}

		// check if a user with the `id` exists
		exists, err := app.Users.Exists(id)
		if err != nil {
			app.ServerError(w, err)
			return
		}

		if exists {
			// add isAuthenticated: true to the request context
			ctx := context.WithValue(r.Context(), isAuthenticatedContextKey, true)
			r = r.WithContext(ctx)
		}

		next.ServeHTTP(w, r)
	})
}

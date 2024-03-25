package main

import (
	"net/http"

	"github.com/justinas/nosurf"
)

// secureHeaders middleware sets important header fields in every response to avoid various types of
// attacks
func secureHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// load resources only from self or certain resources from google fonts
		csp := "default-src 'self'; style-src 'self' fonts.googleapis.com; font-src fonts.gstatic.com"
		w.Header().Set("Content-Security-Policy", csp)
		// strip URL information when redirecting out of the app
		w.Header().Set("Referrer-Policy", "origin-when-cross-origin")
		// block MIME type sniffing to avoid content sniffing attacks
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("X-Frame-Options", "deny")
		w.Header().Set("X-XSS-Protection", "0")

		next.ServeHTTP(w, r)
	})
}

func noSurf(next http.Handler) http.Handler {
	csrfHandler := nosurf.New(next)
	csrfHandler.SetBaseCookie(http.Cookie{
		HttpOnly: true,
		Secure:   true,
		Path:     "/",
	})

	return csrfHandler
}

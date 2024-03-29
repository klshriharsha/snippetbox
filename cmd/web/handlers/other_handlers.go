package handlers

import (
	"net/http"

	"github.com/klshriharsha/snippetbox/cmd/web/config"
)

// pingHandler simply responds with "OK"
func pingHandler(_ *config.Application) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("OK"))
	})
}

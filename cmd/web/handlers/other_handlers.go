package handlers

import (
	"net/http"
)

// pingHandler simply responds with "OK"
func pingHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("OK"))
}

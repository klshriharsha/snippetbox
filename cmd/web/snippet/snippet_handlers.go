package snippet

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/klshriharsha/snippetbox/cmd/web/config"
)

func SnippetViewHandler(_ *config.Application) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.Atoi(r.URL.Query().Get("id"))
		if err != nil || id < 1 {
			http.NotFound(w, r)
			return
		}

		fmt.Fprintf(w, "view snippet with id %d", id)
	}
}

func SnippetCreateHandler(_ *config.Application) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			w.Header().Set("Allow", http.MethodPost)
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}

		w.Write([]byte("create a new snippet"))
	}
}

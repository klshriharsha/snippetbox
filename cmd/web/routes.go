package main

import (
	"net/http"

	"github.com/klshriharsha/snippetbox/cmd/web/config"
	"github.com/klshriharsha/snippetbox/cmd/web/home"
	"github.com/klshriharsha/snippetbox/cmd/web/snippet"
)

func routes(app *config.Application) *http.ServeMux {
	mux := http.NewServeMux()

	fs := http.FileServer(staticFileSystem{http.Dir("./ui/static/")})
	// file server looks for the file under `./ui/static/`
	// so strip the `/static` prefix from request URL
	mux.Handle("/static/", http.StripPrefix("/static", fs))

	mux.HandleFunc("/", home.HomeHandler(app))
	mux.HandleFunc("/snippet/view", snippet.SnippetViewHandler(app))
	mux.HandleFunc("/snippet/create", snippet.SnippetCreateHandler(app))

	return mux
}

package main

import (
	"flag"
	"log"
	"net/http"
	"os"

	"github.com/klshriharsha/snippetbox/cmd/web/config"
	"github.com/klshriharsha/snippetbox/cmd/web/home"
	"github.com/klshriharsha/snippetbox/cmd/web/snippet"
)

func main() {
	// setup a commandline flag to override the network address.
	// run `go run ./cmd/web -help` for documentation on commandline flags
	addr := flag.String("addr", ":4000", "HTTP network address")
	flag.Parse()

	// create informational and error loggers
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	mux := http.NewServeMux()

	fs := http.FileServer(staticFileSystem{http.Dir("./ui/static/")})
	// file server looks for the file under `./ui/static/`
	// so strip the `/static` prefix from request URL
	mux.Handle("/static/", http.StripPrefix("/static", fs))

	app := &config.Application{
		ErrorLog: errorLog,
		InfoLog:  infoLog,
	}

	mux.HandleFunc("/", home.HomeHandler(app))
	mux.HandleFunc("/snippet/view", snippet.SnippetViewHandler(app))
	mux.HandleFunc("/snippet/create", snippet.SnippetCreateHandler(app))

	// create an http server with custom error logger
	srv := &http.Server{
		Addr:     *addr,
		Handler:  mux,
		ErrorLog: errorLog,
	}

	infoLog.Printf("starting server on %s", *addr)
	err := srv.ListenAndServe()
	errorLog.Fatal(err)
}

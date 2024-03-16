package main

import (
	"flag"
	"log"
	"net/http"
	"os"

	"github.com/klshriharsha/snippetbox/cmd/web/config"
)

func main() {
	// setup a commandline flag to override the network address.
	// run `go run ./cmd/web -help` for documentation on commandline flags
	addr := flag.String("addr", ":4000", "HTTP network address")
	flag.Parse()

	// create informational and error loggers
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	// for injecting dependencies to handlers
	app := &config.Application{
		ErrorLog: errorLog,
		InfoLog:  infoLog,
	}

	// create an http server with custom error logger
	srv := &http.Server{
		Addr:     *addr,
		Handler:  routes(app),
		ErrorLog: errorLog,
	}

	infoLog.Printf("starting server on %s", *addr)
	err := srv.ListenAndServe()
	errorLog.Fatal(err)
}

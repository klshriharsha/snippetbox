package main

import (
	"database/sql"
	"flag"
	"log"
	"net/http"
	"net/url"
	"os"

	"github.com/go-playground/form"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/klshriharsha/snippetbox/cmd/web/config"
	"github.com/klshriharsha/snippetbox/cmd/web/render"
	"github.com/klshriharsha/snippetbox/internal/models"
)

func main() {
	// setup a commandline flag to override the network address.
	// run `go run ./cmd/web -help` for documentation on commandline flags
	addr := flag.String("addr", ":4000", "HTTP network address")
	pgURL := flag.String(
		"pgurl",
		"postgresql://web:3c523592-852d-42be-915c-d5931792e39e@localhost:5432/postgres",
		"PostgreSQL URL",
	)
	flag.Parse()

	// create informational and error loggers
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	// attempt to connect to the PostgreSQL database
	db, err := connectDB(*pgURL)
	if err != nil {
		errorLog.Fatal(err)
	}
	defer db.Close()

	// parse all templates and cache them in memory to avoid disk access at runtime
	templateCache, err := render.NewTemplateCache()
	if err != nil {
		errorLog.Fatal(err)
	}

	// for decoding form data received in POST body into an interface
	formDecoder := form.NewDecoder()

	// for injecting dependencies to handlers
	app := &config.Application{
		ErrorLog:      errorLog,
		InfoLog:       infoLog,
		Snippets:      &models.SnippetModel{DB: db},
		TemplateCache: templateCache,
		FormDecoder:   formDecoder,
	}

	// create an http server with custom error logger
	srv := &http.Server{
		Addr:     *addr,
		Handler:  routes(app),
		ErrorLog: errorLog,
	}

	infoLog.Printf("starting server on %s", *addr)
	err = srv.ListenAndServe()
	errorLog.Fatal(err)
}

// connectDB creates an SQL connection pool and ensures a successful ping
func connectDB(pgURL string) (*sql.DB, error) {
	parsedURL, err := url.Parse(pgURL)
	if err != nil {
		return nil, err
	}

	conn, err := sql.Open("pgx", parsedURL.String())
	if err != nil {
		return nil, err
	}
	if err := conn.Ping(); err != nil {
		return nil, err
	}

	return conn, nil
}

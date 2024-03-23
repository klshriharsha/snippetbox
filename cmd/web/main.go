package main

import (
	"context"
	"crypto/tls"
	"flag"
	"log"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/alexedwards/scs/pgxstore"
	"github.com/alexedwards/scs/v2"
	"github.com/go-playground/form"
	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/klshriharsha/snippetbox/cmd/web/config"
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
	pool, err := connectDB(*pgURL)
	if err != nil {
		errorLog.Fatal(err)
	}
	defer pool.Close()

	// parse all templates and cache them in memory to avoid disk access at runtime
	templateCache, err := config.NewTemplateCache()
	if err != nil {
		errorLog.Fatal(err)
	}

	// for decoding form data received in POST body into an interface
	formDecoder := form.NewDecoder()

	// initialize a session manager
	sessionManager := scs.New()
	sessionManager.Store = pgxstore.New(pool)
	sessionManager.Lifetime = 12 * time.Hour

	// for injecting dependencies to handlers
	app := &config.Application{
		ErrorLog: errorLog,
		InfoLog:  infoLog,

		Snippets: &models.SnippetModel{DB: pool},
		Users:    &models.UserModel{DB: pool},

		TemplateCache:  templateCache,
		FormDecoder:    formDecoder,
		SessionManager: sessionManager,
	}

	// TLS configuration for higher security
	tlsConfig := tls.Config{
		CurvePreferences: []tls.CurveID{tls.X25519, tls.CurveP256},
	}

	// create an http server with custom error logger
	srv := &http.Server{
		Addr:         *addr,
		Handler:      routes(app),
		ErrorLog:     errorLog,
		TLSConfig:    &tlsConfig,
		IdleTimeout:  time.Minute,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	infoLog.Printf("starting server on %s", *addr)
	err = srv.ListenAndServeTLS("./tls/cert.pem", "./tls/key.pem")
	errorLog.Fatal(err)
}

// connectDB creates an SQL connection pool and ensures a successful ping
func connectDB(pgURL string) (*pgxpool.Pool, error) {
	parsedURL, err := url.Parse(pgURL)
	if err != nil {
		return nil, err
	}

	conn, err := pgxpool.New(context.Background(), parsedURL.String())
	if err != nil {
		return nil, err
	}
	if err := conn.Ping(context.Background()); err != nil {
		return nil, err
	}

	return conn, nil
}

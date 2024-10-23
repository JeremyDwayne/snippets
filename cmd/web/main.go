package main

import (
	"crypto/tls"
	"database/sql"
	"net/http"
	"os"
	"text/template"
	"time"

	"github.com/alexedwards/scs/sqlite3store"
	"github.com/alexedwards/scs/v2"
	"github.com/charmbracelet/log"
	"github.com/go-playground/form/v4"
	"github.com/jeremydwayne/snippets/internal/models"
	_ "github.com/mattn/go-sqlite3"
)

type application struct {
	logger         *log.Logger
	snippets       models.SnippetModelInterface
	users          models.UserModelInterface
	templateCache  map[string]*template.Template
	formDecoder    *form.Decoder
	sessionManager *scs.SessionManager
}

func main() {
	logger := log.NewWithOptions(os.Stdout, log.Options{
		ReportTimestamp: true,
	})

	dbName := "file:./local.db"
	db, err := openDB(dbName, os.Getenv("TURSO_DATABASE_URL"), os.Getenv("TURSO_AUTH_TOKEN"))
	if err != nil {
		log.Error(err.Error())
		os.Exit(1)
	}
	defer db.Close()

	templateCache, err := newTemplateCache()
	if err != nil {
		log.Error(err.Error())
		os.Exit(1)
	}

	formDecoder := form.NewDecoder()

	sessionManager := scs.New()
	sessionManager.Store = sqlite3store.New(db)
	sessionManager.Lifetime = 12 * time.Hour

	app := &application{
		logger:         logger,
		snippets:       &models.SnippetModel{DB: db},
		users:          &models.UserModel{DB: db},
		templateCache:  templateCache,
		formDecoder:    formDecoder,
		sessionManager: sessionManager,
	}

	stdlog := logger.StandardLog(log.StandardLogOptions{
		ForceLevel: log.ErrorLevel,
	})

	tlsConfig := &tls.Config{
		CurvePreferences: []tls.CurveID{tls.X25519, tls.CurveP256},
	}

	srv := &http.Server{
		Addr:         os.Getenv("HTTP_LISTEN_ADDR"),
		Handler:      app.routes(),
		ErrorLog:     stdlog,
		TLSConfig:    tlsConfig,
		IdleTimeout:  time.Minute,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	log.Info("Starting server", "addr", srv.Addr)

	err = srv.ListenAndServeTLS("./tls/cert.pem", "./tls/key.pem")
	logger.Error(err.Error())
	os.Exit(1)
}

func openDB(dbName string, primaryUrl string, authToken string) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", dbName)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		db.Close()
		return nil, err
	}

	return db, nil
}

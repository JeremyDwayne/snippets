package main

import (
	"crypto/tls"
	"database/sql"
	"embed"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"text/template"
	"time"

	// "github.com/alexedwards/scs/sqlite3store"
	"github.com/alexedwards/scs/v2"
	"github.com/charmbracelet/log"
	"github.com/go-playground/form/v4"
	"github.com/jeremydwayne/snippets/internal/libsqlstore"
	"github.com/jeremydwayne/snippets/internal/models"
	"github.com/jeremydwayne/snippets/internal/sqlc"
	"github.com/pressly/goose/v3"
	"github.com/tursodatabase/go-libsql"
)

type application struct {
	logger         *log.Logger
	db             *sql.DB
	snippets       models.SnippetModelInterface
	users          models.UserModelInterface
	templateCache  map[string]*template.Template
	formDecoder    *form.Decoder
	sessionManager *scs.SessionManager
}

//go:embed migrations
var migrations embed.FS

func main() {
	logger := log.NewWithOptions(os.Stdout, log.Options{
		ReportTimestamp: true,
	})

	dbName := "snippets.local.db"
	dbURL := os.Getenv("TURSO_DATABASE_URL")
	dbAuthToken := os.Getenv("TURSO_AUTH_TOKEN")

	dir, err := os.MkdirTemp("", "libsql-*")
	if err != nil {
		log.Error("Error creating temporary directory:", err)
		os.Exit(1)
	}
	defer os.RemoveAll(dir)

	dbPath := filepath.Join(dir, dbName)

	connector, err := libsql.NewEmbeddedReplicaConnector(dbPath, dbURL,
		libsql.WithAuthToken(dbAuthToken),
		libsql.WithEncryption(os.Getenv("DATABASE_SECRET")),
		libsql.WithSyncInterval(time.Minute),
	)
	if err != nil {
		log.Error("Error creating connector:", err)
		os.Exit(1)
	}
	defer connector.Close()

	db := sql.OpenDB(connector)
	defer db.Close()

	err = db.Ping()
	if err != nil {
		db.Close()
	}

	goose.SetBaseFS(migrations)
	err = goose.SetDialect("turso")
	if err != nil {
		db.Close()
	}

	err = goose.Version(db, "migrations")
	if err == goose.ErrNoCurrentVersion {
		log.Info("Inintializing Database")
	} else if err != nil {
		db.Close()
		log.Fatal(err)
	}

	err = goose.Up(db, "migrations")
	if err == goose.ErrAlreadyApplied {
		log.Info("No new migrations")
	} else if err != nil {
		db.Close()
		log.Fatal(err)
	} else {
		log.Info("Migrations ran")
	}

	if err != nil {
		log.Error(err.Error())
		os.Exit(1)
	}

	templateCache, err := newTemplateCache()
	if err != nil {
		log.Error(err.Error())
		os.Exit(1)
	}

	formDecoder := form.NewDecoder()

	sessionManager := scs.New()
	sessionManager.Store = libsqlstore.New(db)
	sessionManager.Lifetime = 12 * time.Hour

	queries := sqlc.New(db)

	app := &application{
		logger:         logger,
		db:             db,
		snippets:       &models.SnippetModel{DB: queries},
		users:          &models.UserModel{DB: queries},
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
		Addr:         fmt.Sprintf(":%s", os.Getenv("PORT")),
		Handler:      app.routes(),
		ErrorLog:     stdlog,
		TLSConfig:    tlsConfig,
		IdleTimeout:  time.Minute,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	log.Info("Starting server", "addr", srv.Addr)

	err = srv.ListenAndServe()
	logger.Error(err.Error())
	os.Exit(1)
}

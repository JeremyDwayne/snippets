package main

import (
	"crypto/tls"
	"database/sql"
	"fmt"
	"net/http"
	"os"
	"text/template"
	"time"

	"github.com/alexedwards/scs/sqlite3store"
	"github.com/alexedwards/scs/v2"
	"github.com/charmbracelet/log"
	"github.com/go-playground/form/v4"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/sqlite3"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	db "github.com/jeremydwayne/snippets/db/sqlc"
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

	database, err := openDB(os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Error(err.Error())
		os.Exit(1)
	}
	defer database.Close()

	templateCache, err := newTemplateCache()
	if err != nil {
		log.Error(err.Error())
		os.Exit(1)
	}

	formDecoder := form.NewDecoder()

	sessionManager := scs.New()
	sessionManager.Store = sqlite3store.New(database)
	sessionManager.Lifetime = 12 * time.Hour

	queries := db.New(database)

	app := &application{
		logger:         logger,
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
		Addr:         os.Getenv("HTTP_LISTEN_ADDR"),
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

func openDB(dbName string) (*sql.DB, error) {
	dbUrl := fmt.Sprintf("sqlite3://%s", dbName)
	log.Info(dbUrl)
	migrator, err := migrate.New("file://db/migrations", dbUrl)
	if err != nil {
		return nil, err
	}

	dbVersion, dbDirty, err := migrator.Version()
	if err == migrate.ErrNilVersion {
		log.Info("Inintializing Database")
	} else if err != nil {
		log.Fatal(err)
	}

	if dbDirty {
		dbForceVersion := dbVersion - 1
		log.Info("Database is dirty, forcing version", "version", dbForceVersion)
		err = migrator.Force(int(dbForceVersion))
		if err != nil {
			log.Fatal(err)
		}
	}

	log.Info("Database at version", "version", dbVersion)

	err = migrator.Up()
	if err == migrate.ErrNoChange {
		log.Info("No new migrations")
	} else if err != nil {
		log.Fatal(err)
	} else {
		log.Info("Migrations run")
	}

	database, err := sql.Open("sqlite3", dbName)
	if err != nil {
		return nil, err
	}

	err = database.Ping()
	if err != nil {
		database.Close()
		return nil, err
	}

	return database, nil
}

package main

import (
	"database/sql"
	"log/slog"
	"net/http"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jeremydwayne/snippets/internal/models"
)

type application struct {
	logger   *slog.Logger
	config   *config
	snippets *models.SnippetModel
}

type config struct {
	addr      string
	staticDir string
}

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	config := &config{
		addr:      os.Getenv("HTTP_LISTEN_ADDR"),
		staticDir: os.Getenv("STATIC_DIR"),
	}

	db, err := openDB(os.Getenv("DATABASE_URL"))
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}
	defer db.Close()

	app := &application{
		logger:   logger,
		config:   config,
		snippets: &models.SnippetModel{DB: db},
	}

	logger.Info("Starting server", "addr", config.addr)

	err = http.ListenAndServe(config.addr, app.routes())
	logger.Error(err.Error())
	os.Exit(1)
}

func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("mysql", dsn)
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

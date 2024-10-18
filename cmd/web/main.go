package main

import (
	"log/slog"
	"net/http"
	"os"
)

type application struct {
	logger *slog.Logger
	config *config
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

	app := &application{
		logger: logger,
		config: config,
	}

	logger.Info("Starting server", "addr", config.addr)

	err := http.ListenAndServe(config.addr, app.routes())
	logger.Error(err.Error())
	os.Exit(1)
}

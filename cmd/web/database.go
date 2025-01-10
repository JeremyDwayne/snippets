package main

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/charmbracelet/log"
	"github.com/pressly/goose/v3"
	"github.com/tursodatabase/go-libsql"
)

func openDB() (*sql.DB, error) {
	dbName := "snippets.local.db"
	primaryUrl := os.Getenv("TURSO_DATABASE_URL")
	authToken := os.Getenv("TURSO_AUTH_TOKEN")

	dir, err := os.MkdirTemp("", "libsql-*")
	if err != nil {
		fmt.Println("Error creating temporary directory:", err)
		os.Exit(1)
	}
	defer os.RemoveAll(dir)

	dbPath := filepath.Join(dir, dbName)

	connector, err := libsql.NewEmbeddedReplicaConnector(dbPath, primaryUrl,
		libsql.WithAuthToken(authToken),
		libsql.WithEncryption(os.Getenv("DATABASE_SECRET")),
		libsql.WithSyncInterval(time.Minute),
	)
	if err != nil {
		fmt.Println("Error creating connector:", err)
		os.Exit(1)
	}
	fmt.Println("after connector")
	defer connector.Close()

	log.Info("opendb")
	db := sql.OpenDB(connector)
	defer db.Close()

	err = db.Ping()
	if err != nil {
		db.Close()
		return nil, err
	}

	fmt.Println("before migrator")
	migrations := os.DirFS("db/migrations")
	goose.SetBaseFS(migrations)
	err = goose.SetDialect("turso")
	if err != nil {
		return nil, err
	}

	err = goose.Version(db, "db/migrations")
	if err == goose.ErrNoCurrentVersion {
		log.Info("Inintializing Database")
	} else if err != nil {
		log.Fatal(err)
	}

	err = goose.Up(db, "db/migrations")
	if err == goose.ErrAlreadyApplied {
		log.Info("No new migrations")
	} else if err != nil {
		log.Fatal(err)
	} else {
		log.Info("Migrations ran")
	}

	return db, nil
}

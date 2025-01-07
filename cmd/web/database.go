package main

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"

	"github.com/charmbracelet/log"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/source/file" // File source for migrations
	_ "github.com/jeremydwayne/snippets/internal/libsqlmigrate"
	"github.com/tursodatabase/go-libsql"
)

func openDB() (*sql.DB, error) {
	dbName := "local.db"
	primaryUrl := os.Getenv("TURSO_DATABASE_URL")
	authToken := os.Getenv("TURSO_AUTH_TOKEN")
	fmt.Println("env vars")

	dir, err := os.MkdirTemp("", "libsql-*")
	if err != nil {
		fmt.Println("Error creating temporary directory:", err)
		os.Exit(1)
	}
	defer os.RemoveAll(dir)
	fmt.Println("make temp dir")

	dbPath := filepath.Join(dir, dbName)

	connector, err := libsql.NewEmbeddedReplicaConnector(dbPath, primaryUrl,
		libsql.WithAuthToken(authToken),
	)
	fmt.Println("connector")

	if err != nil {
		fmt.Println("Error creating connector:", err)
		os.Exit(1)
	}
	defer connector.Close()

	// dbUrl := fmt.Sprintf("sqlite3://%s", dbName)
	// log.Info(dbUrl)
	fmt.Println("before migrator")
	migrator, err := migrate.New("file://db/migrations", "libsql://"+primaryUrl+"?authToken="+authToken)
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
		log.Info("error")
		log.Fatal(err)
	} else {
		log.Info("Migrations run")
	}

	log.Info("opendb")
	db := sql.OpenDB(connector)
	defer db.Close()

	err = db.Ping()
	if err != nil {
		db.Close()
		return nil, err
	}

	return db, nil
}

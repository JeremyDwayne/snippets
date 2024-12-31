package main

import (
	"database/sql"
	"fmt"

	"github.com/charmbracelet/log"
	"github.com/golang-migrate/migrate/v4"
)

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

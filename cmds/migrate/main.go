package main

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/sqlite3"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func main() {
	db, err := sql.Open("sqlite3", "diarier.db")
	if err != nil {
		log.Fatalf("Could not connect to db: %v\n", err)
	}

	driver, err := sqlite3.WithInstance(db, &sqlite3.Config{})
	if err != nil {
		log.Fatalf("Could not start sql driver: %v\n", err)
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://migrations", // replace with the path to your migrations folder
		"sqlite3",
		driver,
	)
	if err != nil {
		log.Fatalf("migration failed: %v\n", err)
	}

	err = m.Up()
	if err != nil {
		log.Fatalf("migration failed: %v\n", err)
	}
	fmt.Println("Successfully migrated")
}

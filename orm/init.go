package orm

import (
	"database/sql"
	"log"
)

const DBPATH = "diarier.db"

func GetDB() *sql.DB {
	db, err := sql.Open("sqlite3", DBPATH)
	if err != nil {
		log.Fatal(err)
	}

	if db == nil {
		log.Fatal("db nil")
	}
	return db
}

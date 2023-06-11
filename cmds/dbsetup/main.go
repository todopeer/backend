package main

import (
	"log"

	_ "github.com/mattn/go-sqlite3"
	"github.com/todopeer/backend/orm"
)

func main() {
	db := orm.GetDB()
	orm.CreateTables(db, log.Default())
}

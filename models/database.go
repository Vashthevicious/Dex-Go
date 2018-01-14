package models

import (
	"database/sql"
	"github.com/fortinj1354/Dex-Go/settings"
	_ "github.com/go-sql-driver/mysql"
	"time"
)

var db *sql.DB

func MakeDB(dbConnectString string) {
	newDb, err := sql.Open("mysql", dbConnectString)
	if err != nil {
		panic(err)
	}
	db = newDb
	timeout, err := time.ParseDuration(settings.GetDatabaseTimeout())
	if err != nil {
		panic(err)
	}
	db.SetConnMaxLifetime(timeout)
}

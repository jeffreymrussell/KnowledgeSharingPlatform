package bootstrap

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
)

func SetupDatabase(url string) *sql.DB {
	db, err := sql.Open("sqlite3", url)
	if err != nil {
		panic("Failed to open sqlite (" + url + ") Error: " + err.Error())
	}
	return db
}

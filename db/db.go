package db

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

func Conn() (*sql.DB, error) {
	db, err := sql.Open("sqlite3", "site.db")

	if err != nil {
		return db, err
	}

	return db, nil
}
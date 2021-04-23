package repo

import (
	"database/sql"
	"os"

	"golang.org/x/crypto/bcrypt"

	_ "github.com/mattn/go-sqlite3"
)

var salt = os.Getenv("SALT")

func Conn() (*sql.DB, error) {
	db, err := sql.Open("sqlite3", "site.db")

	if err != nil {
		return db, err
	}

	return db, nil
}

func HashFunc(str string) ([]byte, error) {
	byteStr := []byte(str+salt)

	hashedStr, err := bcrypt.GenerateFromPassword(byteStr, bcrypt.DefaultCost)
	if err != nil {
		return []byte(""), err
	}

	return hashedStr, nil
}

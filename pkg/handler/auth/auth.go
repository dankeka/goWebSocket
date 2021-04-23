package auth

import (
	"database/sql"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/dgrijalva/jwt-go"
)

func CheckAuth(r *http.Request, db *sql.DB, jwtToken jwt.MapClaims) (bool, error) {
	var err error
	var expJwt float64

	expJwt, err = strconv.ParseFloat(fmt.Sprintf("%v", jwtToken["exp"]), 64)
	
	if err != nil {
		return false, err
	}

	if expJwt <= float64(time.Now().Add(1 * time.Second).Unix()) {
		return false, nil
	}

	row := db.QueryRow("SELECT id FROM User WHERE id=$1", jwtToken["id"])
	
	var id sql.NullInt64
	err = row.Scan(&id)

	if err != nil {
		return false, err
	}
	
	if id.Int64 == 0 || !id.Valid {
		return false, fmt.Errorf("error: такого пользователя не существует")
	}
	
	return true, nil
}
package auth

import (
	"fmt"
	"os"
	"time"

	"github.com/dankeka/goWebSocket/pkg/types"
	"github.com/dgrijalva/jwt-go"
)

type JWTtoken struct {
	jwt.StandardClaims
	UserId   int    `json:"id"`
	UserName string `json:"name"`
}

func GenerateJWT(u types.User) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &JWTtoken{
		StandardClaims: jwt.StandardClaims {
			ExpiresAt: time.Now().Add(time.Hour * 720).Unix(), // add 30 days
			IssuedAt: time.Now().Unix(),
		},
		UserId: int(u.ID),
		UserName: u.Name,
	})

	return token.SignedString([]byte(os.Getenv("JWT_KEY")))
}

func ParseJWT(accessToken string) (jwt.MapClaims, error) {
	token, errParse := jwt.Parse(
		accessToken, 
		func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("error parse jwt token")
			}

			return []byte([]byte(os.Getenv("JWT_KEY"))), nil
		},
	)

	if errParse != nil {
		return map[string]interface{}{}, errParse
	}

	claims, ok := token.Claims.(jwt.MapClaims)

	if !ok {
		return map[string]interface{}{}, fmt.Errorf("error claims")
	}

	return claims, nil
}
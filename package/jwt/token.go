package token

import (
	"log"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/joho/godotenv"
)

type Claims struct {
	ID int64 `json:"ID"`
	jwt.StandardClaims
}

var key = []byte(os.Getenv("ACCESS_SECRET"))

// GenerateJWT for user
func GenerateJWT(id int64) (string, error) {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	expirationTime := time.Now().Add(2 * time.Hour)

	atClaims := &Claims{}
	atClaims.ID = id

	atClaims.StandardClaims.ExpiresAt = expirationTime.Unix()

	at := jwt.NewWithClaims(jwt.SigningMethodHS256, atClaims)

	AccessToken, err := at.SignedString(key)
	if err != nil {
		log.Fatal(err)
	}

	return AccessToken, nil
}

// VerifyToken params
func VerifyToken(tokenString string) (jwt.Claims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return key, nil
	})
	if err != nil {
		return nil, err
	}
	return token.Claims, err
}

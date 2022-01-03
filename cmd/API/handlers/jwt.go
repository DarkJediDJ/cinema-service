package handlers

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/joho/godotenv"
)

type Claims struct {
	UserID      int  `json:"UserID"`
	AddHalls    bool `json:"AddHalls"`
	AddMovies   bool `json:"AddMovies"`
	AddSessions bool `json:"AddSessions"`
	jwt.StandardClaims
}

var key = []byte(os.Getenv("ACCESS_SECRET"))

var ExpirationTime = time.Now().Add(5 * time.Minute)

func GenerateJWT(addHalls bool, addMovies bool, addSessions bool, userID int) (string, error) {
	var err error
	e := godotenv.Load(".env")
	if e != nil {
		log.Fatal("Error loading .env file")
	}
	expirationTime := time.Now().Add(2 * time.Hour)
	atClaims := &Claims{}
	atClaims.UserID = userID
	atClaims.AddHalls = addHalls
	atClaims.AddMovies = addMovies
	atClaims.AddSessions = addSessions
	atClaims.StandardClaims.ExpiresAt = expirationTime.Unix()
	at := jwt.NewWithClaims(jwt.SigningMethodHS256, atClaims)
	AccessToken, err := at.SignedString(key)
	if err != nil {
		log.Fatal(err)
	}

	return AccessToken, nil
}

func RefreshJWT(response http.ResponseWriter, tknStr string) (string, error) {
	claims := &Claims{}
	tkn, err := jwt.ParseWithClaims(tknStr, claims, func(token *jwt.Token) (interface{}, error) {
		return key, nil
	})
	if !tkn.Valid {
		response.WriteHeader(http.StatusUnauthorized)
		log.Fatal(err)
	}
	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			response.WriteHeader(http.StatusUnauthorized)
			log.Fatal(err)
		}
		response.WriteHeader(http.StatusBadRequest)
		log.Fatal(err)
	}
	if time.Unix(claims.ExpiresAt, 0).Sub(time.Now()) > 30*time.Second {
		response.WriteHeader(http.StatusBadRequest)
		log.Fatal(err)
	}
	claims.ExpiresAt = ExpirationTime.Unix()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(os.Getenv("ACCESS_SECRET")))
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		log.Fatal(err)
	}
	return tokenString, nil
}

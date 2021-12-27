package handlers

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"os"

	config "github.com/darkjedidj/cinema-service/internal"
	"github.com/dgrijalva/jwt-go"
	_ "github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
)

//Registration ...
func Authentification(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("Content-Type", "application/json")
	var user config.User
	var dbUser config.User
	json.NewDecoder(request.Body).Decode(&user)
	db, err := sql.Open("postgres", "host=localhost port=5432 user=admin dbname=admin password=admin sslmode=disable")
	if err != nil {
		log.Fatalln(err)
	}
	defer db.Close()
	rows, err := db.Query(`SELECT "Login", "Password","AddHalls","AddMovies","AddSessions" FROM "User" WHERE "Login" = $1`, user.Login)
	if err != nil {
		panic(err)
	}
	defer rows.Close()
	for rows.Next() {
		if err := rows.Scan(&dbUser.Login, &dbUser.Password, &dbUser.AddHalls, &dbUser.AddMovies, &dbUser.AddSessions); err != nil {
			panic(err)
		}
	}
	userPass := []byte(user.Password)
	dbPass := []byte(dbUser.Password)

	passErr := bcrypt.CompareHashAndPassword(dbPass, userPass)
	if passErr != nil {
		log.Println(passErr)
		response.Write([]byte(`{"response":"Wrong Password!"}`))
		return
	}
	jwtToken, err := GenerateJWT(dbUser.AddHalls, dbUser.AddMovies, dbUser.AddSessions)
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(`{"message":"` + err.Error() + `"}`))
		return
	}
	response.Write([]byte(`{"token":"` + jwtToken + `"}`))
}

func GenerateJWT(addHalls bool, addMovies bool, addSessions bool) (string, error) {
	var err error
	os.Setenv("ACCESS_SECRET", "jdnfksdmfksd") //this should be in an env file
	atClaims := jwt.MapClaims{}
	atClaims["addHalls"] = true
	at := jwt.NewWithClaims(jwt.SigningMethodHS256, atClaims)
	token, err := at.SignedString([]byte(os.Getenv("ACCESS_SECRET")))
	if err != nil {
		return "", err
	}
	return token, nil
}

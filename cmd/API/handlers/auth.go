package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	config "github.com/darkjedidj/cinema-service/internal"
	db "github.com/darkjedidj/cinema-service/internal/queries"
	_ "github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
)

//Authentification ...
func Authentification(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("Content-Type", "application/json")
	var user config.User
	json.NewDecoder(request.Body).Decode(&user)
	dbUser := db.SelectUser(user)
	userPass := []byte(user.Password)
	dbPass := []byte(dbUser.Password)

	passErr := bcrypt.CompareHashAndPassword(dbPass, userPass)
	if passErr != nil {
		log.Println(passErr)
		response.Write([]byte(`{"response":"Wrong Password!"}`))
		return
	}

	jwtToken, err := GenerateJWT(dbUser.AddHalls, dbUser.AddMovies, dbUser.AddSessions, dbUser.UserID)
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(`{"message":"` + err.Error() + `"}`))
		return
	}

	http.SetCookie(response, &http.Cookie{
		Name:    "token",
		Value:   jwtToken,
		Expires: CoockieTime,
	})

	_, err = http.Get("http://localhost:8085/refresh")
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		return
	}
	response.Write([]byte(jwtToken))

}
